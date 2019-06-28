package utils

import (
	"github.com/mlogclub/mlog/utils/config"
	"github.com/mlogclub/simple"
	"strings"
)

// 是否是正式环境
func IsProd() bool {
	return config.Conf.Env == "prod"
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

// 截取markdown摘要
func GetMarkdownSummary(markdown string) string {
	if len(markdown) == 0 {
		return ""
	}
	markdownResult := simple.Markdown(markdown)
	return markdownResult.SummaryText
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
