package container

import (
	"testing"
)

func TestNewArrayStack(t *testing.T) {
	var s Stack = NewArrayStack()
	s.Push("hello")
	s.Push("word")

	v, err := s.Pop()
	if err != nil {
		t.Error(err)
	}
	t.Log(v)

	v, err = s.Peek()
	if err != nil {
		t.Error(err)
	}

	t.Log(v)
	t.Log(s.Size())
	t.Log(s.Values())
	s.Clear()
	t.Log(s.Empty())
}
