package game

import "testing"

func TestQueue_Push(t *testing.T) {

	q := Queue{}

	p1 := &People{}
	q.Push(p1)
	if q.head == p1 && q.tail == p1 {
		t.Log("TestQueue_Push one is Right!")
	} else {
		t.Error("TestQueue_Push one is Wrong!")
	}

	p2 := &People{}
	q.Push(p2)
	if q.head == p1 && q.tail == p2 {
		t.Log("TestQueue_Push two is Right!")
	} else {
		t.Error("TestQueue_Push two is Wrong!")
	}

	p3 := &People{}
	q.Push(p3)
	if p1.next == p2 && p2.next == p3 {
		t.Log("TestQueue_Push three is Right!")
	} else {
		t.Error("TestQueue_Push three is Wrong!")
	}
}

func TestQueue_Pop(t *testing.T) {

	q := Queue{}
	p1 := &People{}
	q.Push(p1)

	if q.Pop() == p1 && q.head == nil && q.tail == nil {
		t.Log("TestQueue_Pop is Right!")
	} else {
		t.Error("TestQueue_Pop is Wrong!")
	}
}

func TestQueue_Pop2(t *testing.T) {

	q := Queue{}

	p1 := &People{}
	q.Push(p1)
	p2 := &People{}
	q.Push(p2)

	q.Pop()

	if q.head == p2 && q.tail == p2 {
		t.Log("TestQueue_Pop2 is Right!")
	} else {
		t.Error("TestQueue_Pop2 is Wrong!")
	}
}

func TestQueue_Pop3(t *testing.T) {

	q := Queue{}

	p1 := &People{}
	q.Push(p1)
	p2 := &People{}
	q.Push(p2)
	p3 := &People{}
	q.Push(p3)

	q.Pop()

	if q.head == p2 && q.tail == p3 {
		t.Log("TestQueue_Pop3 is Right!")
	} else {
		t.Error("TestQueue_Pop3 is Wrong!")
	}
}
