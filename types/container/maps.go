package container

import (
	"encoding/json"

	"github.com/spf13/cast"
)

type Map interface {
	Put(key interface{}, value interface{})
	Get(key interface{}) (value interface{}, found bool)
	GetString(key interface{}) string
	GetInt(key interface{}) int
	GetFloat(key interface{}) float64
	GetBool(key interface{}) bool
	Remove(key interface{})
	Keys() []interface{}

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
		v, ok := value.(int64)
		if ok {
			return int(v)
		}
		return value.(int)
	}
	return 0
}

func (m HashMap) GetFloat(key interface{}) float64 {
	value, found := m.Get(key)
	if found {
		return value.(float64)
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

// New instantiates a new empty HashMap
func NewHashMap() HashMap {
	return HashMap{}
}

type JsonMap map[string]interface{}

func (m JsonMap) Put(key interface{}, value interface{}) {
	m[cast.ToString(key)] = value
}

func (m JsonMap) Get(key interface{}) (value interface{}, found bool) {
	value, found = m[cast.ToString(key)]
	return
}

func (m JsonMap) GetString(key interface{}) string {
	value, found := m.Get(key)
	if found {
		return value.(string)
	}
	return ""
}

func (m JsonMap) GetInt(key interface{}) int {
	value, found := m.Get(key)
	if found {
		v, ok := value.(int64)
		if ok {
			return int(v)
		}
		return value.(int)
	}
	return 0
}

func (m JsonMap) GetFloat(key interface{}) float64 {
	value, found := m.Get(key)
	if found {
		return value.(float64)
	}
	return 0
}

func (m JsonMap) GetBool(key interface{}) bool {
	value, found := m.Get(key)
	if found {
		return value.(bool)
	}
	return false
}

func (m JsonMap) Remove(key interface{}) {
	delete(m, cast.ToString(key))
}

func (m JsonMap) Keys() []interface{} {
	keys := make([]interface{}, m.Size())
	count := 0
	for key := range m {
		keys[count] = key
		count++
	}
	return keys
}

func (m JsonMap) Empty() bool {
	return m.Size() == 0
}

func (m JsonMap) Size() int {
	return len(m)
}

func (m JsonMap) Clear() {
	m = make(map[string]interface{})
}

func (m JsonMap) Values() []interface{} {
	values := make([]interface{}, m.Size())
	count := 0
	for _, value := range m {
		values[count] = value
		count++
	}
	return values
}

func (m JsonMap) Map() map[string]interface{} {
	return m
}

func (m JsonMap) Marshal() ([]byte, error) {
	return json.Marshal(m)
}

// New instantiates a new empty JsonMap
func NewJsonMap() JsonMap {
	return JsonMap{}
}
