package golru

import (
	"sync"
	"time"
)

// simple implementation of lru cache
// hash_map + doubled linked elemList
// doubled linked elemList maintains the item elemList sorted by travel time
// hash_map uses to get the item efficiently
// lock for any Get/Put operation
type simpleExpireLruCache struct {
	elemM map[interface{}] *element // key -> *element{Val: *simpleLruItem, exp: XXX}
	elemL *list
	max int
	exp int64 // unit second
	mu sync.RWMutex
}

var _ LruCache = new(simpleExpireLruCache)

type simpleLruItem struct {
	Key interface{}
	Val interface{}
	exp int64
}

func NewSimpleExpireLruCache(size int, exp int64) *simpleExpireLruCache {
	lru := &simpleExpireLruCache{
		elemM: nil,
		elemL: nil,
		max:   size,
		exp : exp,
	}
	lru.elemL = newList()
	lru.elemM = make(map[interface{}]*element)

	return lru
}

func (s *simpleExpireLruCache) Size() int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.elemL.Size()
}

func (s *simpleExpireLruCache) Empty() bool {
	return s.Size() == 0
}

func (s *simpleExpireLruCache) Cap() int {
	return s.max
}

func (s *simpleExpireLruCache) MGet(keyList []interface{}, defaultVal interface{}) []interface{} {
	ret := make([]interface{}, len(keyList))
	for idx := range ret {
		ret[idx] = s.Get(keyList[idx], defaultVal)
	}
	return ret
}

func (s *simpleExpireLruCache) Get(key, defaultVal interface{}) interface{} {
	s.mu.RLock()
	elem, exist := s.elemM[key]
	if !exist {
		return defaultVal
	}
	s.mu.RUnlock()

	item := elem.Val().(*simpleLruItem)

	if item.exp < time.Now().Unix() {
		return defaultVal
	}

	// promotion
	defer s.promotion(elem)

	return item.Val
}

func (s *simpleExpireLruCache) MPut(keyList, valList []interface{}, dura int64) {
	size := len(keyList)
	if len(valList) < size {
		size = len(valList)
	}

	for i := 0; i < size; i++ {
		s.Put(keyList[i], valList[i], dura)
	}
}

func (s *simpleExpireLruCache) Put(key, val interface{}, dura int64) {
	exp := s.exp
	if dura != 0 {
		exp = dura
	}
	elem, exist := s.elemM[key]
	if exist {
		item := elem.Val().(*simpleLruItem)
		item.Val = val
		item.exp = time.Now().Unix() + exp
		defer s.promotion(elem)
		return
	}

	item := &simpleLruItem{
		Key: key,
		Val: val,
		exp: time.Now().Unix() + exp,
	}
	elem = newElement(item)

	s.mu.Lock()
	s.elemM[key] = elem
	s.elemL.InsertFront(elem)
	s.mu.Unlock()

	defer s.cleanUpTail()
}

func (s *simpleExpireLruCache) promotion(elem *element) {
	s.mu.Lock()
	defer s.mu.Unlock()

	elem.Bubble()
}

func (s *simpleExpireLruCache) cleanUpTail() {
	size := s.elemL.Size()
	if size < s.max {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	deleNum := size - s.max + 7
	for i := 0; i < deleNum; i++ {
		tail := s.elemL.Back()
		item := tail.Val().(*simpleLruItem)
		s.elemL.DeleteTail()
		delete(s.elemM, item.Key)
	}
}
