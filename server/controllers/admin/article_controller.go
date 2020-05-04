package admin

import (
	"strconv"

	"bbs-go/model"

	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple"

	"bbs-go/common"
	"bbs-go/controllers/render"
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
		builder.Put("summary", common.GetSummary(article.ContentType, article.Content))

		// 标签
		tagIds := cache.ArticleTagCache.Get(article.Id)
		tags := cache.TagCache.GetList(tagIds)
		builder.Put("tags", render.BuildTags(tags))

		results = append(results, builder.Build())
	}

	return simple.JsonData(&simple.PageResult{Results: results, Page: paging})
}

func (c *ArticleController) PostCreate() *simple.JsonResult {
	return simple.JsonErrorMsg("未实现")
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

	if err := simple.ReadForm(c.Ctx, t); err != nil {
		return simple.JsonErrorMsg(err.Error())
	}

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

func (c *ArticleController) PostPending() *simple.JsonResult {
	id := c.Ctx.PostValueInt64Default("id", 0)
	if id <= 0 {
		return simple.JsonErrorMsg("id is required")
	}
	err := services.ArticleService.UpdateColumn(id, "status", model.StatusOk)
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	return simple.JsonSuccess()
}
