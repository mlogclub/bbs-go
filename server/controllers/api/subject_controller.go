package api

import (
	"strconv"

	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple"

	"github.com/mlogclub/bbs-go/common/urls"
	"github.com/mlogclub/bbs-go/model"
	"github.com/mlogclub/bbs-go/services"
)

type SubjectController struct {
	Ctx iris.Context
}

func (this *SubjectController) GetBy(subjectId int64) *simple.JsonResult {
	s := services.SubjectService.Get(subjectId)
	return simple.JsonData(s)
}

func (this *SubjectController) GetContents() *simple.JsonResult {
	cursor := simple.FormValueInt64Default(this.Ctx, "cursor", 0)
	subjectId := simple.FormValueInt64Default(this.Ctx, "subjectId", 0)
	contents, cursor := services.SubjectContentService.GetSubjectContents(subjectId, cursor)

	var results []map[string]interface{}
	for _, c := range contents {
		url := ""
		if c.EntityType == model.EntityTypeArticle {
			url = urls.ArticleUrl(c.EntityId)
		} else if c.EntityType == model.EntityTypeTopic {
			url = urls.TopicUrl(c.EntityId)
		}
		item := map[string]interface{}{
			"subjectContentId": c.Id,
			"subjectId":        c.SubjectId,
			"entityType:":      c.EntityType,
			"entityId":         c.EntityId,
			"url":              url,
			"title":            c.Title,
			"summary":          c.Summary,
			"createTime":       c.CreateTime,
		}
		s := services.SubjectService.Get(c.SubjectId)
		if s != nil {
			item["subjectTitle"] = s.Title
		}
		results = append(results, item)
	}

	return simple.JsonCursorData(results, strconv.FormatInt(cursor, 10))
}
