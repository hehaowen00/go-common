package ring

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

func (r *Ring[T]) Clear() {
	clear(r.data)
	r.head = 0
	r.tail = 0
	r.len = 0
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

func (r *Ring[T]) Enqueue(data []T) {
	if r.cap-r.len < len(data) {
		v := r.cap + len(data)
		r.resize(pow2(int64(v)))
	}

	for _, v := range data {
		r.data[r.tail] = v
		r.tail++
	}

	r.len += len(data)
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

func (r *Ring[T]) Dequeue() []T {
	var data []T

	if r.len == 0 {
		return nil
	}

	if r.head < r.tail {
		data = append(data, r.data[r.head:r.tail]...)
	} else {
		data = append(data, r.data[r.head:r.cap]...)
		data = append(data, r.data[0:r.tail]...)
	}

	r.head = 0
	r.tail = 0
	r.len = 0
	clear(r.data)

	return data
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
	r.cap = cap
}

func pow2(v int64) int {
	v--
	v |= v >> 1
	v |= v >> 2
	v |= v >> 4
	v |= v >> 8
	v |= v >> 16
	v |= v >> 32
	v++
	return int(v)
}
