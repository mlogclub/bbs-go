package controllers

import (
	"strconv"

	"github.com/kataras/iris"

	"github.com/mlogclub/mlog/model"
	"github.com/mlogclub/mlog/services"
)

type ArticleShareController struct {
	Ctx                 iris.Context
	ArticleShareService *services.ArticleShareService
}

// 分享详情
func (this *ArticleShareController) GetBy(id int64) {
	share := this.ArticleShareService.Get(id)
	if share == nil || share.Status != model.ArticleShareStatusOk || share.ArticleId < 0 {
		this.Ctx.StatusCode(404)
		return
	}
	this.Ctx.Redirect("/article/"+strconv.FormatInt(share.ArticleId, 10), iris.StatusMovedPermanently)
	return
}
