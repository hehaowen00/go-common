package stack

type Stack[T any] struct {
	data []T
}

func NewStack[T any]() *Stack[T] {
	return &Stack[T]{}
}

func (s *Stack[T]) Len() int {
	return len(s.data)
}

func (s *Stack[T]) IsEmpty() bool {
	return len(s.data) == 0
}

func (s *Stack[T]) Peek() (T, bool) {
	if len(s.data) == 0 {
		return *new(T), false
	}

	return s.data[len(s.data)-1], true
}

func (s *Stack[T]) Push(item T) {
	s.data = append(s.data, item)
}

func (s *Stack[T]) Pop() (T, bool) {
	if len(s.data) == 0 {
		return *new(T), false
	}

	res := s.data[len(s.data)-1]
	s.data = s.data[:len(s.data)-1]

	return res, true
}
