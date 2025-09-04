package numbers

import (
	"strconv"

	"golang.org/x/exp/constraints"
)

// ToInt64 str to int64，如果转换失败，默认值为0
// str 字符串
func ToInt64(str string) int64 {
	return ToInt64ByDefault(str, 0)
}

// ToInt64ByDefault str to int64
// str 字符串
// def 如果转换失败使用的默认值
func ToInt64ByDefault(str string, def int64) int64 {
	val, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		val = def
	}
	return val
}

// ToInt str to int，如果转换失败，默认值为0
// str 字符串
func ToInt(str string) int {
	return ToIntByDefault(str, 0)
}

// ToIntByDefault str to int
// str 字符串
// def 如果转换失败使用的默认值
func ToIntByDefault(str string, def int) int {
	val, err := strconv.Atoi(str)
	if err != nil {
		val = def
	}
	return val
}

func Min[T constraints.Ordered](a, b T) T {
	if a < b {
		return a
	}
	return b
}

func Max[T constraints.Ordered](a, b T) T {
	if a > b {
		return a
	}
	return b
}
