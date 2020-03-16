package util

import (
	"bytes"
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/zfd81/rooster/types/container"
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
	stack := container.NewArrayStack()
	if str == "" {
		return "", nil
	}
	strLen := len(str) //字符串长度
	index := 0         //关键字出现的顺序
	start := 0         //替换内容的开始位置
	end := 0           //替换内容的结束位置
	buffer := make([]byte, 0, strLen)
	for i := 0; i < strLen; i++ {
		char := str[i]
		if char == '\\' { //转义符判断
			if start > 0 {
				return "", fmt.Errorf("Syntax error,near %d '%s'", start, str[start-1:])
			}
			if !stack.Empty() { //判断转义符前是否有转义符
				stack.Pop()
				buffer = append(buffer, '\\')
			}
			if i+1 == strLen { //判断最后一位
				buffer = append(buffer, char)
			} else {
				stack.Push(char)
			}
		} else if char == keyword { //关键字判断
			//判断是否最后一位
			if i+1 == strLen {
				return "", fmt.Errorf("Syntax error,near %d '%s'", i, str[i:])
			}
			if stack.Empty() {
				index++
				start = i + 1
				end = start
			} else { //关键字被转义
				stack.Pop()
				buffer = append(buffer, keyword)
			}
		} else {
			if !stack.Empty() {
				stack.Pop()
				buffer = append(buffer, '\\')
			}
			if start == 0 { //处理非替换内容
				buffer = append(buffer, char)
			} else { //处理替换内容
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
						buffer = append(buffer, char)
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

func ToUnderscore(str string) string {
	strLen := len(str)
	buffer := make([]byte, 0, strLen+5)
	flag := true
	for i := 0; i < strLen; i++ {
		var b byte = str[i]
		if i > 0 && b >= 'A' && b <= 'Z' {
			if str[i-1] != '_' && flag {
				buffer = append(buffer, '_')
				b += 32
				flag = false
			}
			buffer = append(buffer, b)
		} else {
			buffer = append(buffer, b)
			flag = true
		}
	}
	return string(buffer)
}

func ToCamelCase(str string) string {
	strLen := len(str)
	buffer := make([]byte, 0, strLen)
	flag := false
	limit := strLen - 1
	for i := 0; i < strLen; i++ {
		var b byte = str[i]
		if b == '_' && i > 0 && i < limit && !flag { //第一位和最后一位的"_"不做处理
			flag = true
		} else {
			if flag == true || i == 0 {
				if bool(b >= 'a' && b <= 'z') { //首字母大写
					b -= 32
					flag = false
				}
			}
			buffer = append(buffer, b)
		}
	}
	return string(buffer)
}
