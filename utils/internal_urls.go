package utils

import (
	"strconv"

	"github.com/mlogclub/mlog/utils/config"
)

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
