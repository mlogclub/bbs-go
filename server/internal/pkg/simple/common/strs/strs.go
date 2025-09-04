package strs

import (
	"strings"
	"unicode"

	uuid "github.com/iris-contrib/go.uuid"
)

/*
IsBlank checks if a string is whitespace or empty (""). Observe the following behavior:

	goutils.IsBlank("")        = true
	goutils.IsBlank(" ")       = true
	goutils.IsBlank("bob")     = false
	goutils.IsBlank("  bob  ") = false

Parameter:

	str - the string to check

Returns:

	true - if the string is whitespace or empty ("")
*/
func IsBlank(str string) bool {
	strLen := len(str)
	if str == "" || strLen == 0 {
		return true
	}
	for i := 0; i < strLen; i++ {
		if unicode.IsSpace(rune(str[i])) == false {
			return false
		}
	}
	return true
}

func IsNotBlank(str string) bool {
	return !IsBlank(str)
}

func IsAnyBlank(strs ...string) bool {
	for _, str := range strs {
		if IsBlank(str) {
			return true
		}
	}
	return false
}

func DefaultIfBlank(str, def string) string {
	if IsBlank(str) {
		return def
	} else {
		return str
	}
}

// IsEmpty checks if a string is empty (""). Returns true if empty, and false otherwise.
func IsEmpty(str string) bool {
	return len(str) == 0
}

func IsNotEmpty(str string) bool {
	return !IsEmpty(str)
}

func Substr(s string, start, length int) string {
	bt := []rune(s)
	if start < 0 {
		start = 0
	}
	if start > len(bt) {
		start = start % len(bt)
	}
	var end int
	if (start + length) > (len(bt) - 1) {
		end = len(bt)
	} else {
		end = start + length
	}
	return string(bt[start:end])
}

func Equals(a, b string) bool {
	return a == b
}

func EqualsIgnoreCase(a, b string) bool {
	return a == b || strings.ToUpper(a) == strings.ToUpper(b)
}

func UUID() string {
	u, _ := uuid.NewV4()
	return strings.ReplaceAll(u.String(), "-", "")
}

// RuneLen 字符成长度
func RuneLen(s string) int {
	bt := []rune(s)
	return len(bt)
}

func LeftPad(str string, length int, padStr string) string {
	if length <= len(str) {
		return str
	}
	lenDiff := length - len(str)
	times := (lenDiff + len(padStr) - 1) / len(padStr)
	return strings.Repeat(padStr, times)[:lenDiff] + str
}

func RightPad(str string, length int, padStr string) string {
	if length <= len(str) {
		return str
	}
	lenDiff := length - len(str)
	times := (lenDiff + len(padStr) - 1) / len(padStr)
	return str + strings.Repeat(padStr, times)[:lenDiff]
}
