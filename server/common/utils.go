package common

import (
	"bbs-go/model/constants"
	"errors"
	"regexp"
	"strings"

	"github.com/mlogclub/simple"
	"github.com/mlogclub/simple/markdown"

	"bbs-go/config"
)

// 是否是正式环境
func IsProd() bool {
	return config.Instance.Env == "prod"
}

// index of
func IndexOf(userIds []int64, userId int64) int {
	if userIds == nil || len(userIds) == 0 {
		return -1
	}
	for i, v := range userIds {
		if v == userId {
			return i
		}
	}
	return -1
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

// 获取html内容摘要
func GetHtmlSummary(html string) string {
	if len(html) == 0 {
		return ""
	}
	text := simple.GetHtmlText(html)
	return simple.GetSummary(text, 256)
}

// 获取用户角色
func GetUserRoles(roles string) []string {
	if len(roles) == 0 {
		return nil
	}
	ss := strings.Split(roles, ",")
	if len(ss) == 0 {
		return nil
	}
	var ret []string
	for _, s := range ss {
		s = strings.TrimSpace(s)
		if len(s) > 0 {
			ret = append(ret, s)
		}
	}
	return ret
}

// 是否是内部图片
func IsInternalImage(imageUrl string) bool {
	// TODO @ 2019/12/31 这个地方硬编码了，要修改
	return strings.Contains(imageUrl, "file.mlog.club") || strings.Contains(imageUrl, "static.mlog.club")
}

// 应用图片样式
func ApplyImageStyle(imageUrl, styleName string) string {
	if !IsInternalImage(imageUrl) {
		return imageUrl
	}
	splitterIndex := strings.LastIndex(imageUrl, "!")
	if splitterIndex >= 0 {
		imageUrl = imageUrl[:splitterIndex]
	}
	return imageUrl + "!" + styleName
}
