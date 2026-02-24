package bbsurls

import (
	"bbs-go/internal/cache"
	"bbs-go/internal/models/constants"
	"log/slog"
	"net/url"
	"strconv"
	"strings"

	"bbs-go/internal/pkg/idcodec"
)

// 是否是内部链接
func IsInternalUrl(href string) bool {
	if IsAnchor(href) {
		return true
	}
	u, err := url.Parse(getBaseURL())
	if err != nil {
		slog.Error(err.Error(), slog.Any("err", err))
		return false
	}
	if strings.TrimSpace(u.Host) == "" {
		return strings.HasPrefix(href, "/")
	}
	return strings.Contains(href, u.Host)
}

// 是否是锚链接
func IsAnchor(href string) bool {
	return strings.Index(href, "#") == 0
}

func AbsUrl(path string) string {
	baseURL := getBaseURL()
	if baseURL == "/" {
		return path
	}
	return baseURL + path
}

func getBaseURL() string {
	baseURL := strings.TrimSpace(cache.SysConfigCache.GetStr(constants.SysConfigBaseURL))
	if baseURL == "" {
		return "/"
	}
	for len(baseURL) > 1 && strings.HasSuffix(baseURL, "/") {
		baseURL = strings.TrimSuffix(baseURL, "/")
	}
	return baseURL
}

// 用户主页
func UserUrl(userId int64) string {
	return AbsUrl("/user/" + idcodec.Encode(userId))
}

// 文章详情
func ArticleUrl(articleId int64) string {
	return AbsUrl("/article/" + strconv.FormatInt(articleId, 10))
}

// 标签文章列表
func TagArticlesUrl(tagId int64) string {
	return AbsUrl("/articles/" + strconv.FormatInt(tagId, 10))
}

// 话题详情
func TopicUrl(topicId int64) string {
	return AbsUrl("/topic/" + idcodec.Encode(topicId))
}

func UrlJoin(parts ...string) string {
	sep := "/"
	var ss []string
	for i, part := range parts {
		part = strings.TrimSpace(part)
		var (
			from = 0
			to   = len(part)
		)
		if strings.Index(part, sep) == 0 {
			from = 1
		}
		if strings.LastIndex(part, sep) == len(part)-1 {
			to = len(part) - 1
		}
		part = part[from:to]

		ss = append(ss, part)
		if i != len(parts)-1 {
			ss = append(ss, sep)
		}
	}
	return strings.Join(ss, "")
}
