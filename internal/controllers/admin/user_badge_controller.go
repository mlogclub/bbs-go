package admin

import (
	"bbs-go/internal/cache"
	"bbs-go/internal/models"
	"bbs-go/internal/services"
	"strconv"

	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple/web"
	"github.com/mlogclub/simple/web/params"
)

type UserBadgeController struct {
	Ctx iris.Context
}

func (c *UserBadgeController) GetBy(id int64) *web.JsonResult {
	t := services.UserBadgeService.Get(id)
	if t == nil {
		return web.JsonErrorMsg("Not found, id=" + strconv.FormatInt(id, 10))
	}
	return web.JsonData(t)
}

func (c *UserBadgeController) AnyList() *web.JsonResult {
	list, paging := services.UserBadgeService.FindPageByCnd(params.NewPagedSqlCnd(c.Ctx,
		params.QueryFilter{
			ParamName: "userId",
			Op:        params.Eq,
		},
		params.QueryFilter{
			ParamName: "badgeId",
			Op:        params.Eq,
		},
		params.QueryFilter{
			ParamName: "sourceType",
			Op:        params.Eq,
		},
	).Desc("id"))
	return web.JsonData(&web.PageResult{Results: web.ConvertList(list, func(item models.UserBadge) map[string]any {
		b := web.NewRspBuilder(item)
		if badge := cache.BadgeCache.GetByID(item.BadgeId); badge != nil {
			b.Put("icon", badge.Icon)
		}
		return b.Build()
	}), Page: paging})
}
