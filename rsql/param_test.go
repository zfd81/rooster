package rsql

import (
	"testing"
	"time"
)

func TestNewMapParams(t *testing.T) {
	m := make(map[string]interface{})
	NewMapParams(m)
}

func TestNewStructParams(t *testing.T) {
	u := &User{
		Id:       22,
		Name:     "user22",
		Password: "pwd22",
		FullName: "用户22",
		Number:   "num22",
		Model: Model{
			Creator:          1,
			CreatedDate:      time.Now(),
			Modifier:         1,
			LastmodifiedDate: time.Now(),
		},
		Field1: "test",
		Field2: 999,
	}
	p := NewStructParams(u)

	for k, v := range p {
		t.Log(k, v)
	}
}
