package admin

import (
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"

	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple"

	"bbs-go/controllers/render"
	"bbs-go/model"
	"bbs-go/services"
	"bbs-go/services/cache"
)

type ArticleController struct {
	Ctx iris.Context
}

func (c *ArticleController) GetBy(id int64) *simple.JsonResult {
	t := services.ArticleService.Get(id)
	if t == nil {
		return simple.JsonErrorMsg("Not found, id=" + strconv.FormatInt(id, 10))
	}
	return simple.JsonData(t)
}

func (c *ArticleController) AnyList() *simple.JsonResult {
	list, paging := services.ArticleService.FindPageByParams(simple.NewQueryParams(c.Ctx).
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

		// 标签
		tagIds := cache.ArticleTagCache.Get(article.Id)
		tags := cache.TagCache.GetList(tagIds)
		builder.Put("tags", render.BuildTags(tags))

		results = append(results, builder.Build())
	}

	return simple.JsonData(&simple.PageResult{Results: results, Page: paging})
}

func (c *ArticleController) PostCreate() *simple.JsonResult {
	return simple.JsonErrorMsg("为实现")
}

func (c *ArticleController) PostUpdate() *simple.JsonResult {
	id := c.Ctx.PostValueInt64Default("id", 0)
	if id <= 0 {
		return simple.JsonErrorMsg("id is required")
	}
	t := services.ArticleService.Get(id)
	if t == nil {
		return simple.JsonErrorMsg("entity not found")
	}

	c.Ctx.ReadForm(t)

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

func (c *ArticleController) PostDelete() *simple.JsonResult {
	id := c.Ctx.PostValueInt64Default("id", 0)
	if id <= 0 {
		return simple.JsonErrorMsg("id is required")
	}
	err := services.ArticleService.Delete(id)
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	return simple.JsonSuccess()
}
