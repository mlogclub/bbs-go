
package admin

import (
	"github.com/mlogclub/bbs-go/model"
	"github.com/mlogclub/bbs-go/services"
	"github.com/mlogclub/simple"
	"github.com/kataras/iris/v12"
	"strconv"
)

type TopicLikeController struct {
	Ctx             iris.Context
}

func (this *TopicLikeController) GetBy(id int64) *simple.JsonResult {
	t := services.TopicLikeService.Get(id)
	if t == nil {
		return simple.JsonErrorMsg("Not found, id=" + strconv.FormatInt(id, 10))
	}
	return simple.JsonData(t)
}

func (this *TopicLikeController) AnyList() *simple.JsonResult {
	list, paging := services.TopicLikeService.FindPageByParams(simple.NewQueryParams(this.Ctx).PageByReq().Desc("id"))
	return simple.JsonData(&simple.PageResult{Results: list, Page: paging})
}

func (this *TopicLikeController) PostCreate() *simple.JsonResult {
	t := &model.TopicLike{}
	err := this.Ctx.ReadForm(t)
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}

	err = services.TopicLikeService.Create(t)
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	return simple.JsonData(t)
}

func (this *TopicLikeController) PostUpdate() *simple.JsonResult {
	id, err := simple.FormValueInt64(this.Ctx, "id")
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	t := services.TopicLikeService.Get(id)
	if t == nil {
		return simple.JsonErrorMsg("entity not found")
	}

	err = this.Ctx.ReadForm(t)
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}

	err = services.TopicLikeService.Update(t)
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	return simple.JsonData(t)
}

