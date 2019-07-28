package api

import (
	"strings"

	"github.com/kataras/iris"
	"github.com/mlogclub/simple"

	"github.com/mlogclub/mlog/controllers/render"
	"github.com/mlogclub/mlog/model"
	"github.com/mlogclub/mlog/services"
)

type TopicController struct {
	Ctx iris.Context
}

// 发表帖子
func (this *TopicController) PostCreate() *simple.JsonResult {
	user := services.UserTokenService.GetCurrent(this.Ctx)
	if user == nil {
		return simple.Error(simple.ErrorNotLogin)
	}
	title := strings.TrimSpace(simple.FormValue(this.Ctx, "title"))
	content := strings.TrimSpace(simple.FormValue(this.Ctx, "content"))
	tags := simple.FormValueStringArray(this.Ctx, "tags")
	topic, err := services.TopicService.Publish(user.Id, tags, title, content)
	if err != nil {
		return simple.Error(err)
	}
	return simple.JsonData(render.BuildTopic(topic))
}

// 帖子详情
func (this *TopicController) GetBy(topicId int64) *simple.JsonResult {
	topic := services.TopicService.Get(topicId)
	if topic == nil || topic.Status != model.TopicStatusOk {
		return simple.ErrorMsg("主题不存在")
	}
	return simple.JsonData(render.BuildTopic(topic))
}

// 最新帖子
func (this *TopicController) GetRecent() *simple.JsonResult {
	topics, err := services.TopicService.QueryCnd(simple.NewQueryCnd("status = ?", model.TopicStatusOk).Order("id desc").Size(20))
	if err != nil {
		return simple.ErrorMsg(err.Error())
	}
	return simple.JsonPageData(render.BuildSimpleTopics(topics), nil)
}

// 帖子列表
func (this *TopicController) GetTopics() *simple.JsonResult {
	page := simple.FormValueIntDefault(this.Ctx, "page", 1)

	topics, paging := services.TopicService.Query(simple.NewParamQueries(this.Ctx).
		Eq("status", model.TopicStatusOk).
		Page(page, 20).Desc("last_comment_time"))

	return simple.JsonPageData(render.BuildSimpleTopics(topics), paging)
}

// 标签帖子列表
func (this *TopicController) GetTagTopics() *simple.JsonResult {
	page := simple.FormValueIntDefault(this.Ctx, "page", 1)
	tagId, err := simple.FormValueInt64(this.Ctx, "tagId")
	if err != nil {
		return simple.ErrorMsg(err.Error())
	}
	topics, paging := services.TopicService.GetTagTopics(tagId, page)
	return simple.JsonPageData(render.BuildSimpleTopics(topics), paging)
}

// 收藏
func (this *TopicController) GetFavoriteBy(topicId int64) *simple.JsonResult {
	user := services.UserTokenService.GetCurrent(this.Ctx)
	if user == nil {
		return simple.Error(simple.ErrorNotLogin)
	}
	err := services.FavoriteService.AddTopicFavorite(user.Id, topicId)
	if err != nil {
		return simple.ErrorMsg(err.Error())
	}
	return simple.Success()
}
