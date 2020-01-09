package simple

import (
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	uuid "github.com/iris-contrib/go.uuid"
)

// 截取字符串
func Substr(s string, start, length int) string {
	bt := []rune(s)
	if start < 0 {
		start = 0
	}
	if start > len(bt) {
		start = start % len(bt)
	}
	var end int
	if (start + length) > (len(bt) - 1) {
		end = len(bt)
	} else {
		end = start + length
	}
	return string(bt[start:end])
}

// uuid
func Uuid() string {
	u, _ := uuid.NewV4()
	s := u.String()
	s = strings.ReplaceAll(s, "-", "")
	return s
}

// 字符成长度
func RuneLen(s string) int {
	bt := []rune(s)
	return len(bt)
}

// 获取summary
func GetSummary(s string, length int) string {
	summary := Substr(s, 0, length)
	if RuneLen(s) > length {
		summary += "..."
	}
	return summary
}

// 获取html文本
func GetHtmlText(html string) string {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return ""
	}
	return doc.Text()
}

// AbsoluteURL
// u : /1/2/3 ...
// baseURL : www.baidu.com
// Url : www.baidu.com/a/b/c
func AbsoluteURL(u string, baseURL, URL string) (string, error) {
	if strings.HasPrefix(u, "#") {
		return "", nil
	}

	var (
		_baseURL *url.URL
		_URL     *url.URL
	)

	_baseURL, err := url.Parse(baseURL)
	if err != nil {
		return "", err
	}

	if len(URL) > 0 {
		_URL, err = url.Parse(URL)
		if err != nil {
			return "", err
		}
	}

	var _base *url.URL
	if _baseURL != nil {
		_base = _baseURL
	} else {
		_base = _URL
	}

	if _base == nil {
		return u, nil
	}

	absURL, err := _base.Parse(u)
	if err != nil {
		return "", err
	}
	absURL.Fragment = ""
	if absURL.Scheme == "//" {
		absURL.Scheme = _URL.Scheme
	}
	return absURL.String(), nil
}
