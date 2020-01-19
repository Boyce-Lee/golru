package golru

type element struct {
	pre      *element
	next     *element
	val      interface{}
	elemList *list
}

func newElement(val interface{}) *element {
	return &element{
		pre:      nil,
		next:     nil,
		val:      val,
		elemList: nil,
	}
}

// not thread-safe
type list struct {
	// sentinel of head
	// sentinel.next is the first element
	// sentinel.pre is the tail element
	sentinel *element
	size     int
}

func newList() *list {
	l := &list{}
	sentinel := &element{}
	sentinel.pre, sentinel.next = sentinel, sentinel
	sentinel.elemList = l
	l.sentinel = sentinel
	return l
}

func (e *element) Val() interface{} {
	return e.val
}

func (e *element) Next() *element {
	return e.next
}

func (e *element) Pre() *element {
	return e.pre
}

func (e *element) Bubble() {
	if e.elemList == nil {
		return
	}
	e.remove()
	e.elemList.InsertFront(e)
}

func (e *element) Sink() {
	e.remove()
	e.elemList.InsertTail(e)
}

func (e *element) remove() {
	pre, next := e.pre, e.next
	if pre != nil {
		pre.next = next
	}
	if next != nil {
		next.pre = pre
	}

	e.elemList.size--
}

func (l *list) Size() int {
	return l.size
}

func (l *list) InsertFront(elem *element) {
	first := l.sentinel.next

	l.sentinel.next = elem
	elem.pre = l.sentinel

	elem.next = first
	first.pre = elem

	elem.elemList = l

	l.size++
}

func (l *list) InsertTail(elem *element) {
	tail := l.sentinel.pre

	l.sentinel.pre = elem
	elem.next = l.sentinel

	tail.next = elem
	elem.pre = tail

	elem.elemList = l

	l.size++
}

func (l *list) DeleteTail() {
	if l.size == 0 {
		return
	}

	tail := l.sentinel.pre
	tPre := tail.pre

	tPre.next = l.sentinel
	l.sentinel.pre = tPre
}

func (l *list) Front() *element {
	if l.size == 0 {
		return nil
	}
	return l.sentinel.next
}

func (l *list) Back() *element {
	if l.size == 0 {
		return nil
	}
	return l.sentinel.pre
}