package stream

import (
	"sort"
	"sync"
)

type PlaceholderType = struct{}

var Placeholder PlaceholderType

// drain drains the given channel.
func drain(channel <-chan interface{}) {
	for range channel {
	}
}

type Stream struct {
	source <-chan interface{}
}

// Concat returns a concatenated Stream.
func Concat(s Stream, others ...Stream) Stream {
	return s.Concat(others...)
}

// From constructs a Stream from the given GenerateFunc.
func From(generate GenerateFunc) Stream {
	source := make(chan interface{})

	GoSafe(func() {
		defer close(source)
		generate(source)
	})

	return Range(source)
}

// Just converts the given arbitrary items to a Stream.
func Just(items ...interface{}) Stream {
	source := make(chan interface{}, len(items))
	for _, item := range items {
		source <- item
	}
	close(source)

	return Range(source)
}

// Range converts the given channel to a Stream.
func Range(source <-chan interface{}) Stream {
	return Stream{
		source: source,
	}
}

// AllMach returns whether all elements of this stream match the provided predicate.
func (s Stream) AllMach(predicate func(item interface{}) bool) bool {
	for item := range s.source {
		if !predicate(item) {
			// make sure the former goroutine not block, and current func returns fast.
			go drain(s.source)
			return false
		}
	}

	return true
}

// AnyMach returns whether any elements of this stream match the provided predicate.
func (s Stream) AnyMach(predicate func(item interface{}) bool) bool {
	for item := range s.source {
		if predicate(item) {
			// make sure the former goroutine not block, and current func returns fast.
			go drain(s.source)
			return true
		}
	}

	return false
}

// Buffer buffers the items into a queue with size n.
func (s Stream) Buffer(n int) Stream {
	if n < 0 {
		n = 0
	}

	source := make(chan interface{}, n)
	go func() {
		for item := range s.source {
			source <- item
		}
		close(source)
	}()

	return Range(source)
}

// Concat returns a Stream that concatenated other streams
func (s Stream) Concat(others ...Stream) Stream {
	source := make(chan interface{})

	go func() {
		group := NewRoutineGroup()
		group.Run(func() {
			for item := range s.source {
				source <- item
			}
		})

		for _, each := range others {
			each := each
			group.Run(func() {
				for item := range each.source {
					source <- item
				}
			})
		}

		group.Wait()
		close(source)
	}()

	return Range(source)
}

// Count counts the number of elements in the result.
func (s Stream) Count() (count int) {
	for range s.source {
		count++
	}
	return
}

// Distinct removes the duplicated items base on the given KeyFunc.
func (s Stream) Distinct(fn KeyFunc) Stream {
	source := make(chan interface{})

	GoSafe(func() {
		defer close(source)

		keys := make(map[interface{}]PlaceholderType)
		for item := range s.source {
			key := fn(item)
			if _, ok := keys[key]; !ok {
				source <- item
				keys[key] = Placeholder
			}
		}
	})

	return Range(source)
}

// Done waits all upstreaming operations to be done.
func (s Stream) Done() {
	drain(s.source)
}

// Filter filters the items by the given FilterFunc.
func (s Stream) Filter(fn FilterFunc, opts ...Option) Stream {
	return s.Walk(func(item interface{}, pipe chan<- interface{}) {
		if fn(item) {
			pipe <- item
		}
	}, opts...)
}

// First returns the first item, nil if no items.
func (s Stream) First() interface{} {
	for item := range s.source {
		// make sure the former goroutine not block, and current func returns fast.
		go drain(s.source)
		return item
	}

	return nil
}

// ForAll handles the streaming elements from the source and no later streams.
func (s Stream) ForAll(fn ForAllFunc) {
	fn(s.source)
	// avoid goroutine leak on fn not consuming all items.
	go drain(s.source)
}

// ForEach seals the Stream with the ForEachFunc on each item, no successive operations.
func (s Stream) ForEach(fn ForEachFunc) {
	for item := range s.source {
		fn(item)
	}
}

// Group groups the elements into different groups based on their keys.
func (s Stream) Group(fn KeyFunc) Stream {
	groups := make(map[interface{}][]interface{})
	for item := range s.source {
		key := fn(item)
		groups[key] = append(groups[key], item)
	}

	source := make(chan interface{})
	go func() {
		for _, group := range groups {
			source <- group
		}
		close(source)
	}()

	return Range(source)
}

// Head returns the first n elements in p.
func (s Stream) Head(n int64) Stream {
	if n < 1 {
		panic("n must be greater than 0")
	}

	source := make(chan interface{})

	go func() {
		for item := range s.source {
			n--
			if n >= 0 {
				source <- item
			}
			if n == 0 {
				// let successive method go ASAP even we have more items to skip
				close(source)
				// why we don't just break the loop, and drain to consume all items.
				// because if breaks, this former goroutine will block forever,
				// which will cause goroutine leak.
				drain(s.source)
			}
		}
		// not enough items in s.source, but we need to let successive method to go ASAP.
		if n > 0 {
			close(source)
		}
	}()

	return Range(source)
}

// Last returns the last item, or nil if no items.
func (s Stream) Last() (item interface{}) {
	for item = range s.source {
	}
	return
}

// Map converts each item to another corresponding item, which means it's a 1:1 model.
func (s Stream) Map(fn MapFunc, opts ...Option) Stream {
	return s.Walk(func(item interface{}, pipe chan<- interface{}) {
		pipe <- fn(item)
	}, opts...)
}

// Merge merges all the items into a slice and generates a new stream.
func (s Stream) Merge() Stream {
	var items []interface{}
	for item := range s.source {
		items = append(items, item)
	}

	source := make(chan interface{}, 1)
	source <- items
	close(source)

	return Range(source)
}

// NoneMatch returns whether all elements of this stream don't match the provided predicate.
func (s Stream) NoneMatch(predicate func(item interface{}) bool) bool {
	for item := range s.source {
		if predicate(item) {
			// make sure the former goroutine not block, and current func returns fast.
			go drain(s.source)
			return false
		}
	}

	return true
}

// Parallel applies the given ParallelFunc to each item concurrently with given number of workers.
func (s Stream) Parallel(fn ParallelFunc, opts ...Option) {
	s.Walk(func(item interface{}, pipe chan<- interface{}) {
		fn(item)
	}, opts...).Done()
}

// Reduce is an utility method to let the caller deal with the underlying channel.
func (s Stream) Reduce(fn ReduceFunc) (interface{}, error) {
	return fn(s.source)
}

// Reverse reverses the elements in the stream.
func (s Stream) Reverse() Stream {
	var items []interface{}
	for item := range s.source {
		items = append(items, item)
	}
	// reverse, official method
	for i := len(items)/2 - 1; i >= 0; i-- {
		opp := len(items) - 1 - i
		items[i], items[opp] = items[opp], items[i]
	}

	return Just(items...)
}

// Skip returns a Stream that skips size elements.
func (s Stream) Skip(n int64) Stream {
	if n < 0 {
		panic("n must not be negative")
	}
	if n == 0 {
		return s
	}

	source := make(chan interface{})

	go func() {
		for item := range s.source {
			n--
			if n >= 0 {
				continue
			} else {
				source <- item
			}
		}
		close(source)
	}()

	return Range(source)
}

// Sort sorts the items from the underlying source.
func (s Stream) Sort(less LessFunc) Stream {
	var items []interface{}
	for item := range s.source {
		items = append(items, item)
	}
	sort.Slice(items, func(i, j int) bool {
		return less(items[i], items[j])
	})

	return Just(items...)
}

// Split splits the elements into chunk with size up to n,
// might be less than n on tailing elements.
func (s Stream) Split(n int) Stream {
	if n < 1 {
		panic("n should be greater than 0")
	}

	source := make(chan interface{})
	go func() {
		var chunk []interface{}
		for item := range s.source {
			chunk = append(chunk, item)
			if len(chunk) == n {
				source <- chunk
				chunk = nil
			}
		}
		if chunk != nil {
			source <- chunk
		}
		close(source)
	}()

	return Range(source)
}

// Tail returns the last n elements in p.
func (s Stream) Tail(n int64) Stream {
	if n < 1 {
		panic("n should be greater than 0")
	}

	source := make(chan interface{})

	go func() {
		ring := NewRing(int(n))
		for item := range s.source {
			ring.Add(item)
		}
		for _, item := range ring.Take() {
			source <- item
		}
		close(source)
	}()

	return Range(source)
}

// Walk lets the callers handle each item, the caller may write zero, one or more items base on the given item.
func (s Stream) Walk(fn WalkFunc, opts ...Option) Stream {
	option := buildOptions(opts...)
	if option.unlimitedWorkers {
		return s.walkUnlimited(fn, option)
	}

	return s.walkLimited(fn, option)
}

func (s Stream) walkLimited(fn WalkFunc, option *rxOptions) Stream {
	pipe := make(chan interface{}, option.workers)

	go func() {
		var wg sync.WaitGroup
		pool := make(chan PlaceholderType, option.workers)

		for item := range s.source {
			// important, used in another goroutine
			val := item
			pool <- Placeholder
			wg.Add(1)

			// better to safely run caller defined method
			GoSafe(func() {
				defer func() {
					wg.Done()
					<-pool
				}()

				fn(val, pipe)
			})
		}

		wg.Wait()
		close(pipe)
	}()

	return Range(pipe)
}

func (s Stream) walkUnlimited(fn WalkFunc, option *rxOptions) Stream {
	pipe := make(chan interface{}, option.workers)

	go func() {
		var wg sync.WaitGroup

		for item := range s.source {
			// important, used in another goroutine
			val := item
			wg.Add(1)
			// better to safely run caller defined method
			GoSafe(func() {
				defer wg.Done()
				fn(val, pipe)
			})
		}

		wg.Wait()
		close(pipe)
	}()

	return Range(pipe)
}
