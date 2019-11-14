package api

import (
	"strconv"

	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple"

	"github.com/mlogclub/bbs-go/controllers/render"
	"github.com/mlogclub/bbs-go/model"
	"github.com/mlogclub/bbs-go/services"
)

type CommentController struct {
	Ctx iris.Context
}

func (this *CommentController) GetList() *simple.JsonResult {
	var (
		err        error
		cursor     int64
		entityType string
		entityId   int64
	)
	cursor = simple.FormValueInt64Default(this.Ctx, "cursor", 0)

	if entityType, err = simple.FormValueRequired(this.Ctx, "entityType"); err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	if entityId, err = simple.FormValueInt64(this.Ctx, "entityId"); err != nil {
		return simple.JsonErrorMsg(err.Error())
	}

	comments, cursor := services.CommentService.GetComments(entityType, entityId, cursor)
	return simple.JsonCursorData(render.BuildComments(comments), strconv.FormatInt(cursor, 10))
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
