package event

import (
	"bbs-go/model"
	"fmt"
	"github.com/mlogclub/simple"
	"github.com/mlogclub/simple/json"
	"reflect"
	"testing"
)

func TestEvent(t *testing.T) {
	//var w sync.WaitGroup
	//w.Add(1)
	RegHandler(reflect.TypeOf(model.User{}), func(i interface{}) {
		fmt.Println("处理用户1")
		fmt.Println(json.ToStr(i))
	})
	RegHandler(reflect.TypeOf(model.User{}), func(i interface{}) {
		fmt.Println("处理用户2")
		fmt.Println(json.ToStr(i))
	})
	Send(model.User{
		Username: simple.SqlNullString("test"),
	})
	//w.Wait()
}
