package admin

import (
	"github.com/kataras/iris"
	"github.com/mlogclub/mlog/controllers/render"
	"github.com/mlogclub/mlog/model"
	"github.com/mlogclub/mlog/services"
	"github.com/mlogclub/simple"
	"strconv"
)

type TopicController struct {
	Ctx iris.Context
}

func (this *TopicController) GetBy(id int64) *simple.JsonResult {
	t := services.TopicService.Get(id)
	if t == nil {
		return simple.ErrorMsg("Not found, id=" + strconv.FormatInt(id, 10))
	}
	return simple.JsonData(t)
}

func (this *TopicController) AnyList() *simple.JsonResult {
	list, paging := services.TopicService.Query(simple.NewParamQueries(this.Ctx).
		EqAuto("status").LikeAuto("title").PageAuto().Desc("id"))

	var results []map[string]interface{}
	for _, topic := range list {
		builder := simple.NewRspBuilderExcludes(topic, "content")

		// 用户
		builder = builder.Put("user", render.BuildUserDefaultIfNull(topic.UserId))

		// 简介
		mr := simple.NewMd().Run(topic.Content)
		builder.Put("summary", mr.SummaryText)

		// 标签
		tags := services.TopicService.GetTopicTags(topic.Id)
		builder.Put("tags", render.BuildTags(tags))

		results = append(results, builder.Build())
	}

	return simple.JsonData(&simple.PageResult{Results: results, Page: paging})

}

func (this *TopicController) PostCreate() *simple.JsonResult {
	t := &model.Topic{}
	err := this.Ctx.ReadForm(t)
	if err != nil {
		return simple.ErrorMsg(err.Error())
	}

	err = services.TopicService.Create(t)
	if err != nil {
		return simple.ErrorMsg(err.Error())
	}
	return simple.JsonData(t)
}

func (this *TopicController) PostUpdate() *simple.JsonResult {
	id, err := simple.FormValueInt64(this.Ctx, "id")
	if err != nil {
		return simple.ErrorMsg(err.Error())
	}
	t := services.TopicService.Get(id)
	if t == nil {
		return simple.ErrorMsg("entity not found")
	}

	err = this.Ctx.ReadForm(t)
	if err != nil {
		return simple.ErrorMsg(err.Error())
	}

	err = services.TopicService.Update(t)
	if err != nil {
		return simple.ErrorMsg(err.Error())
	}
	return simple.JsonData(t)
}

func (this *TopicController) PostDelete() *simple.JsonResult {
	id, err := simple.FormValueInt64(this.Ctx, "id")
	if err != nil {
		return simple.ErrorMsg(err.Error())
	}

	err = services.TopicService.Delete(id)
	if err != nil {
		return simple.ErrorMsg(err.Error())
	}
	return simple.Success()
}
