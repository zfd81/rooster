package stream

type (
	// FilterFunc defines the method to filter a Stream.
	FilterFunc func(item interface{}) bool
	// ForAllFunc defines the method to handle all elements in a Stream.
	ForAllFunc func(pipe <-chan interface{})
	// ForEachFunc defines the method to handle each element in a Stream.
	ForEachFunc func(item interface{})
	// GenerateFunc defines the method to send elements into a Stream.
	GenerateFunc func(source chan<- interface{})
	// KeyFunc defines the method to generate keys for the elements in a Stream.
	KeyFunc func(item interface{}) interface{}
	// LessFunc defines the method to compare the elements in a Stream.
	LessFunc func(a, b interface{}) bool
	// MapFunc defines the method to map each element to another object in a Stream.
	MapFunc func(item interface{}) interface{}
	// Option defines the method to customize a Stream.
	Option func(opts *rxOptions)
	// ParallelFunc defines the method to handle elements parallelly.
	ParallelFunc func(item interface{})
	// ReduceFunc defines the method to reduce all the elements in a Stream.
	ReduceFunc func(pipe <-chan interface{}) (interface{}, error)
	// WalkFunc defines the method to walk through all the elements in a Stream.
	WalkFunc func(item interface{}, pipe chan<- interface{})
)
