package container

type void struct{}

var empty void

type Container interface {
	Empty() bool
	Size() int
	Clear()
	Values() []interface{}
}
