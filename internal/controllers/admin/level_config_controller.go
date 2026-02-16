package admin

import (
	"bbs-go/internal/models"
	"bbs-go/internal/services"
	"log/slog"
	"strconv"

	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple/web"
	"github.com/mlogclub/simple/web/params"
)

type LevelConfigController struct {
	Ctx iris.Context
}

func (c *LevelConfigController) GetBy(id int64) *web.JsonResult {
	t := services.LevelConfigService.Get(id)
	if t == nil {
		return web.JsonErrorMsg("Not found, id=" + strconv.FormatInt(id, 10))
	}
	return web.JsonData(t)
}

func (c *LevelConfigController) AnyList() *web.JsonResult {
	list, paging := services.LevelConfigService.FindPageByCnd(params.NewPagedSqlCnd(c.Ctx,
		params.QueryFilter{
			ParamName: "id",
		},
		params.QueryFilter{
			ParamName: "status",
			Op:        params.Eq,
		},
	).Asc("level"))
	return web.JsonData(&web.PageResult{Results: list, Page: paging})
}

// PostSave_all 批量保存等级配置（Level 必须从 1 开始且连续，NeedExp 必须严格递增）
func (c *LevelConfigController) PostSave_all() *web.JsonResult {
	var items []models.LevelConfig
	if err := c.Ctx.ReadJSON(&items); err != nil {
		return web.JsonError(err)
	}
	if err := services.LevelConfigService.SaveAll(items); err != nil {
		slog.Error("save level config failed", slog.Any("err", err))
		return web.JsonError(err)
	}
	return web.JsonSuccess()
}
