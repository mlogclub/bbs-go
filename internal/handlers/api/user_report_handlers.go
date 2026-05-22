package api

import (
	"bbs-go/internal/models"
	"bbs-go/internal/models/req"
	"bbs-go/internal/pkg/common"
	"bbs-go/internal/services"

	"github.com/gin-gonic/gin"

	"bbs-go/internal/pkg/ginx"

	"github.com/mlogclub/simple/common/dates"
)

func UserReportSubmit(ctx *gin.Context) {
	var req req.UserReportReq
	if err := ginx.Bind(ctx, &req); err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}
	report := &models.UserReport{
		DataId:     req.DataId,
		DataType:   req.DataType,
		Reason:     req.Reason,
		CreateTime: dates.NowTimestamp(),
	}

	if user := common.GetCurrentUser(ctx); user != nil {
		report.UserId = user.Id
	}
	services.UserReportService.Create(report)
	ginx.WriteJSON(ctx, nil)

}
