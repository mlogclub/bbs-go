package validate

import (
	"errors"
	"regexp"
	"strings"

	"github.com/mlogclub/simple/common/strs"
)

// IsUsername 验证用户名合法性，用户名必须由5-12位(数字、字母、_、-)组成，且必须以字母开头。
func IsUsername(username string) error {
	if strs.IsBlank(username) {
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

// IsEmail 验证是否是合法的邮箱
func IsEmail(email string) (err error) {
	if strs.IsBlank(email) {
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

// IsValidPassword 是否是合法的密码
func IsValidPassword(password, rePassword string) error {
	if err := IsPassword(password); err != nil {
		return err
	}
	if password != rePassword {
		return errors.New("两次输入密码不匹配")
	}
	return nil
}

func IsPassword(password string) error {
	if strs.IsBlank(password) {
		return errors.New("请输入密码")
	}
	if strs.RuneLen(password) < 6 {
		return errors.New("密码过于简单")
	}
	if strs.RuneLen(password) > 1024 {
		return errors.New("密码长度不能超过128")
	}
	return nil
}

// IsURL 是否是合法的URL
func IsURL(url string) error {
	if strs.IsBlank(url) {
		return errors.New("URL格式错误")
	}
	indexOfHttp := strings.Index(url, "http://")
	indexOfHttps := strings.Index(url, "https://")
	if indexOfHttp == 0 || indexOfHttps == 0 {
		return nil
	}
	return errors.New("URL格式错误")
}
