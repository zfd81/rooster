package container

import (
	"errors"
)

const (
	DefaultStackCapacity = 10
)

type Stack interface {

	//Pushes an item onto the top of this stack.
	Push(value interface{})

	//Removes the object at the top of this stack and returns that object as the value of this function.
	Pop() (interface{}, error)

	//Looks at the object at the top of this stack without removing it from the stack.
	Peek() (interface{}, error)

	Container
	// Empty() bool
	// Size() int
	// Clear()
	// Values() []interface{}
}

// Stack holds elements in an slice
type ArrayStack []interface{}

func (s *ArrayStack) Push(value interface{}) {
	*s = append(*s, value)
}

func (s *ArrayStack) Pop() (interface{}, error) {
	if len(*s) == 0 {
		return nil, errors.New("Out of index, len is 0")
	}
	value := (*s)[len(*s)-1]
	*s = (*s)[:len(*s)-1]
	return value, nil
}

func (s *ArrayStack) Peek() (interface{}, error) {
	if len(*s) == 0 {
		return nil, errors.New("Out of index, len is 0")
	}
	return (*s)[len(*s)-1], nil
}

// Empty returns true if stack does not contain any elements.
func (s *ArrayStack) Empty() bool {
	return len(*s) == 0
}

// Size returns number of elements within the stack.
func (s *ArrayStack) Size() int {
	return len(*s)
}

// Clear removes all elements from the stack.
func (s *ArrayStack) Clear() {
	*s = (*s)[0:0]
}

// Values returns all elements in the stack (LIFO order).
func (s *ArrayStack) Values() []interface{} {
	return *s
}

// New instantiates a new empty stack
func NewArrayStack() *ArrayStack {
	s := make(ArrayStack, 0, DefaultStackCapacity)
	return &s
}
