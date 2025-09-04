package jsons

import (
	"fmt"
	"testing"
)

func TestToJsonStr(t *testing.T) {
	ret := ToJsonStr(nil)
	fmt.Println(ret)
	fmt.Println(ret == "")
	if ret != "" {
		t.Fail()
	}
}
