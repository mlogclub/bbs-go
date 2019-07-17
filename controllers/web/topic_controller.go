package web

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/context"
	"github.com/mlogclub/mlog/controllers/render"
	"github.com/mlogclub/mlog/model"
	"github.com/mlogclub/mlog/services"
	"github.com/mlogclub/mlog/utils"
	"github.com/mlogclub/mlog/utils/session"
	"github.com/mlogclub/simple"
	"strings"
)

type TopicController struct {
	Ctx iris.Context
}

func (this *TopicController) GetBy(topicId int64) {
	topic := services.TopicService.Get(topicId)

	if topic == nil || topic.Status != model.TopicStatusOk {
		this.Ctx.StatusCode(404)
		return
	}

	// 浏览数量+1
	services.TopicService.IncrViewCount(topicId)

	render.View(this.Ctx, "topic/detail.html", iris.Map{
		model.TplSiteTitle:  topic.Title,
		"CommentEntityType": model.EntityTypeTopic,
		"CommentEntityId":   topic.Id,
		"Topic":             render.BuildTopic(topic),
	})
}

func (this *TopicController) GetCreate() {
	user := session.GetCurrentUser(this.Ctx)
	if user == nil {
		this.Ctx.Redirect("/user/signin?redirectUrl=/topic/create", iris.StatusSeeOther)
		return
	}
	render.View(this.Ctx, "topic/create.html", iris.Map{
		model.TplSiteTitle: "发表主题",
	})
	return
}

func (this *TopicController) PostCreate() *simple.JsonResult {
	user := session.GetCurrentUser(this.Ctx)
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
	return simple.NewEmptyRspBuilder().Put("topicId", topic.Id).JsonResult()
}

func (this *TopicController) GetEditBy(topicId int64) {
	user := session.GetCurrentUser(this.Ctx)
	if user == nil {
		this.Ctx.StatusCode(403)
		return
	}

	topic := services.TopicService.Get(topicId)
	if topic == nil || topic.Status != model.TopicStatusOk || topic.UserId != user.Id {
		this.Ctx.StatusCode(404)
		return
	}

	tags := services.TopicService.GetTopicTags(topicId)
	var tagNames []string
	if len(tags) > 0 {
		for _, tag := range tags {
			tagNames = append(tagNames, tag.Name)
		}
	}

	render.View(this.Ctx, "topic/edit.html", iris.Map{
		"Topic": iris.Map{
			"TopicId": topic.Id,
			"Title":   topic.Title,
			"Content": topic.Content,
			"Tags":    tagNames,
		},
	})
}

func (this *TopicController) PostEditBy(topicId int64) *simple.JsonResult {
	user := session.GetCurrentUser(this.Ctx)
	if user == nil {
		return simple.Error(simple.ErrorNotLogin)
	}

	topic := services.TopicService.Get(topicId)
	if topic == nil || topic.Status != model.TopicStatusOk || topic.UserId != user.Id {
		return simple.ErrorMsg("话题不存在或已被删除")
	}

	title := strings.TrimSpace(simple.FormValue(this.Ctx, "title"))
	content := strings.TrimSpace(simple.FormValue(this.Ctx, "content"))
	tags := simple.FormValueStringArray(this.Ctx, "tags")

	err := services.TopicService.Edit(topicId, tags, title, content)
	if err != nil {
		return simple.Error(err)
	}
	return simple.NewEmptyRspBuilder().Put("topicId", topic.Id).JsonResult()
}

// 收藏
func (this *TopicController) PostFavoriteBy(topicId int64) *simple.JsonResult {
	user := session.GetCurrentUser(this.Ctx)
	if user == nil {
		return simple.Error(simple.ErrorNotLogin)
	}
	err := services.FavoriteService.AddTopicFavorite(user.Id, topicId)
	if err != nil {
		return simple.ErrorMsg(err.Error())
	}
	return simple.Success()
}

// 主题列表
func GetTopics(ctx context.Context) {
	page := ctx.Params().GetIntDefault("page", 1)

	topics, paging := services.TopicService.Query(simple.NewParamQueries(ctx).
		Eq("status", model.TopicStatusOk).
		Page(page, 20).Desc("last_comment_time"))

	render.View(ctx, "topic/index.html", iris.Map{
		model.TplSiteTitle: "主题",
		"Topics":           render.BuildTopics(topics),
		"Page":             paging,
		"PrePageUrl":       utils.BuildTopicsUrl(page - 1),
		"NextPageUrl":      utils.BuildTopicsUrl(page + 1),
	})
}

// 标签主题列表
func GetTagTopics(ctx iris.Context) {
	tagId := ctx.Params().GetInt64Default("tagId", 0)
	page := ctx.Params().GetIntDefault("page", 1)

	tag := services.TagService.Get(tagId)

	title := "主题"
	if tag != nil {
		title = tag.Name + " - " + title
	}

	topics, paging := services.TopicService.GetTagTopics(tagId, page)

	render.View(ctx, "topic/index.html", iris.Map{
		model.TplSiteTitle: title,
		"Topics":           render.BuildTopics(topics),
		"Page":             paging,
		"PrePageUrl":       utils.BuildTagTopicsUrl(tagId, page-1),
		"NextPageUrl":      utils.BuildTagTopicsUrl(tagId, page+1),
	})
}
