package utils

import (
	"github.com/mlogclub/mlog/utils/config"
	"strconv"
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

// 首页文章列表
func BuildArticlesUrl(page int) string {
	if page > 1 {
		return BuildAbsUrl("/articles/" + strconv.Itoa(page))
	}
	return BuildAbsUrl("/articles")
}

// 分类文章列表
func BuildCategoryArticlesUrl(categoryId int64, page int) string {
	if page < 1 {
		page = 1
	}
	return BuildAbsUrl("/articles/cat/" + strconv.FormatInt(categoryId, 10) + "/" + strconv.Itoa(page))
}

// 标签文章列表
func BuildTagArticlesUrl(tagId int64, page int) string {
	if page < 1 {
		page = 1
	}
	return BuildAbsUrl("/articles/tag/" + strconv.FormatInt(tagId, 10) + "/" + strconv.Itoa(page))
}

// 用户文章列表
func BuildUserArticlesUrl(userId int64, page int) string {
	path := "/user/" + strconv.FormatInt(userId, 10) + "/articles"
	if page > 1 {
		path = path + "/" + strconv.Itoa(page)
	}
	return BuildAbsUrl(path)
}

// 用户收藏列表
func BuildUserFavoritesUrl(userId int64, page int) string {
	path := "/user/" + strconv.FormatInt(userId, 10) + "/favorites"
	if page > 1 {
		path = path + "/" + strconv.Itoa(page)
	}
	return BuildAbsUrl(path)
}

// 标签列表
func BuildTagsUrl(page int) string {
	path := "/tags"
	if page > 1 {
		path = path + "/" + strconv.Itoa(page)
	}
	return BuildAbsUrl(path)
}

// 消息列表
func BuildMessagesUrl(page int) string {
	path := "/messages"
	if page > 1 {
		path = path + "/" + strconv.Itoa(page)
	}
	return BuildAbsUrl(path)
}

// 话题详情
func BuildTopicUrl(topicId int64) string {
	return BuildAbsUrl("/topic/" + strconv.FormatInt(topicId, 10))
}

// 话题列表
func BuildTopicsUrl(page int) string {
	if page > 1 {
		return BuildAbsUrl("/topics/" + strconv.Itoa(page))
	}
	return BuildAbsUrl("/topics")
}

// 用户帖子列表
func BuildUserTopicsUrl(userId int64, page int) string {
	path := "/user/" + strconv.FormatInt(userId, 10) + "/topics"
	if page > 1 {
		path = path + "/" + strconv.Itoa(page)
	}
	return BuildAbsUrl(path)
}
