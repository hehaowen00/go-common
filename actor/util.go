package actor

func GetMessage[T any](m *Message) (T, bool) {
	v, ok := m.Data.(T)
	return v, ok
}

func GetState[T any](s *State) (*T, bool) {
	v, ok := s.state.(*T)
	return v, ok
}
