package container

type Map interface {
	Put(key interface{}, value interface{})
	Get(key interface{}) (value interface{}, found bool)
	GetString(key interface{}) string
	GetInt(key interface{}) int
	GetBool(key interface{}) bool
	Remove(key interface{})
	Keys() []interface{}
	Map() map[interface{}]interface{}

	Container
	// Empty() bool
	// Size() int
	// Clear()
	// Values() []interface{}
}

type HashMap map[interface{}]interface{}

func (m HashMap) Put(key interface{}, value interface{}) {
	m[key] = value
}

func (m HashMap) Get(key interface{}) (value interface{}, found bool) {
	value, found = m[key]
	return
}

func (m HashMap) GetString(key interface{}) string {
	value, found := m.Get(key)
	if found {
		return value.(string)
	}
	return ""
}

func (m HashMap) GetInt(key interface{}) int {
	value, found := m.Get(key)
	if found {
		return value.(int)
	}
	return 0
}

func (m HashMap) GetBool(key interface{}) bool {
	value, found := m.Get(key)
	if found {
		return value.(bool)
	}
	return false
}

func (m HashMap) Remove(key interface{}) {
	delete(m, key)
}

func (m HashMap) Keys() []interface{} {
	keys := make([]interface{}, m.Size())
	count := 0
	for key := range m {
		keys[count] = key
		count++
	}
	return keys
}

func (m HashMap) Empty() bool {
	return m.Size() == 0
}

func (m HashMap) Size() int {
	return len(m)
}

func (m HashMap) Clear() {
	m = make(map[interface{}]interface{})
}

func (m HashMap) Values() []interface{} {
	values := make([]interface{}, m.Size())
	count := 0
	for _, value := range m {
		values[count] = value
		count++
	}
	return values
}

func (m HashMap) Map() map[interface{}]interface{} {
	return m
}
