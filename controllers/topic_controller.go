package controllers

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/context"
	"github.com/mlogclub/mlog/controllers/render"
	"github.com/mlogclub/mlog/model"
	"github.com/mlogclub/mlog/services"
	"github.com/mlogclub/mlog/utils"
	"github.com/mlogclub/simple"
	"strings"
)

type TopicController struct {
	Ctx             iris.Context
	TopicService    *services.TopicService
	FavoriteService *services.FavoriteService
}

func (this *TopicController) GetBy(topicId int64) {
	topic := this.TopicService.Get(topicId)

	if topic == nil || topic.Status != model.TopicStatusOk {
		this.Ctx.StatusCode(404)
		return
	}

	// 浏览数量+1
	this.TopicService.IncrViewCount(topicId)

	render.View(this.Ctx, "topic/detail.html", iris.Map{
		utils.GlobalFieldSiteTitle: topic.Title,
		"CommentEntityType":        model.EntityTypeTopic,
		"CommentEntityId":          topic.Id,
		"Topic":                    render.BuildTopic(topic),
	})
}

func (this *TopicController) GetCreate() {
	user := utils.GetCurrentUser(this.Ctx)
	if user == nil {
		this.Ctx.Redirect("/user/signin?redirectUrl=/topic/create", iris.StatusSeeOther)
		return
	}
	render.View(this.Ctx, "topic/create.html", iris.Map{
		utils.GlobalFieldSiteTitle: "发起讨论",
	})
	return
}

func (this *TopicController) PostCreate() *simple.JsonResult {
	user := utils.GetCurrentUser(this.Ctx)
	if user == nil {
		return simple.Error(simple.ErrorNotLogin)
	}
	title := strings.TrimSpace(simple.FormValue(this.Ctx, "title"))
	content := strings.TrimSpace(simple.FormValue(this.Ctx, "content"))
	tags := simple.FormValueStringArray(this.Ctx, "tags")

	topic, err := this.TopicService.Publish(user.Id, tags, title, content)
	if err != nil {
		return simple.Error(err)
	}
	return simple.NewEmptyRspBuilder().Put("topicId", topic.Id).JsonResult()
}

// 收藏
func (this *TopicController) PostFavoriteBy(topicId int64) *simple.JsonResult {
	user := utils.GetCurrentUser(this.Ctx)
	if user == nil {
		return simple.Error(simple.ErrorNotLogin)
	}
	err := this.FavoriteService.AddTopicFavorite(user.Id, topicId)
	if err != nil {
		return simple.ErrorMsg(err.Error())
	}
	return simple.Success()
}

// 帖子列表
func GetTopics(ctx context.Context) {
	page := ctx.Params().GetIntDefault("page", 1)

	topics, paging := services.TopicServiceInstance.Query(simple.NewParamQueries(ctx).
		Eq("status", model.TopicStatusOk).
		Page(page, 20).Desc("id"))

	render.View(ctx, "topic/index.html", iris.Map{
		utils.GlobalFieldSiteTitle: "讨论",
		"Topics":                   render.BuildTopics(topics),
		"Page":                     paging,
		"PrePageUrl":               utils.BuildTopicsUrl(page - 1),
		"NextPageUrl":              utils.BuildTopicsUrl(page + 1),
	})
}
