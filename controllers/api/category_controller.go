package api

import (
	"github.com/kataras/iris/context"
	"github.com/mlogclub/simple"

	"github.com/mlogclub/bbs-go/controllers/render"
	"github.com/mlogclub/bbs-go/services/cache"
)

type CategoryController struct {
	Ctx context.Context
}

func (this *CategoryController) GetBy(categoryId int64) *simple.JsonResult {
	category := cache.CategoryCache.Get(categoryId)
	if category == nil {
		return simple.JsonErrorMsg("分类不存在")
	}
	return simple.JsonData(render.BuildCategory(category))
}
