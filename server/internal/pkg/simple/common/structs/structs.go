package structs

import (
	"errors"
	"fmt"
	"log"
	"reflect"

	"bbs-go/internal/pkg/simple/common/arrays"
)

func StructToMap(obj interface{}, excludes ...string) map[string]interface{} {
	var data = make(map[string]interface{})
	keys := reflect.TypeOf(obj)
	values := reflect.ValueOf(obj)
	fillMap(data, keys, values, excludes...)
	return data
}

func fillMap(data map[string]interface{}, keys reflect.Type, values reflect.Value, excludes ...string) {
	if values.Kind() == reflect.Ptr {
		values = values.Elem()
	}
	if keys.Kind() == reflect.Ptr {
		keys = keys.Elem()
	}

	for i := 0; i < keys.NumField(); i++ {
		keyField := keys.Field(i)
		valueField := values.Field(i)

		if keyField.Anonymous {
			fillMap(data, keyField.Type, valueField, excludes...)
		} else {
			if !arrays.ContainsIgnoreCase(keyField.Name, excludes) {
				jsonTag := keyField.Tag.Get("json")
				if len(jsonTag) > 0 {
					data[jsonTag] = valueField.Interface()
				} else {
					data[keyField.Name] = valueField.Interface()
				}
			}
		}
	}
}

func MapToStruct(obj interface{}, data map[string]interface{}) error {
	for k, v := range data {
		err := setField(obj, k, v)
		if err != nil {
			return err
		}
	}
	return nil
}

func setField(obj interface{}, name string, value interface{}) error {
	structValue := reflect.ValueOf(obj).Elem()
	structFieldValue := structValue.FieldByName(name)

	if !structFieldValue.IsValid() {
		return fmt.Errorf("No such field: %s in obj ", name)
	}

	if !structFieldValue.CanSet() {
		return fmt.Errorf("Cannot set %s field value ", name)
	}

	structFieldType := structFieldValue.Type()
	val := reflect.ValueOf(value)
	if structFieldType != val.Type() {
		return errors.New("Provided value type didn't match obj field type ")
	}
	structFieldValue.Set(val)
	return nil
}

func StructFields(s interface{}) []reflect.StructField {
	t := StructTypeOf(s)
	if t.Kind() != reflect.Struct {
		log.Println("Check type error not Struct")
		return nil
	}

	var results []reflect.StructField
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		results = append(results, f)
		// if f.Anonymous {
		// 	fields := StructFields(f.Type)
		// 	results = append(results, fields...)
		// }
	}
	return results
}

func StructName(s interface{}) string {
	t := StructTypeOf(s)
	return t.Name()
}

func StructTypeOf(s interface{}) reflect.Type {
	t := reflect.TypeOf(s)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t
}
