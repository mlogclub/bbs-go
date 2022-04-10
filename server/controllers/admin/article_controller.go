package admin

import (
	"bbs-go/model/constants"
	"bbs-go/pkg/sitemap"
	"strconv"

	"bbs-go/model"

	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple/common/dates"
	"github.com/mlogclub/simple/web"
	"github.com/mlogclub/simple/web/params"

	"bbs-go/cache"
	"bbs-go/controllers/render"
	"bbs-go/pkg/common"
	"bbs-go/services"
)

type ArticleController struct {
	Ctx iris.Context
}

func (c *ArticleController) GetSitemap() *web.JsonResult {
	go func() {
		sitemap.Generate()
	}()
	return web.JsonSuccess()
}

func (c *ArticleController) GetBy(id int64) *web.JsonResult {
	t := services.ArticleService.Get(id)
	if t == nil {
		return web.JsonErrorMsg("Not found, id=" + strconv.FormatInt(id, 10))
	}
	return web.JsonData(t)
}

func (c *ArticleController) AnyList() *web.JsonResult {
	var (
		id     = params.FormValueInt64Default(c.Ctx, "id", 0)
		userId = params.FormValueInt64Default(c.Ctx, "userId", 0)
	)
	params := params.NewQueryParams(c.Ctx)
	if id > 0 {
		params.Eq("id", id)
	}
	if userId > 0 {
		params.Eq("user_id", userId)
	}
	params.EqByReq("status").EqByReq("title").PageByReq().Desc("id")

	if id <= 0 && userId <= 0 {
		return web.JsonErrorMsg("请指定查询的【文章编号】或【作者编号】")
	}
	list, paging := services.ArticleService.FindPageByParams(params)
	results := c.buildArticles(list)
	return web.JsonPageData(results, paging)
}

// GetRecent 展示最近一页数据
func (c *ArticleController) GetRecent() *web.JsonResult {
	params := params.NewQueryParams(c.Ctx).EqByReq("id").EqByReq("user_id").EqByReq("status").Desc("id").Limit(20)
	list := services.ArticleService.Find(&params.Cnd)
	results := c.buildArticles(list)
	return web.JsonData(results)
}

// 构建文章列表返回数据
func (c *ArticleController) buildArticles(articles []model.Article) []map[string]interface{} {
	var results []map[string]interface{}
	for _, article := range articles {
		builder := web.NewRspBuilderExcludes(article, "content")

		// 用户
		builder = builder.Put("user", render.BuildUserInfoDefaultIfNull(article.UserId))

		// 简介
		builder.Put("summary", common.GetSummary(article.ContentType, article.Content))

		// 标签
		tagIds := cache.ArticleTagCache.Get(article.Id)
		tags := cache.TagCache.GetList(tagIds)
		builder.Put("tags", render.BuildTags(tags))

		results = append(results, builder.Build())
	}
	return results
}

func (c *ArticleController) PostUpdate() *web.JsonResult {
	id := c.Ctx.PostValueInt64Default("id", 0)
	if id <= 0 {
		return web.JsonErrorMsg("id is required")
	}
	t := services.ArticleService.Get(id)
	if t == nil {
		return web.JsonErrorMsg("entity not found")
	}

	if err := params.ReadForm(c.Ctx, t); err != nil {
		return web.JsonErrorMsg(err.Error())
	}

	// 数据校验
	if len(t.Title) == 0 {
		return web.JsonErrorMsg("标题不能为空")
	}
	if len(t.Content) == 0 {
		return web.JsonErrorMsg("内容不能为空")
	}
	if len(t.ContentType) == 0 {
		return web.JsonErrorMsg("请选择内容格式")
	}

	t.UpdateTime = dates.NowTimestamp()
	err := services.ArticleService.Update(t)
	if err != nil {
		return web.JsonErrorMsg(err.Error())
	}

	return web.JsonData(t)
}

func (c *ArticleController) GetTags() *web.JsonResult {
	articleId := params.FormValueInt64Default(c.Ctx, "articleId", 0)
	tags := services.ArticleService.GetArticleTags(articleId)
	return web.JsonData(render.BuildTags(tags))
}

func (c *ArticleController) PutTags() *web.JsonResult {
	var (
		articleId = params.FormValueInt64Default(c.Ctx, "articleId", 0)
		tags      = params.FormValueStringArray(c.Ctx, "tags")
	)
	services.ArticleService.PutTags(articleId, tags)
	return web.JsonData(render.BuildTags(services.ArticleService.GetArticleTags(articleId)))
}

func (c *ArticleController) PostDelete() *web.JsonResult {
	id := c.Ctx.PostValueInt64Default("id", 0)
	if id <= 0 {
		return web.JsonErrorMsg("id is required")
	}
	err := services.ArticleService.Delete(id)
	if err != nil {
		return web.JsonErrorMsg(err.Error())
	}
	return web.JsonSuccess()
}

func (c *ArticleController) PostPending() *web.JsonResult {
	id := c.Ctx.PostValueInt64Default("id", 0)
	if id <= 0 {
		return web.JsonErrorMsg("id is required")
	}
	err := services.ArticleService.UpdateColumn(id, "status", constants.StatusOk)
	if err != nil {
		return web.JsonErrorMsg(err.Error())
	}
	return web.JsonSuccess()
}
