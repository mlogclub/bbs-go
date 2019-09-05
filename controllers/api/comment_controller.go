package api

import (
	"strconv"

	"github.com/kataras/iris/context"
	"github.com/mlogclub/simple"

	"github.com/mlogclub/bbs-go/controllers/render"
	"github.com/mlogclub/bbs-go/model"
	"github.com/mlogclub/bbs-go/services"
)

type CommentController struct {
	Ctx context.Context
}

func (this *CommentController) GetList() *simple.JsonResult {
	cursor := simple.FormValueInt64Default(this.Ctx, "cursor", 0)

	entityType, err := simple.FormValueRequired(this.Ctx, "entityType")
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	entityId, err := simple.FormValueInt64(this.Ctx, "entityId")
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}

	list, err := services.CommentService.List(entityType, entityId, cursor)
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}

	next := cursor
	var results []model.CommentResponse
	for _, comment := range list {
		results = append(results, *render.BuildComment(comment))
		next = comment.Id
	}
	return simple.JsonCursorData(results, strconv.FormatInt(next, 10))
}

func (this *CommentController) PostCreate() *simple.JsonResult {
	user := services.UserTokenService.GetCurrent(this.Ctx)
	if user == nil {
		return simple.JsonError(simple.ErrorNotLogin)
	}

	form := &model.CreateCommentForm{}
	err := this.Ctx.ReadForm(form)
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}

	comment, err := services.CommentService.Publish(user.Id, form)
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}

	return simple.JsonData(render.BuildComment(*comment))
}
