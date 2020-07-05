package common

import (
	"bbs-go/model/constants"
	"github.com/mlogclub/simple"
	"github.com/mlogclub/simple/markdown"

	"bbs-go/config"
	"math/rand"
	"strconv"
)

// 是否是正式环境
func IsProd() bool {
	return config.Instance.Env == "prod"
}

func GetSummary(contentType string, content string) (summary string) {
	if contentType == constants.ContentTypeMarkdown {
		summary = markdown.GetSummary(content, 256)
	} else if contentType == constants.ContentTypeHtml {
		summary = simple.GetSummary(simple.GetHtmlText(content), 256)
	} else {
		summary = simple.GetSummary(content, 256)
	}
	return
}

// 截取markdown摘要
func GetMarkdownSummary(markdownStr string) string {
	return markdown.GetSummary(markdownStr, 256)
}

// 生成随机验证码
func RandomCode(len int) string {
	if len <= 0 {
		len = 4
	}
	var code string
	for i := 0; i < len; i++ {
		code += strconv.Itoa(rand.Intn(10))
	}
	return code
}
