package common

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

func (s *Stack[T]) Peek() Option[T] {
	if len(s.data) == 0 {
		return None[T]()
	}

	return Some(s.data[len(s.data)-1])
}

func (s *Stack[T]) Push(item T) {
	s.data = append(s.data, item)
}

func (s *Stack[T]) Pop() Option[T] {
	if len(s.data) == 0 {
		return None[T]()
	}

	res := Some(s.data[len(s.data)-1])
	s.data = s.data[:len(s.data)-1]

	return res
}
