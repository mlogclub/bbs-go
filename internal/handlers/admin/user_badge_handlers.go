package admin

import (
	"bbs-go/internal/cache"
	"bbs-go/internal/models"
	"bbs-go/internal/services"
	"strconv"

	"github.com/gin-gonic/gin"

	"bbs-go/internal/pkg/ginx"
	"bbs-go/internal/pkg/params"

	"github.com/mlogclub/simple/web"
)

func UserBadgeDetail(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}

	t := services.UserBadgeService.Get(id)
	if t == nil {
		ginx.WriteJSON(ctx, ginx.ErrorMessage("Not found, id="+strconv.FormatInt(id, 10)))
		return
	}
	ginx.WriteJSON(ctx, t)

}

func UserBadgeList(ctx *gin.Context) {
	list, paging := services.UserBadgeService.FindPageByCnd(params.NewPagedSqlCnd(ctx,
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
	ginx.WriteJSON(ctx, &web.PageResult{Results: web.ConvertList(list, func(item models.UserBadge) map[string]any {
		b := web.NewRspBuilder(item)
		if badge := cache.BadgeCache.GetByID(item.BadgeId); badge != nil {
			b.Put("icon", badge.Icon)
		}
		return b.Build()
	}), Page: paging})

}
