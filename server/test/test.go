package main

import (
	"github.com/mlogclub/bbs-go/common/config"
	"github.com/mlogclub/bbs-go/common/email"
)

func init() {
	config.InitConfig("./bbs-go.yaml")
}

func main() {
	email.SendEmail("gaoyoubo@foxmail.com", "企业邮箱测试", "<b>Hello world3</b>")
}
