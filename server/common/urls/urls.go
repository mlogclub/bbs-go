package urls

import (
	"net/url"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/mlogclub/bbs-go/common/config"
)

// 是否是内部链接
func IsInternalUrl(href string) bool {
	if IsAnchor(href) {
		return true
	}
	u, err := url.Parse(config.Conf.BaseUrl)
	if err != nil {
		logrus.Error(err)
		return false
	}
	return strings.Contains(href, u.Host)
}

// 是否是锚链接
func IsAnchor(href string) bool {
	return strings.Index(href, "#") == 0
}

func AbsUrl(path string) string {
	return config.Conf.BaseUrl + path
}

// 用户主页
func UserUrl(userId int64) string {
	return AbsUrl("/user/" + strconv.FormatInt(userId, 10))
}

// 文章详情
func ArticleUrl(articleId int64) string {
	return AbsUrl("/article/" + strconv.FormatInt(articleId, 10))
}

// 标签文章列表
func TagArticlesUrl(tagId int64) string {
	return AbsUrl("/articles/tag/" + strconv.FormatInt(tagId, 10))
}

// 话题详情
func TopicUrl(topicId int64) string {
	return AbsUrl("/topic/" + strconv.FormatInt(topicId, 10))
}

// 项目详情
func ProjectUrl(projectId int64) string {
	return AbsUrl("/project/" + strconv.FormatInt(projectId, 10))
}
