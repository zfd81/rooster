package util

import (
	"bytes"
	"fmt"
	"strings"
	"unicode/utf8"
)

type ReplacerFunc func(index int, start int, end int, content string) (string, error)

func Left(str string, length int) string {
	if str == "" || length < 0 {
		return ""
	}
	strRune := []rune(str)
	if length < len(strRune) {
		return string(strRune[:length])
	} else {
		return str
	}
}

func Right(str string, length int) string {
	if str == "" || length < 0 {
		return ""
	}
	strRune := []rune(str)
	strLen := len(strRune)
	if length < strLen {
		return string(strRune[strLen-length:])
	} else {
		return str
	}
}

func Substr(str string, beginIndex int, endIndex int) (string, error) {
	if str == "" {
		return "", nil
	}
	if beginIndex < 0 {
		return "", fmt.Errorf("String index out of range: %d", beginIndex)
	}
	strRune := []rune(str)
	strLen := len(strRune)
	if endIndex > strLen {
		return "", fmt.Errorf("String index out of range: %d", endIndex)
	}
	subLen := endIndex - beginIndex
	if subLen < 0 {
		return "", fmt.Errorf("String index out of range: %d", subLen)
	}
	return string(strRune[beginIndex:endIndex]), nil
}

func ReplaceBetween(str string, open string, close string, replacer ReplacerFunc) (string, error) {
	if str == "" {
		return "", nil
	}
	strLen := utf8.RuneCountInString(str)
	openLen := len(open)
	closeLen := len(close)
	pos := 0
	index := 0
	var buffer bytes.Buffer
	for {
		if pos < strLen-closeLen {
			start := strings.Index(str[pos:], open)
			if start < 0 {
				break
			}
			start = pos + start + openLen
			end := strings.Index(str[start:], close)
			if end < 0 {
				break
			}
			end += start
			buffer.WriteString(str[pos : start-openLen])
			content := str[start:end]
			index++
			newContent, err := replacer(index, start-openLen, end, content)
			if err != nil {
				return "", err
			}
			buffer.WriteString(newContent)
			pos = end + closeLen
		} else {
			break
		}
	}
	buffer.WriteString(str[pos:])
	return buffer.String(), nil
}

func ReplaceByKeyword(str string, keyword byte, replacer ReplacerFunc) (string, error) {
	if str == "" {
		return "", nil
	}
	strLen := len(str)
	index := 0
	start := 0
	end := 0
	buffer := make([]byte, 0, strLen)
	for i := 0; i < strLen; i++ {
		if str[i] == keyword {
			//判断是否最后一位
			if i+1 == strLen {
				return "", fmt.Errorf("Syntax error,near %d '%s'", i, str[i:])
			}
			index++
			start = i + 1
			end = start
		} else {
			if start == 0 {
				buffer = append(buffer, str[i])
			} else {
				char := str[i]
				if (char >= 48 && char <= 57) || (char >= 65 && char <= 90) || (char >= 97 && char <= 122) || char == 95 {
					//判断最后一位
					if i+1 == strLen {
						content := str[start:]
						newContent, err := replacer(index, start, i, content)
						if err != nil {
							return "", err
						}
						buffer = append(buffer, []byte(newContent)...)
					} else {
						end = i
					}
				} else {
					if end > start {
						content := str[start:i]
						newContent, err := replacer(index, start, end, content)
						if err != nil {
							return "", err
						}
						buffer = append(buffer, []byte(newContent)...)
						start = 0
						end = 0
						buffer = append(buffer, str[i])
					} else {
						return "", fmt.Errorf("Syntax error,near %d '%s'", start, str[start-1:])
					}
				}
			}
		}
	}
	return string(buffer), nil
}

func IndexOf(str string, substr string, fromIndex int) int {
	strLen := utf8.RuneCountInString(str)
	if fromIndex >= strLen {
		if substr == "" {
			return strLen
		}
		return -1
	}
	if fromIndex < 0 {
		fromIndex = 0
	}
	if substr == "" {
		return fromIndex
	}
	index := strings.Index(str[fromIndex:], substr)
	if index < 0 {
		return -1
	}
	return fromIndex + index
}
