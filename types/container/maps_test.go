package container

import (
	"encoding/json"
	"testing"
)

func TestNewHashMap(t *testing.T) {
	m := NewHashMap()
	m.Put("name", "zfd")
	m.Put("age", 40)
	m.Put("sex", "man")
	jsonByte, err := json.Marshal(m)
	if err != nil {
		t.Error(err)
	}
	t.Log(string(jsonByte))
}

func TestNewJsonMap(t *testing.T) {
	m := NewJsonMap()
	m.Put("name", "zfd")
	m.Put("age", 40)
	m.Put("sex", "man")
	jsonByte, err := m.Marshal()
	if err != nil {
		t.Error(err)
	}
	t.Log(string(jsonByte))
}
