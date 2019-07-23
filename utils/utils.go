package utils

import (
	"github.com/mlogclub/mlog/utils/config"
	"github.com/mlogclub/simple"
	"regexp"
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
	mdResult := simple.NewMd().Run(markdown)
	return mdResult.SummaryText
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

// 验证用户名合法性，用户名必须由5-12位(数字、字母、_、-)组成，且必须以字母开头。
func IsValidateUsername(username string) bool {
	matched, err := regexp.MatchString("^[0-9a-zA-Z_-]{5,12}$", username)
	if err != nil || !matched {
		return false
	}
	matched, err = regexp.MatchString("^[a-zA-Z]", username)
	if err != nil || !matched {
		return false
	}
	return true
}
