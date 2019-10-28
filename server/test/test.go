package main

import (
	"fmt"

	"github.com/mlogclub/simple"

	"github.com/mlogclub/bbs-go/common/qq"
)

/*
ref: /user/signin
code: 47A71FC84C952C346DC4AFE53113CC38
state: 0049c10c16db4b02a97d271684eb5fe1
*/
func main() {
	// access_token=BC6A85A56265EEA32A12DA61AD8C6154&expires_in=7776000&refresh_token=55CD278593E6519ACEE403C4A0A8AB22
	u, e := qq.GetUserInfo("BC6A85A56265EEA32A12DA61AD8C6154")
	if e != nil {
		fmt.Println(e)
	} else {
		fmt.Println(simple.FormatJson(u))
	}

	// ub := simple.ParseUrl("?access_token=BC6A85A56265EEA32A12DA61AD8C6154&expires_in=7776000&refresh_token=55CD278593E6519ACEE403C4A0A8AB22")
	// fmt.Println(ub.GetQuery().Get("access_token"))

}
