// 实现简单的三元表达式功能
package simple

import (
	"reflect"
	"strconv"
)

// If - (a ? b : c) Or (a && b)
func If(args ...interface{}) interface{} {
	var condition = callFn(args[0])
	if len(args) == 1 {
		return condition
	}
	var trueVal = args[1]
	var falseVal interface{}
	if len(args) > 2 {
		falseVal = args[2]
	} else {
		falseVal = nil
	}
	if condition == nil {
		return callFn(falseVal)
	} else if v, ok := condition.(bool); ok {
		if v == false {
			return callFn(falseVal)
		}
	} else if isFalse(condition) {
		return callFn(falseVal)
	} else if v, ok := condition.(error); ok {
		if v != nil {
			return condition
		}
	}
	return callFn(trueVal)
}

func isFalse(f interface{}) bool {
	v := reflect.ValueOf(f)
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.String:
		str := v.String()
		if str == "" {
			return true
		}
		zero, err := strconv.ParseFloat(str, 10)
		if zero == 0 && err == nil {
			return true
		}
		boolean, err := strconv.ParseBool(str)
		return boolean == false && err == nil
	default:
		return false
	}
}

// callFn if args[i] == func, run it
func callFn(f interface{}) interface{} {
	if f != nil {
		t := reflect.TypeOf(f)
		if t.Kind() == reflect.Func && t.NumIn() == 0 {
			function := reflect.ValueOf(f)
			in := make([]reflect.Value, 0)
			out := function.Call(in)
			if num := len(out); num > 0 {
				list := make([]interface{}, num)
				for i, value := range out {
					list[i] = value.Interface()
				}
				if num == 1 {
					return list[0]
				}
				return list
			}
			return nil
		}
	}
	return f
}
