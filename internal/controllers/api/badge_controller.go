package api

import (
	"bbs-go/internal/cache"
	"bbs-go/internal/models"
	"bbs-go/internal/models/resp"
	"bbs-go/internal/pkg/common"

	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple/web"
)

type BadgeController struct {
	Ctx iris.Context
}

// GetBadges 获取全部勋章列表
func (c *BadgeController) GetBadges() *web.JsonResult {
	userId := common.GetID(c.Ctx, "userId")

	// 用户已获得的勋章（走缓存）
	userBadges := map[int64]*models.UserBadge{}
	if userId > 0 {
		list := cache.UserBadgeCache.GetByUser(userId)
		for i := range list {
			ub := list[i]
			userBadges[ub.BadgeId] = &ub
		}
	}

	// 全部勋章
	badges := cache.BadgeCache.GetAll()

	var ret []resp.BadgeResponse
	for i := range badges {
		b := badges[i]
		item := resp.BadgeResponse{
			Id:          b.Id,
			Name:        b.Name,
			Title:       b.Title,
			Description: b.Description,
			Icon:        b.Icon,
			SortNo:      b.SortNo,
			Status:      b.Status,
		}
		if ub, ok := userBadges[b.Id]; ok {
			item.Owned = true
			item.Worn = ub.IsWorn
			item.ObtainTime = ub.CreateTime
		}
		ret = append(ret, item)
	}
	return web.JsonData(ret)
}
