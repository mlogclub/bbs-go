package event

import (
	"bbs-go/model"
	"fmt"
	"reflect"
	"testing"

	"github.com/mlogclub/simple/common/jsons"
	"github.com/mlogclub/simple/sqls"
)

func TestEvent(t *testing.T) {
	//var w sync.WaitGroup
	//w.Add(1)
	RegHandler(reflect.TypeOf(model.User{}), func(i interface{}) {
		fmt.Println("处理用户1")
		fmt.Println(jsons.ToStr(i))
	})
	RegHandler(reflect.TypeOf(model.User{}), func(i interface{}) {
		fmt.Println("处理用户2")
		fmt.Println(jsons.ToStr(i))
	})
	Send(model.User{
		Username: sqls.SqlNullString("test"),
	})
	//w.Wait()
}
