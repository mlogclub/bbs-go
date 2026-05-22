package admin

import (
	"bbs-go/internal/models"
	"bbs-go/internal/services"
	"log/slog"
	"strconv"

	"github.com/gin-gonic/gin"

	"bbs-go/internal/pkg/ginx"
	"bbs-go/internal/pkg/params"

	"github.com/mlogclub/simple/web"
)

// PostSave_all 批量保存等级配置（Level 必须从 1 开始且连续，NeedExp 必须严格递增）
func LevelConfigDetail(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}

	t := services.LevelConfigService.Get(id)
	if t == nil {
		ginx.WriteJSON(ctx, ginx.ErrorMessage("Not found, id="+strconv.FormatInt(id, 10)))
		return
	}
	ginx.WriteJSON(ctx, t)

}

func LevelConfigList(ctx *gin.Context) {
	list, paging := services.LevelConfigService.FindPageByCnd(params.NewPagedSqlCnd(ctx,
		params.QueryFilter{
			ParamName: "id",
		},
		params.QueryFilter{
			ParamName: "status",
			Op:        params.Eq,
		},
	).Asc("level"))
	ginx.WriteJSON(ctx, &web.PageResult{Results: list, Page: paging})

}

func LevelConfigSaveAll(ctx *gin.Context) {
	var items []models.LevelConfig
	if err := ginx.BindJSON(ctx, &items); err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}
	if err := services.LevelConfigService.SaveAll(items); err != nil {
		slog.Error("save level config failed", slog.Any("err", err))
		ginx.WriteJSON(ctx, err)
		return
	}
	ginx.WriteJSON(ctx, nil)

}
