package common

type Queue[T any] struct {
	head *node[T]
	tail *node[T]

	freedHead *node[T]
	freedTail *node[T]

	len  int
	free int
}

type node[T any] struct {
	data T
	next *node[T]
}

func NewQueue[T any]() *Queue[T] {
	return &Queue[T]{}
}

func (q *Queue[T]) Len() int {
	return q.len
}

func (q *Queue[T]) Enqueue(item T) {
	var n *node[T]

	if q.freedHead != nil {
		n = q.freedHead
		q.freedHead = n.next
		q.free--
	} else {
		n = new(node[T])
	}

	n.data = item

	if q.len == 0 {
		q.head = n
		q.tail = n
	} else {
		q.tail.next = n
		q.tail = n
	}

	q.len++
}

func (q *Queue[T]) Dequeue() (T, bool) {
	var tmp T
	if q.len == 0 {
		return tmp, false
	}

	n := q.head
	item := n.data
	q.head = q.head.next

	n.data = tmp
	n.next = nil

	q.len--

	if q.freedHead == nil {
		q.freedHead = n
		q.freedTail = n
	} else {
		q.freedTail.next = n
		q.freedTail = n
	}

	q.free++

	return item, true
}
