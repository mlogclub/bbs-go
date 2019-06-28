package main

import (
	"encoding/json"
	"fmt"
	config2 "github.com/mlogclub/mlog/utils/config"
)

func main() {
	// utils.BaiduUrlPush([]string{
	// 	"https://www.mlog.club/share/136",
	// 	"https://www.mlog.club/share/135",
	// })
	// tagIds := strings.Split("", ",")
	// fmt.Println(tagIds)
	// fmt.Println(len(tagIds))
	// fmt.Println(simple.Contains("1", tagIds))

	// var tagIds []string
	// fmt.Print("xxx")
	// fmt.Print(strings.Join(tagIds, ","))
	// fmt.Println("yyy")

	// s := "1,2,,3"
	// ss := strings.Split(s, ",")
	// fmt.Println(len(ss))

	// fmt.Println(strings.Split("管理员", ","))

	config := config2.GetConfig("./mlog.yaml")
	data, _ := json.Marshal(config)
	fmt.Println(string(data))
}
