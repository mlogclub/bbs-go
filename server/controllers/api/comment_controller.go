package api

import (
	"bbs-go/model"
	"bbs-go/model/constants"
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

func (c *CommentController) GetFuck() *simple.JsonResult {
	go func() {
		users := services.UserService.Find(simple.NewSqlCnd().Eq("forbidden_end_time", -1))
		for _, user := range users {
			// 删除评论
			services.CommentService.ScanByUser(user.Id, func(comments []model.Comment) {
				for _, comment := range comments {
					if comment.Status != constants.StatusDeleted {
						_ = services.CommentService.Delete(comment.Id)
					}
				}
			})
		}
	}()
	return simple.JsonSuccess()
}

func (c *CommentController) GetList() *simple.JsonResult {
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

	comments, cursor, hasMore := services.CommentService.GetComments(entityType, entityId, cursor)
	return simple.JsonCursorData(render.BuildComments(comments), strconv.FormatInt(cursor, 10), hasMore)
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

	return simple.JsonData(render.BuildComment(*comment))
}
