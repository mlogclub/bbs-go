package api

import (
	"github.com/kataras/iris"
	"github.com/mlogclub/simple"

	"github.com/mlogclub/mlog/controllers/render"
	"github.com/mlogclub/mlog/model"
	"github.com/mlogclub/mlog/services"
)

type TopicController struct {
	Ctx iris.Context
}

func (this *TopicController) GetBy(topicId int64) *simple.JsonResult {
	topic := services.TopicService.Get(topicId)
	if topic == nil || topic.Status != model.TopicStatusOk {
		return simple.ErrorMsg("主题不存在")
	}
	return simple.JsonData(render.BuildTopic(topic))
}

func (this *TopicController) GetList() *simple.JsonResult {
	return simple.ErrorMsg("unsupported")
}
