package common

type Ring[T any] struct {
	data []T
	head int
	tail int
	cap  int
	len  int
}

func NewRing[T any]() *Ring[T] {
	return &Ring[T]{
		data: make([]T, 1),
		cap:  1,
	}
}

func (r *Ring[T]) Len() int {
	return len(r.data)
}

func (r *Ring[T]) Push(item T) {
	if r.len > 0 && r.tail == r.head {
		r.resize(r.cap * 2)
	}

	r.data[r.tail] = item
	r.tail++

	if r.tail == len(r.data) {
		r.tail = r.tail % len(r.data)
	}

	r.len++
}

func (r *Ring[T]) Pop() (T, bool) {
	var tmp T

	if r.len == 0 {
		return tmp, false
	}

	res := r.data[r.head]
	r.data[r.head] = tmp
	r.head = (r.head + 1) % len(r.data)
	r.len--

	return res, true
}

func (r *Ring[T]) resize(cap int) {
	newData := make([]T, cap)
	count := 0

	for i := 0; i < r.len; i++ {
		start := (r.head + i) % len(r.data)
		newData[count] = r.data[start]
		count++
	}

	r.head = 0
	r.tail = count
	r.data = newData
	r.cap *= 2
}
