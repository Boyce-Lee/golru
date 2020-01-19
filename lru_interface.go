package golru

type LruCache interface {
	Size() int
	Cap() int
	Empty() bool
	Get(key, defaultVal interface{}) interface{}
	MGet(keyList []interface{}, defVal interface{}) []interface{}
	Put(key, value interface{}, exp int64)
	MPut(keyList, valueList []interface{}, exp int64)
}
