package api

import (
	"bbs-go/model"
	"bbs-go/spam"
	"strconv"

	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple"
	"github.com/mlogclub/simple/mvc"
	"github.com/mlogclub/simple/mvc/params"

	"bbs-go/controllers/render"
	"bbs-go/services"
)

type CommentController struct {
	Ctx iris.Context
}

func (c *CommentController) GetComments() *mvc.JsonResult {
	var (
		err        error
		cursor     int64
		entityType string
		entityId   int64
	)
	cursor = params.FormValueInt64Default(c.Ctx, "cursor", 0)

	if entityType, err = params.FormValueRequired(c.Ctx, "entityType"); err != nil {
		return mvc.JsonErrorMsg(err.Error())
	}
	if entityId, err = params.FormValueInt64(c.Ctx, "entityId"); err != nil {
		return mvc.JsonErrorMsg(err.Error())
	}
	currentUser := services.UserTokenService.GetCurrent(c.Ctx)
	comments, cursor, hasMore := services.CommentService.GetComments(entityType, entityId, cursor)
	return mvc.JsonCursorData(render.BuildComments(comments, currentUser, true, false), strconv.FormatInt(cursor, 10), hasMore)
}

func (c *CommentController) GetReplies() *mvc.JsonResult {
	var (
		cursor    = params.FormValueInt64Default(c.Ctx, "cursor", 0)
		commentId = params.FormValueInt64Default(c.Ctx, "commentId", 0)
	)
	currentUser := services.UserTokenService.GetCurrent(c.Ctx)
	comments, cursor, hasMore := services.CommentService.GetReplies(commentId, cursor, 10)
	return mvc.JsonCursorData(render.BuildComments(comments, currentUser, false, true), strconv.FormatInt(cursor, 10), hasMore)
}

func (c *CommentController) PostCreate() *mvc.JsonResult {
	user := services.UserTokenService.GetCurrent(c.Ctx)
	if err := services.UserService.CheckPostStatus(user); err != nil {
		return mvc.JsonError(err)
	}
	form := model.GetCreateCommentForm(c.Ctx)
	if err := spam.CheckComment(user, form); err != nil {
		return mvc.JsonErrorMsg(err.Error())
	}

	comment, err := services.CommentService.Publish(user.Id, form)
	if err != nil {
		return mvc.JsonErrorMsg(err.Error())
	}

	return mvc.JsonData(render.BuildComment(comment))
}

func (c *CommentController) PostLikeBy(commentId int64) *mvc.JsonResult {
	user := services.UserTokenService.GetCurrent(c.Ctx)
	if user == nil {
		return mvc.JsonError(simple.ErrorNotLogin)
	}
	err := services.UserLikeService.CommentLike(user.Id, commentId)
	if err != nil {
		return mvc.JsonErrorMsg(err.Error())
	}
	return mvc.JsonSuccess()
}
