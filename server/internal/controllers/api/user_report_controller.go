package api

import (
	"bbs-go/internal/models"
	"bbs-go/internal/services"

	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple/common/dates"
	"github.com/mlogclub/simple/web"
	"github.com/mlogclub/simple/web/params"
)

type UserReportController struct {
	Ctx iris.Context
}

func (c *UserReportController) PostSubmit() *web.JsonResult {
	var (
		dataId, _ = params.FormValueInt64(c.Ctx, "dataId")
		dataType  = params.FormValue(c.Ctx, "dataId")
		reason    = params.FormValue(c.Ctx, "reason")
	)
	report := &models.UserReport{
		DataId:     dataId,
		DataType:   dataType,
		Reason:     reason,
		CreateTime: dates.NowTimestamp(),
	}

	if user := services.UserTokenService.GetCurrent(c.Ctx); user != nil {
		report.UserId = user.Id
	}
	services.UserReportService.Create(report)
	return web.JsonSuccess()
}
