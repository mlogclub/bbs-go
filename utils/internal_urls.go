package utils

import (
	"net/url"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/mlogclub/mlog/utils/config"
)

func IsInternalUrl(rawUrl string) bool {
	u, err := url.Parse(config.Conf.BaseUrl)
	if err != nil {
		logrus.Error(err)
		return false
	}
	if strings.Index(rawUrl, "#") == 0 { // 如果是“#”开头，说明是锚链接
		return true
	}
	return strings.Contains(rawUrl, u.Host)
}

func BuildAbsUrl(path string) string {
	return config.Conf.BaseUrl + path
}

// 用户主页
func BuildUserUrl(userId int64) string {
	return BuildAbsUrl("/user/" + strconv.FormatInt(userId, 10))
}

// 文章详情
func BuildArticleUrl(articleId int64) string {
	return BuildAbsUrl("/article/" + strconv.FormatInt(articleId, 10))
}

// 话题详情
func BuildTopicUrl(topicId int64) string {
	return BuildAbsUrl("/topic/" + strconv.FormatInt(topicId, 10))
}
