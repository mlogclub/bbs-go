package admin

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"bbs-go/internal/pkg/ginx"
	"bbs-go/internal/pkg/params"

	"github.com/mlogclub/simple/web"

	"bbs-go/internal/handlers/render"
	"bbs-go/internal/services"
)

func UserScoreLogDetail(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}

	t := services.UserScoreLogService.Get(id)
	if t == nil {
		ginx.WriteJSON(ctx, ginx.ErrorMessage("Not found, id="+strconv.FormatInt(id, 10)))
		return
	}
	ginx.WriteJSON(ctx, t)

}

func UserScoreLogList(ctx *gin.Context) {
	list, paging := services.UserScoreLogService.FindPageByParams(params.NewQueryParams(ctx).
		EqByReq("user_id").EqByReq("source_type").EqByReq("source_id").EqByReq("type").PageByReq().Desc("id"))

	var results []map[string]interface{}
	for _, userScoreLog := range list {
		user := render.BuildUserInfoDefaultIfNull(userScoreLog.UserId)
		item := web.NewRspBuilder(userScoreLog).Put("user", user).Build()
		results = append(results, item)
	}

	ginx.WriteJSON(ctx, &web.PageResult{Results: results, Page: paging})

}
