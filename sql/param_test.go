package sql

import (
	"testing"
)

func TestNewMapParams(t *testing.T) {
	m := make(map[string]interface{})
	NewMapParams(m)
}
