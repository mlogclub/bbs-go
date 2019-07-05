package controllers

import (
	"github.com/kataras/iris"
	"github.com/mlogclub/mlog/utils/session"
	"github.com/mlogclub/simple"

	"github.com/mlogclub/mlog/controllers/render"
	"github.com/mlogclub/mlog/model"
	"github.com/mlogclub/mlog/services"
)

type CommentController struct {
	Ctx iris.Context
}

func (this *CommentController) GetList() *simple.JsonResult {
	entityType, err := simple.FormValueRequired(this.Ctx, "entityType")
	if err != nil {
		return simple.ErrorMsg(err.Error())
	}
	entityId, err := simple.FormValueInt64(this.Ctx, "entityId")
	if err != nil {
		return simple.ErrorMsg(err.Error())
	}

	cursor := simple.FormValueInt64Default(this.Ctx, "cursor", 0)

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
	user := session.GetCurrentUser(this.Ctx)
	if user == nil {
		return simple.Error(simple.ErrorNotLogin)
	}
	entityType, err := simple.FormValueRequired(this.Ctx, "entityType")
	if err != nil {
		return simple.ErrorMsg(err.Error())
	}
	entityId, err := simple.FormValueInt64(this.Ctx, "entityId")
	if err != nil {
		return simple.ErrorMsg(err.Error())
	}
	content, err := simple.FormValueRequired(this.Ctx, "content")
	if err != nil {
		return simple.ErrorMsg(err.Error())
	}
	quoteId := simple.FormValueInt64Default(this.Ctx, "quoteId", 0)

	comment := &model.Comment{
		UserId:     user.Id,
		EntityType: entityType,
		EntityId:   entityId,
		Content:    content,
		QuoteId:    quoteId,
		Status:     model.CommentStatusOk,
		CreateTime: simple.NowTimestamp(),
	}
	err = services.CommentService.Create(comment)
	if err != nil {
		return simple.ErrorMsg(err.Error())
	}

	if entityType == model.EntityTypeTopic {
		services.TopicService.SetLastCommentTime(entityId, simple.NowTimestamp())
	}

	services.MessageService.SendCommentMsg(comment)

	return simple.JsonData(render.BuildComment(*comment))
}
