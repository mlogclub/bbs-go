package admin

import (
	"strconv"

	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple"

	"bbs-go/controllers/render"
	"bbs-go/model"
	"bbs-go/services"
)

type TopicController struct {
	Ctx iris.Context
}

func (c *TopicController) GetBy(id int64) *simple.JsonResult {
	t := services.TopicService.Get(id)
	if t == nil {
		return simple.JsonErrorMsg("Not found, id=" + strconv.FormatInt(id, 10))
	}
	return simple.JsonData(t)
}

func (c *TopicController) AnyList() *simple.JsonResult {
	list, paging := services.TopicService.FindPageByParams(simple.NewQueryParams(c.Ctx).
		EqByReq("id").EqByReq("user_id").EqByReq("status").EqByReq("recommend").LikeByReq("title").PageByReq().Desc("id"))

	var results []map[string]interface{}
	for _, topic := range list {
		builder := simple.NewRspBuilderExcludes(topic, "content")

		// 用户
		builder = builder.Put("user", render.BuildUserDefaultIfNull(topic.UserId))

		// 节点
		node := services.TopicNodeService.Get(topic.NodeId)
		builder.Put("node", node)

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

// 推荐
func (c *TopicController) PostRecommend() *simple.JsonResult {
	id, err := simple.FormValueInt64(c.Ctx, "id")
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	err = services.TopicService.SetRecommend(id, true)
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	return simple.JsonSuccess()
}

// 取消推荐
func (c *TopicController) DeleteRecommend() *simple.JsonResult {
	id, err := simple.FormValueInt64(c.Ctx, "id")
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	err = services.TopicService.SetRecommend(id, false)
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	return simple.JsonSuccess()
}

func (c *TopicController) PostCreate() *simple.JsonResult {
	t := &model.Topic{}
	err := c.Ctx.ReadForm(t)
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}

	err = services.TopicService.Create(t)
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	return simple.JsonData(t)
}

func (c *TopicController) PostUpdate() *simple.JsonResult {
	id, err := simple.FormValueInt64(c.Ctx, "id")
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	t := services.TopicService.Get(id)
	if t == nil {
		return simple.JsonErrorMsg("entity not found")
	}

	err = c.Ctx.ReadForm(t)
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}

	err = services.TopicService.Update(t)
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	return simple.JsonData(t)
}

func (c *TopicController) PostDelete() *simple.JsonResult {
	id, err := simple.FormValueInt64(c.Ctx, "id")
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}

	err = services.TopicService.Delete(id)
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	return simple.JsonSuccess()
}

func (c *TopicController) PostUndelete() *simple.JsonResult {
	id, err := simple.FormValueInt64(c.Ctx, "id")
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	err = services.TopicService.Undelete(id)
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	return simple.JsonSuccess()
}
