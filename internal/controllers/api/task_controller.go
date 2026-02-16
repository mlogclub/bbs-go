package api

import (
	"bbs-go/internal/controllers/render"
	"bbs-go/internal/models/constants"
	"bbs-go/internal/models/resp"
	"bbs-go/internal/pkg/common"
	"bbs-go/internal/services"

	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple/common/dates"
	"github.com/mlogclub/simple/sqls"
	"github.com/mlogclub/simple/web"
)

type TaskController struct {
	Ctx iris.Context
}

func (c *TaskController) GetTasks() *web.JsonResult {
	now := dates.NowTimestamp()

	groupName := c.Ctx.URLParam("groupName")
	cnd := sqls.NewCnd().
		Eq("status", constants.StatusOk).
		Asc("sort_no").
		Asc("id")
	if groupName != "" {
		cnd.Eq("group_name", groupName)
	}

	taskConfigs := services.TaskConfigService.Find(cnd)

	user := common.GetCurrentUser(c.Ctx)
	resp := make([]resp.TaskResponse, 0, len(taskConfigs))
	for i := range taskConfigs {
		cfg := taskConfigs[i]

		// 时间窗口过滤
		if cfg.StartTime > 0 && now < cfg.StartTime {
			continue
		}
		if cfg.EndTime > 0 && now > cfg.EndTime {
			continue
		}

		resp = append(resp, render.BuildTask(&cfg, user, now))
	}

	return web.JsonData(resp)
}

// 任务分组列表
func (c *TaskController) GetGroups() *web.JsonResult {
	return web.JsonData(render.BuildTaskGroups())
}
