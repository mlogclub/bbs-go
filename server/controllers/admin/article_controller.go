package admin

import (
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"

	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple"

	"github.com/mlogclub/bbs-go/controllers/render"
	"github.com/mlogclub/bbs-go/model"
	"github.com/mlogclub/bbs-go/services"
	"github.com/mlogclub/bbs-go/services/cache"
	"github.com/mlogclub/bbs-go/services/collect"
)

type ArticleController struct {
	Ctx iris.Context
}

func (this *ArticleController) GetBy(id int64) *simple.JsonResult {
	t := services.ArticleService.Get(id)
	if t == nil {
		return simple.JsonErrorMsg("Not found, id=" + strconv.FormatInt(id, 10))
	}
	return simple.JsonData(t)
}

func (this *ArticleController) AnyList() *simple.JsonResult {
	list, paging := services.ArticleService.FindPageByParams(simple.NewQueryParams(this.Ctx).
		EqByReq("id").EqByReq("user_id").EqByReq("status").LikeByReq("title").PageByReq().Desc("id"))

	var results []map[string]interface{}
	for _, article := range list {
		builder := simple.NewRspBuilderExcludes(article, "content")

		// 用户
		builder = builder.Put("user", render.BuildUserDefaultIfNull(article.UserId))

		// 简介
		if article.ContentType == model.ContentTypeMarkdown {
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
		if article.CategoryId > 0 {
			category := cache.CategoryCache.Get(article.CategoryId)
			builder.Put("category", render.BuildCategory(category))
		}

		// 标签
		tagIds := cache.ArticleTagCache.Get(article.Id)
		tags := cache.TagCache.GetList(tagIds)
		builder.Put("tags", render.BuildTags(tags))

		results = append(results, builder.Build())
	}

	return simple.JsonData(&simple.PageResult{Results: results, Page: paging})
}

func (this *ArticleController) PostCreate() *simple.JsonResult {
	return simple.JsonErrorMsg("为实现")
}

func (this *ArticleController) PostUpdate() *simple.JsonResult {
	id := this.Ctx.PostValueInt64Default("id", 0)
	if id <= 0 {
		return simple.JsonErrorMsg("id is required")
	}
	t := services.ArticleService.Get(id)
	if t == nil {
		return simple.JsonErrorMsg("entity not found")
	}

	this.Ctx.ReadForm(t)

	// 数据校验
	if len(t.Title) == 0 {
		return simple.JsonErrorMsg("标题不能为空")
	}
	if len(t.Content) == 0 {
		return simple.JsonErrorMsg("内容不能为空")
	}
	if len(t.ContentType) == 0 {
		return simple.JsonErrorMsg("请选择内容格式")
	}

	t.UpdateTime = simple.NowTimestamp()
	err := services.ArticleService.Update(t)
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}

	return simple.JsonData(t)
}

func (this *ArticleController) PostDelete() *simple.JsonResult {
	id := this.Ctx.PostValueInt64Default("id", 0)
	if id <= 0 {
		return simple.JsonErrorMsg("id is required")
	}
	err := services.ArticleService.Delete(id)
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	return simple.JsonSuccess()
}

func (this *ArticleController) PostCollect() *simple.JsonResult {
	url := this.Ctx.PostValue("url")
	if len(url) == 0 {
		return simple.JsonErrorMsg("链接不存在")
	}
	title, content, err := collect.Collect(url, true)
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	return simple.NewEmptyRspBuilder().Put("title", title).Put("content", content).JsonResult()
}
