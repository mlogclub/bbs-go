package event

import (
	"bbs-go/model"
	"reflect"
	"testing"
)

func TestEvent(t *testing.T) {
	RegListener(reflect.TypeOf(model.User{}), func(i interface{}) {

	})
	RegListener(reflect.TypeOf(model.User{}), func(i interface{}) {

	})
}
