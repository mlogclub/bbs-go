package common

import (
	"errors"
	"regexp"
	"strings"

	"github.com/mlogclub/simple"

	"github.com/mlogclub/bbs-go/common/config"
	"github.com/mlogclub/bbs-go/model"
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

func GetSummary(contentType string, content string) (summary string) {
	if contentType == model.ContentTypeMarkdown {
		mr := simple.NewMd().Run(content)
		summary = mr.SummaryText
	} else {
		summary = simple.GetSummary(simple.GetHtmlText(content), 256)
	}
	return
}

// 截取markdown摘要
func GetMarkdownSummary(markdown string) string {
	if len(markdown) == 0 {
		return ""
	}
	mdResult := simple.NewMd().Run(markdown)
	return mdResult.SummaryText
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

// 验证用户名合法性，用户名必须由5-12位(数字、字母、_、-)组成，且必须以字母开头。
func IsValidateUsername(username string) error {
	if len(username) == 0 {
		return errors.New("请输入用户名")
	}
	matched, err := regexp.MatchString("^[0-9a-zA-Z_-]{5,12}$", username)
	if err != nil || !matched {
		return errors.New("用户名必须由5-12位(数字、字母、_、-)组成，且必须以字母开头")
	}
	matched, err = regexp.MatchString("^[a-zA-Z]", username)
	if err != nil || !matched {
		return errors.New("用户名必须由5-12位(数字、字母、_、-)组成，且必须以字母开头")
	}
	return nil
}

// 验证是否是合法的邮箱
func IsValidateEmail(email string) (err error) {
	if len(email) == 0 {
		err = errors.New("邮箱格式不符合规范")
		return
	}
	pattern := `^([A-Za-z0-9_\-\.])+\@([A-Za-z0-9_\-\.])+\.([A-Za-z]{2,4})$`
	matched, _ := regexp.MatchString(pattern, email)
	if !matched {
		err = errors.New("邮箱格式不符合规范")
	}
	return
}

// 是否是合法的密码
func IsValidatePassword(password, rePassword string) error {
	if len(password) == 0 {
		return errors.New("请输入密码")
	}
	if simple.RuneLen(password) < 6 {
		return errors.New("密码过于简单")
	}
	if password != rePassword {
		return errors.New("两次输入密码不匹配")
	}
	return nil
}
