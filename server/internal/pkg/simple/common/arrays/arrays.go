package arrays

import (
	"reflect"
	"strings"
)

func Contains(obj interface{}, arr interface{}) bool {
	targetValue := reflect.ValueOf(arr)
	switch reflect.TypeOf(arr).Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < targetValue.Len(); i++ {
			if targetValue.Index(i).Interface() == obj {
				return true
			}
		}
	case reflect.Map:
		if targetValue.MapIndex(reflect.ValueOf(obj)).IsValid() {
			return true
		}
	}
	return false
}

func ContainsIgnoreCase(str string, arr []string) bool {
	if len(str) == 0 {
		return false
	}
	if len(arr) == 0 {
		return false
	}
	str = strings.ToLower(str)
	for i := 0; i < len(arr); i++ {
		if strings.ToLower(arr[i]) == str {
			return true
		}
	}
	return false
}
