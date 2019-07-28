package api

import (
	"github.com/kataras/iris/context"
	"github.com/mlogclub/simple"

	"github.com/mlogclub/mlog/controllers/render"
	"github.com/mlogclub/mlog/model"
	"github.com/mlogclub/mlog/services"
)

type CommentController struct {
	Ctx context.Context
}

func (this *CommentController) GetList() *simple.JsonResult {
	cursor := simple.FormValueInt64Default(this.Ctx, "cursor", 0)

	entityType, err := simple.FormValueRequired(this.Ctx, "entityType")
	if err != nil {
		return simple.ErrorMsg(err.Error())
	}
	entityId, err := simple.FormValueInt64(this.Ctx, "entityId")
	if err != nil {
		return simple.ErrorMsg(err.Error())
	}

	list, err := services.CommentService.List(entityType, entityId, cursor)
	if err != nil {
		return simple.ErrorMsg(err.Error())
	}

	nextCursor := cursor
	var itemList []model.CommentResponse
	for _, comment := range list {
		itemList = append(itemList, *render.BuildComment(comment))
		nextCursor = comment.Id
	}
	return simple.NewEmptyRspBuilder().Put("itemList", itemList).Put("cursor", nextCursor).JsonResult()
}

func (this *CommentController) PostCreate() *simple.JsonResult {
	user := services.UserTokenService.GetCurrent(this.Ctx)
	if user == nil {
		return simple.Error(simple.ErrorNotLogin)
	}

	form := &model.CreateCommentForm{}
	err := this.Ctx.ReadForm(form)
	if err != nil {
		return simple.ErrorMsg(err.Error())
	}

	comment, err := services.CommentService.Publish(user.Id, form)
	if err != nil {
		return simple.ErrorMsg(err.Error())
	}

	return simple.JsonData(render.BuildComment(*comment))
}
