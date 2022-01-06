package api

import (
	"bbs-go/model"
	"bbs-go/spam"
	"strconv"

	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple"

	"bbs-go/controllers/render"
	"bbs-go/services"
)

type CommentController struct {
	Ctx iris.Context
}

func (c *CommentController) GetComments() *simple.JsonResult {
	var (
		err        error
		cursor     int64
		entityType string
		entityId   int64
	)
	cursor = simple.FormValueInt64Default(c.Ctx, "cursor", 0)

	if entityType, err = simple.FormValueRequired(c.Ctx, "entityType"); err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	if entityId, err = simple.FormValueInt64(c.Ctx, "entityId"); err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	currentUser := services.UserTokenService.GetCurrent(c.Ctx)
	comments, cursor, hasMore := services.CommentService.GetComments(entityType, entityId, cursor)
	return simple.JsonCursorData(render.BuildComments(comments, currentUser, true, false), strconv.FormatInt(cursor, 10), hasMore)
}

func (c *CommentController) GetReplies() *simple.JsonResult {
	var (
		cursor    = simple.FormValueInt64Default(c.Ctx, "cursor", 0)
		commentId = simple.FormValueInt64Default(c.Ctx, "commentId", 0)
	)
	currentUser := services.UserTokenService.GetCurrent(c.Ctx)
	comments, cursor, hasMore := services.CommentService.GetReplies(commentId, cursor, 10)
	return simple.JsonCursorData(render.BuildComments(comments, currentUser, false, true), strconv.FormatInt(cursor, 10), hasMore)
}

func (c *CommentController) PostCreate() *simple.JsonResult {
	user := services.UserTokenService.GetCurrent(c.Ctx)
	if err := services.UserService.CheckPostStatus(user); err != nil {
		return simple.JsonError(err)
	}
	form := model.GetCreateCommentForm(c.Ctx)
	if err := spam.CheckComment(user, form); err != nil {
		return simple.JsonErrorMsg(err.Error())
	}

	comment, err := services.CommentService.Publish(user.Id, form)
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}

	return simple.JsonData(render.BuildComment(comment))
}

func (c *CommentController) PostLikeBy(commentId int64) *simple.JsonResult {
	user := services.UserTokenService.GetCurrent(c.Ctx)
	if user == nil {
		return simple.JsonError(simple.ErrorNotLogin)
	}
	err := services.UserLikeService.CommentLike(user.Id, commentId)
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	return simple.JsonSuccess()
}
