package golru

import "testing"

func TestNewList(t *testing.T) {
	l := newList()
	if l.Size() != 0 {
		t.Errorf("size error")
	}

	elem1 := newElement(1)
	l.InsertFront(elem1)
	if l.Size() != 1 {
		t.Errorf("size error")
	}

	elem2 := newElement(2)
	l.InsertTail(elem2)
	if l.Size() != 2 {
		t.Errorf("size error")
	}

	fElem := l.Front()
	if fElem != elem1 {
		t.Errorf("get front error")
	}
	if fElem.Val() != 1 {
		t.Errorf("get front error")
	}

	elem2.Bubble()
	fElem = l.Front()
	if fElem != elem2 {
		t.Errorf("bubble error")
	}
	if fElem.Val() != 2 {
		t.Errorf("bubble error")
	}

	elem2.Sink()
	fElem = l.Front()
	if fElem != elem1 {
		t.Errorf("sink error")
	}
	if fElem.Val() != 1 {
		t.Errorf("sink error")
	}
}
