package admin

import (
	"bbs-go/internal/models"
	"bbs-go/internal/pkg/common"
	"bbs-go/internal/pkg/errs"
	"bbs-go/internal/services"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mlogclub/simple/common/dates"

	"bbs-go/internal/pkg/ginx"
	"bbs-go/internal/pkg/params"

	"github.com/mlogclub/simple/web"
)

func UserReportDetail(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}

	t := services.UserReportService.Get(id)
	if t == nil {
		ginx.WriteJSON(ctx, ginx.ErrorMessage("Not found, id="+strconv.FormatInt(id, 10)))
		return
	}
	ginx.WriteJSON(ctx, t)

}

func UserReportList(ctx *gin.Context) {
	list, paging := services.UserReportService.FindPageByCnd(params.NewPagedSqlCnd(ctx,
		params.QueryFilter{
			ParamName: "dataType",
			Op:        params.Eq,
		},
		params.QueryFilter{
			ParamName: "dataId",
			Op:        params.Eq,
		},
		params.QueryFilter{
			ParamName: "auditStatus",
			Op:        params.Eq,
		},
	).Desc("id"))
	ginx.WriteJSON(ctx, &web.PageResult{Results: list, Page: paging})
}

func UserReportCreate(ctx *gin.Context) {
	t := &models.UserReport{}
	err := ginx.Bind(ctx, t)
	if err != nil {
		ginx.WriteJSON(ctx, ginx.ErrorMessage(err.Error()))
		return
	}

	err = services.UserReportService.Create(t)
	if err != nil {
		ginx.WriteJSON(ctx, ginx.ErrorMessage(err.Error()))
		return
	}
	ginx.WriteJSON(ctx, t)

}

func UserReportAudit(ctx *gin.Context) {
	id, _ := params.GetInt64(ctx, "id")
	if id <= 0 {
		ginx.WriteJSON(ctx, ginx.ErrorMessage("id is required"))
		return
	}

	auditStatus, _ := params.GetInt64(ctx, "auditStatus")
	if auditStatus != 1 && auditStatus != 2 {
		ginx.WriteJSON(ctx, ginx.ErrorMessage("auditStatus must be 1 or 2"))
		return
	}

	user := common.GetCurrentUser(ctx)
	if user == nil {
		ginx.WriteJSON(ctx, errs.NotLogin())
		return
	}

	t := services.UserReportService.Get(id)
	if t == nil {
		ginx.WriteJSON(ctx, ginx.ErrorMessage("entity not found"))
		return
	}

	t.AuditStatus = auditStatus
	t.AuditTime = dates.NowTimestamp()
	t.AuditUserId = user.Id
	if err := services.UserReportService.Update(t); err != nil {
		ginx.WriteJSON(ctx, ginx.ErrorMessage(err.Error()))
		return
	}
	ginx.WriteJSON(ctx, t)

}

func UserReportUpdate(ctx *gin.Context) {
	id, err := params.FormValueInt64(ctx, "id")
	if err != nil {
		ginx.WriteJSON(ctx, ginx.ErrorMessage(err.Error()))
		return
	}
	t := services.UserReportService.Get(id)
	if t == nil {
		ginx.WriteJSON(ctx, ginx.ErrorMessage("entity not found"))
		return
	}

	err = ginx.Bind(ctx, t)
	if err != nil {
		ginx.WriteJSON(ctx, ginx.ErrorMessage(err.Error()))
		return
	}

	err = services.UserReportService.Update(t)
	if err != nil {
		ginx.WriteJSON(ctx, ginx.ErrorMessage(err.Error()))
		return
	}
	ginx.WriteJSON(ctx, t)

}
