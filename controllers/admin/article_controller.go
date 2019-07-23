package admin

import (
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"

	"github.com/kataras/iris"
	"github.com/mlogclub/simple"

	"github.com/mlogclub/mlog/controllers/render"
	"github.com/mlogclub/mlog/model"
	"github.com/mlogclub/mlog/services"
	"github.com/mlogclub/mlog/services/cache"
	"github.com/mlogclub/mlog/services/collect"
)

type ArticleController struct {
	Ctx iris.Context
}

func (this *ArticleController) GetBy(id int64) *simple.JsonResult {
	t := services.ArticleService.Get(id)
	if t == nil {
		return simple.ErrorMsg("Not found, id=" + strconv.FormatInt(id, 10))
	}
	return simple.JsonData(t)
}

func (this *ArticleController) AnyList() *simple.JsonResult {
	list, paging := services.ArticleService.Query(simple.NewParamQueries(this.Ctx).
		EqAuto("status").LikeAuto("title").PageAuto().Desc("id"))

	var results []map[string]interface{}
	for _, article := range list {
		builder := simple.NewRspBuilderExcludes(article, "content")

		// 用户
		builder = builder.Put("user", render.BuildUserDefaultIfNull(article.UserId))

		// 简介
		if article.ContentType == model.ArticleContentTypeMarkdown {
			mr := simple.NewMd().Run(article.Content)
			if len(article.Summary) == 0 {
				builder.Put("summary", mr.SummaryText)
			}
		} else {
			if len(article.Summary) == 0 {
				doc, err := goquery.NewDocumentFromReader(strings.NewReader(article.Content))
				if err != nil {
					builder.Put("summary", simple.GetSummary(doc.Text(), 256))
				}
			}
		}

		// 分类
		builder.Put("category", render.BuildCategory(article.CategoryId))

		// 标签
		tagIds := cache.ArticleTagCache.Get(article.Id)
		tags := cache.TagCache.GetList(tagIds)
		builder.Put("tags", render.BuildTags(tags))

		results = append(results, builder.Build())
	}

	return simple.JsonData(&simple.PageResult{Results: results, Page: paging})
}

func (this *ArticleController) PostCreate() *simple.JsonResult {
	return simple.ErrorMsg("为实现")
}

func (this *ArticleController) PostUpdate() *simple.JsonResult {
	id := this.Ctx.PostValueInt64Default("id", 0)
	if id <= 0 {
		return simple.ErrorMsg("id is required")
	}
	t := services.ArticleService.Get(id)
	if t == nil {
		return simple.ErrorMsg("entity not found")
	}

	this.Ctx.ReadForm(t)

	// 数据校验
	if len(t.Title) == 0 {
		return simple.ErrorMsg("标题不能为空")
	}
	if len(t.Content) == 0 {
		return simple.ErrorMsg("内容不能为空")
	}
	if len(t.ContentType) == 0 {
		return simple.ErrorMsg("请选择内容格式")
	}

	t.UpdateTime = simple.NowTimestamp()
	err := services.ArticleService.Update(t)
	if err != nil {
		return simple.ErrorMsg(err.Error())
	}

	return simple.JsonData(t)
}

func (this *ArticleController) PostDelete() *simple.JsonResult {
	id := this.Ctx.PostValueInt64Default("id", 0)
	if id <= 0 {
		return simple.ErrorMsg("id is required")
	}
	services.ArticleService.UpdateColumn(id, "status", model.ArticleStatusDeleted)
	return simple.Success()
}

func (this *ArticleController) PostCollect() *simple.JsonResult {
	url := this.Ctx.PostValue("url")
	if len(url) == 0 {
		return simple.ErrorMsg("链接不存在")
	}
	title, content, err := collect.Collect(url, true)
	if err != nil {
		return simple.ErrorMsg(err.Error())
	}
	return simple.NewEmptyRspBuilder().Put("title", title).Put("content", content).JsonResult()
}
