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
		return simple.JsonError(simple.ErrorNotLogin)
	}
	title := strings.TrimSpace(simple.FormValue(this.Ctx, "title"))
	content := strings.TrimSpace(simple.FormValue(this.Ctx, "content"))
	tags := simple.FormValueStringArray(this.Ctx, "tags")
	topic, err := services.TopicService.Publish(user.Id, tags, title, content)
	if err != nil {
		return simple.JsonError(err)
	}
	return simple.JsonData(render.BuildSimpleTopic(topic))
}

// 编辑时获取详情
func (this *TopicController) GetEditBy(topicId int64) *simple.JsonResult {
	user := services.UserTokenService.GetCurrent(this.Ctx)
	if user == nil {
		return simple.JsonError(simple.ErrorNotLogin)
	}

	topic := services.TopicService.Get(topicId)
	if topic == nil || topic.Status != model.TopicStatusOk {
		return simple.JsonErrorMsg("话题不存在或已被删除")
	}
	if topic.UserId != user.Id {
		return simple.JsonErrorMsg("无权限")
	}

	tags := services.TopicService.GetTopicTags(topicId)
	var tagNames []string
	if len(tags) > 0 {
		for _, tag := range tags {
			tagNames = append(tagNames, tag.Name)
		}
	}

	return simple.NewEmptyRspBuilder().
		Put("topicId", topic.Id).
		Put("title", topic.Title).
		Put("content", topic.Content).
		Put("tags", tagNames).
		JsonResult()
}

// 编辑帖子
func (this *TopicController) PostEditBy(topicId int64) *simple.JsonResult {
	user := services.UserTokenService.GetCurrent(this.Ctx)
	if user == nil {
		return simple.JsonError(simple.ErrorNotLogin)
	}

	topic := services.TopicService.Get(topicId)
	if topic == nil || topic.Status != model.TopicStatusOk {
		return simple.JsonErrorMsg("话题不存在或已被删除")
	}
	if topic.UserId != user.Id {
		return simple.JsonErrorMsg("无权限")
	}

	title := strings.TrimSpace(simple.FormValue(this.Ctx, "title"))
	content := strings.TrimSpace(simple.FormValue(this.Ctx, "content"))
	tags := simple.FormValueStringArray(this.Ctx, "tags")

	err := services.TopicService.Edit(topicId, tags, title, content)
	if err != nil {
		return simple.JsonError(err)
	}
	return simple.JsonData(render.BuildSimpleTopic(topic))
}

// 删除帖子
func (this *TopicController) PostDeleteBy(topicId int64) *simple.JsonResult {
	user := services.UserTokenService.GetCurrent(this.Ctx)
	if user == nil {
		return simple.JsonError(simple.ErrorNotLogin)
	}
	topic := services.TopicService.Get(topicId)
	if topic == nil || topic.Status != model.TopicStatusOk {
		return simple.JsonSuccess()
	}
	if topic.UserId != user.Id {
		return simple.JsonErrorMsg("无权限")
	}
	err := services.TopicService.Delete(topicId)
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	return simple.JsonSuccess()
}

// 帖子详情
func (this *TopicController) GetBy(topicId int64) *simple.JsonResult {
	topic := services.TopicService.Get(topicId)
	if topic == nil || topic.Status != model.TopicStatusOk {
		return simple.JsonErrorMsg("主题不存在")
	}
	return simple.JsonData(render.BuildTopic(topic))
}

// 最新帖子
func (this *TopicController) GetRecent() *simple.JsonResult {
	topics, err := services.TopicService.QueryCnd(simple.NewQueryCnd("status = ?", model.TopicStatusOk).Order("id desc").Size(10))
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	return simple.JsonData(render.BuildSimpleTopics(topics))
}

// 用户最近的帖子
func (this *TopicController) GetUserRecent() *simple.JsonResult {
	userId, err := simple.FormValueInt64(this.Ctx, "userId")
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	topics, err := services.TopicService.QueryCnd(simple.NewQueryCnd("user_id = ? and status = ?", userId, model.TopicStatusOk).Order("id desc").Size(10))
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	return simple.JsonData(render.BuildSimpleTopics(topics))
}

// 用户帖子列表
func (this *TopicController) GetUserTopics() *simple.JsonResult {
	userId, err := simple.FormValueInt64(this.Ctx, "userId")
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	page := simple.FormValueIntDefault(this.Ctx, "page", 1)

	topics, paging := services.TopicService.Query(simple.NewParamQueries(this.Ctx).
		Eq("user_id", userId).
		Eq("status", model.TopicStatusOk).
		Page(page, 20).Desc("id"))

	return simple.JsonPageData(render.BuildSimpleTopics(topics), paging)
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
		return simple.JsonErrorMsg(err.Error())
	}
	topics, paging := services.TopicService.GetTagTopics(tagId, page)
	return simple.JsonPageData(render.BuildSimpleTopics(topics), paging)
}

// 收藏
func (this *TopicController) GetFavoriteBy(topicId int64) *simple.JsonResult {
	user := services.UserTokenService.GetCurrent(this.Ctx)
	if user == nil {
		return simple.JsonError(simple.ErrorNotLogin)
	}
	err := services.FavoriteService.AddTopicFavorite(user.Id, topicId)
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	return simple.JsonSuccess()
}
