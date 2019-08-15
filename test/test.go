package main

import (
	"fmt"

	"github.com/mlogclub/mlog/utils"
)

func main() {
	url := "https://file.mlog.club/images/2019/08/07/ff14571273b239543a1ce6b4de48d85d.jpg"
	aliyunOss, _ := utils.NewAliyunOss("oss-cn-hongkong.aliyuncs.com", "LTAI4UX8o6kXhXs8", "SDxyWC9avQi1hugOSIW33I9LdKHMhO", "mlogclub", "https://file.mlog.club/")
	ret := aliyunOss.SignUrl(url)
	fmt.Println(ret)
}
