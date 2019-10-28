package main

import (
	"fmt"
	"strings"
)

/*
ref: /user/signin
code: 47A71FC84C952C346DC4AFE53113CC38
state: 0049c10c16db4b02a97d271684eb5fe1
*/
func main() {
	// access_token=BC6A85A56265EEA32A12DA61AD8C6154&expires_in=7776000&refresh_token=55CD278593E6519ACEE403C4A0A8AB22

	// u, e := qq.GetUserInfo("BC6A85A56265EEA32A12DA61AD8C6154")
	// if e != nil {
	// 	fmt.Println(e)
	// } else {
	// 	fmt.Println(simple.FormatJson(u))
	// }

	// ub := simple.ParseUrl("?access_token=BC6A85A56265EEA32A12DA61AD8C6154&expires_in=7776000&refresh_token=55CD278593E6519ACEE403C4A0A8AB22")
	// fmt.Println(ub.GetQuery().Get("access_token"))

	content:= `callback( {"client_id":"101809474","openid":"C10088323A8652C2D404AD1404027F07","unionid":"UID_39164DF896391520A444724916200E24"} );
`


	prefix := "callback("
	suffix := ");"
	content = strings.TrimSpace(content)
	if strings.Index(content, "callback(") == 0 {
		content = content[len(prefix) : len(content)-len(suffix)]
		content = strings.TrimSpace(content)
	}

	fmt.Println(content)
}
