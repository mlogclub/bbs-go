package async_test

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"bbs-go/internal/pkg/simple/common/async"
)

func TestAsync(t *testing.T) {
	f1 := async.Exec(func() (int, error) {
		fmt.Println("执行方法1")
		return 1, nil
	})
	f2 := async.Exec(func() (int, error) {
		fmt.Println("执行方法2")
		time.Sleep(1 * time.Second)
		return 2, errors.New("失败")
	})

	fmt.Println(f1.Await())
	fmt.Println(f2.Await())
}

func TestAsyncWithNotTimeout(t *testing.T) {
	f := async.Exec(func() (string, error) {
		fmt.Println("执行方法")
		time.Sleep(2 * time.Second)
		return "成功", nil
	})

	fmt.Println(f.AwaitTimeout(3 * time.Second))
}

func TestAsyncWithTimeout(t *testing.T) {
	f := async.Exec(func() (string, error) {
		fmt.Println("执行方法")
		time.Sleep(5 * time.Second)
		return "成功", nil
	})
	fmt.Println(f.AwaitTimeout(3 * time.Second))
}
