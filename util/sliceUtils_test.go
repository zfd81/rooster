package util

import (
	"testing"
)

func TestInsertIntSlice(t *testing.T) {
	src := []int{1, 2, 3, 4}
	ins := []int{5, 6}
	src = InsertIntSlice(src, ins, 4)
	t.Log(src)
	i := 0
	for len(src) != 0 {
		if i < 3 {
			src = append(src, i)
		}
		t.Log(src[i])
		i++
	}
}
