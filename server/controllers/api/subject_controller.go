package api

import (
	"github.com/kataras/iris/context"
	"github.com/mlogclub/simple"

	"github.com/mlogclub/bbs-go/common/urls"
	"github.com/mlogclub/bbs-go/model"
	"github.com/mlogclub/bbs-go/services"
)

type SubjectController struct {
	Ctx context.Context
}

func (this *SubjectController) GetAnalyze() *simple.JsonResult {
	user := services.UserTokenService.GetCurrent(this.Ctx)
	if user == nil || user.Id != 1 {
		return simple.JsonErrorMsg("无权限")
	}
	go func() {
		services.ArticleService.Scan(func(articles []model.Article) bool {
			for _, article := range articles {
				if article.Status == model.ArticleStatusPublished {
					services.SubjectContentService.AnalyzeArticle(&article)
				}
			}
			return true
		})
	}()
	return simple.JsonSuccess()
}

func (this *SubjectController) GetBy(subjectId int64) *simple.JsonResult {
	s := services.SubjectService.Get(subjectId)
	return simple.JsonData(s)
}

func (this *SubjectController) GetContents() *simple.JsonResult {
	page := simple.FormValueIntDefault(this.Ctx, "page", 1)
	subjectId := simple.FormValueInt64Default(this.Ctx, "subjectId", 0)

	cnd := simple.NewSqlCnd().Eq("deleted", false).Page(page, 20).Desc("id")
	if subjectId > 0 {
		cnd.Eq("subject_id", subjectId)
	}

	contents, paging := services.SubjectContentService.FindPageByCnd(cnd)

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

	return simple.JsonPageData(results, paging)
}
