package event

import (
	"fmt"
	"reflect"
	"testing"

	"bbs-go/internal/models"

	"bbs-go/internal/pkg/simple/common/jsons"
	"bbs-go/internal/pkg/simple/sqls"
)

func TestEvent(t *testing.T) {
	//var w sync.WaitGroup
	//w.Add(1)
	RegHandler(reflect.TypeOf(models.User{}), func(i interface{}) {
		fmt.Println("处理用户1")
		fmt.Println(jsons.ToStr(i))
	})
	RegHandler(reflect.TypeOf(models.User{}), func(i interface{}) {
		fmt.Println("处理用户2")
		fmt.Println(jsons.ToStr(i))
	})
	Send(models.User{
		Username: sqls.SqlNullString("test"),
	})
	//w.Wait()
}
