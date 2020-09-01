package util

func InsertIntSlice(slice, insertion []int, index int) []int {
	length := len(slice)
	if index > length {
		index = length
	}
	result := make([]int, length+len(insertion))
	at := copy(result, slice[:index])
	at += copy(result[at:], insertion)
	copy(result[at:], slice[index:])
	return result
}
