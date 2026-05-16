package validate

import (
	"bbs-go/internal/pkg/locales"
	"errors"
	"regexp"
	"strings"

	"github.com/mlogclub/simple/common/strs"
)

// IsUsername 验证用户名合法性，用户名必须由5-12位(数字、字母、_、-)组成，且必须以字母开头。
func IsUsername(username string) error {
	if strs.IsBlank(username) {
		return errors.New(locales.Get("user.username_required"))
	}
	matched, err := regexp.MatchString("^[0-9a-zA-Z_-]{5,12}$", username)
	if err != nil || !matched {
		return errors.New(locales.Get("user.username_invalid"))
	}
	matched, err = regexp.MatchString("^[a-zA-Z]", username)
	if err != nil || !matched {
		return errors.New(locales.Get("user.username_invalid"))
	}
	return nil
}

// IsEmail 验证是否是合法的邮箱
func IsEmail(email string) (err error) {
	if strs.IsBlank(email) {
		err = errors.New(locales.Get("user.email_invalid"))
		return
	}
	pattern := `^([A-Za-z0-9_\-\.])+\@([A-Za-z0-9_\-\.])+\.([A-Za-z]{2,4})$`
	matched, _ := regexp.MatchString(pattern, email)
	if !matched {
		err = errors.New(locales.Get("user.email_invalid"))
	}
	return
}

// IsValidPassword 是否是合法的密码
func IsValidPassword(password, rePassword string) error {
	if err := IsPassword(password); err != nil {
		return err
	}
	if password != rePassword {
		return errors.New(locales.Get("user.password_mismatch"))
	}
	return nil
}

func IsPassword(password string) error {
	if strs.IsBlank(password) {
		return errors.New(locales.Get("user.password_required"))
	}
	if strs.RuneLen(password) < 6 {
		return errors.New(locales.Get("user.password_invalid"))
	}
	if strs.RuneLen(password) > 1024 {
		return errors.New(locales.Get("user.password_too_long"))
	}
	return nil
}

// IsURL 是否是合法的URL
func IsURL(url string) error {
	if strs.IsBlank(url) {
		return errors.New(locales.Get("user.url_invalid"))
	}
	indexOfHttp := strings.Index(url, "http://")
	indexOfHttps := strings.Index(url, "https://")
	if indexOfHttp == 0 || indexOfHttps == 0 {
		return nil
	}
	return errors.New(locales.Get("user.url_invalid"))
}
