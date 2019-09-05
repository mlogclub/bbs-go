package urls

import (
	"net/url"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/mlogclub/bbs-go/common/config"
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

// 话题详情
func TopicUrl(topicId int64) string {
	return AbsUrl("/topic/" + strconv.FormatInt(topicId, 10))
}

// 项目详情
func ProjectUrl(projectId int64) string {
	return AbsUrl("/project/" + strconv.FormatInt(projectId, 10))
}
