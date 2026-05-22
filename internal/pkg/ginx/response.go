package ginx

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mlogclub/simple/sqls"
	"github.com/mlogclub/simple/web"
)

type cursorData struct {
	results any
	cursor  string
	hasMore bool
}

type pageData struct {
	results any
	paging  *sqls.Paging
}

func CursorData(results any, cursor string, hasMore bool) any {
	return cursorData{results: results, cursor: cursor, hasMore: hasMore}
}

func PageData(results any, paging *sqls.Paging) any {
	return pageData{results: results, paging: paging}
}

func ErrorMessage(message string) error {
	return errors.New(message)
}

func ErrorCode(code int, message string) error {
	return web.NewError(code, message)
}

func WriteJSON(ctx *gin.Context, result any) {
	ctx.JSON(http.StatusOK, buildJSONResult(result))
}

func WriteHttpStatusJSON(ctx *gin.Context, statusCode int, result any) {
	ctx.JSON(statusCode, buildJSONResult(result))
}

func buildJSONResult(result any) *web.JsonResult {
	switch value := result.(type) {
	case nil:
		return web.JsonSuccess()
	case *web.JsonResult:
		return value
	case web.JsonResult:
		return &value
	case *web.CodeError:
		return web.JsonError(value)
	case web.CodeError:
		return web.JsonError(&value)
	case error:
		return web.JsonError(value)
	case cursorData:
		return web.JsonCursorData(value.results, value.cursor, value.hasMore)
	case pageData:
		return web.JsonPageData(value.results, value.paging)
	case web.RspBuilder:
		return value.JsonResult()
	case *web.RspBuilder:
		return value.JsonResult()
	default:
		return web.JsonData(result)
	}
}
