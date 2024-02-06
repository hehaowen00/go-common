package actor

import (
	"bytes"
	"encoding/gob"
)

func GetMessage[T any](m *Message) (T, error) {
	buf := bytes.NewBuffer(m.Data)

	dec := gob.NewDecoder(buf)

	var v T

	err := dec.Decode(&v)
	if err != nil {
		return v, err
	}

	return v, err
}

func GetState[T any](s *State) (*T, bool) {
	v, ok := s.state.(*T)
	return v, ok
}
