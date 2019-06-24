package utils

import (
	"github.com/mlogclub/simple"
	"strconv"
	"strings"
)

// 是否是正式环境
func IsProd() bool {
	return Conf.Env == "prod"
}

func BuildAbsUrl(path string) string {
	return Conf.BaseUrl + path
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

// 用户话题列表
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

// 截取markdown摘要
func GetMarkdownSummary(markdown string) string {
	if len(markdown) == 0 {
		return ""
	}
	markdownResult := simple.Markdown(markdown)
	return markdownResult.SummaryText
}

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
