package container

type Set interface {
	Add(value interface{})
	Contains(value interface{}) bool
	Remove(value interface{})
	Iterator(f func(index int, value interface{}))
	Container
}

type HashSet map[interface{}]void

func (s HashSet) Add(value interface{}) {
	s[value] = empty
}

func (s HashSet) Contains(value interface{}) bool {
	_, found := s[value]
	return found
}

func (s HashSet) Remove(value interface{}) {
	delete(s, value)
}

func (s HashSet) Iterator(f func(index int, value interface{})) {
	count := 0
	for k := range s {
		f(count, k)
		count++
	}
}

func (s HashSet) Empty() bool {
	return len(s) == 0
}

func (s HashSet) Size() int {
	return len(s)
}

func (s HashSet) Clear() {
	s = make(map[interface{}]void)
}

func (s HashSet) Values() []interface{} {
	values := make([]interface{}, s.Size())
	count := 0
	for k := range s {
		values[count] = k
		count++
	}
	return values
}
