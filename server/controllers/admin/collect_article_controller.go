package admin

import (
	"strconv"

	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple"

	"github.com/mlogclub/bbs-go/controllers/render"
	"github.com/mlogclub/bbs-go/services"
)

type CollectArticleController struct {
	Ctx iris.Context
}

func (this *CollectArticleController) GetBy(id int64) *simple.JsonResult {
	t := services.CollectArticleService.Get(id)
	if t == nil {
		return simple.JsonErrorMsg("Not found, id=" + strconv.FormatInt(id, 10))
	}
	return simple.JsonData(t)
}

func (this *CollectArticleController) AnyList() *simple.JsonResult {
	list, paging := services.CollectArticleService.FindPageByParams(simple.NewQueryParams(this.Ctx).
		EqByReq("rule_id").
		EqByReq("link_id").
		EqByReq("article_id").
		EqByReq("status").
		PageByReq().Desc("id"))
	var results []map[string]interface{}
	for _, article := range list {
		item := simple.NewRspBuilderExcludes(article, "content").Build()
		item["user"] = render.BuildUserDefaultIfNull(article.UserId)
		results = append(results, item)
	}

	return simple.JsonData(&simple.PageResult{Results: results, Page: paging})
}
