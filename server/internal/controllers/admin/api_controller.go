package admin

import (
	"bbs-go/internal/models"
	"bbs-go/internal/models/dto"
	"bbs-go/internal/services"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple/common/arrs"
	"github.com/mlogclub/simple/common/dates"
	"github.com/mlogclub/simple/sqls"
	"github.com/mlogclub/simple/web"
	"github.com/mlogclub/simple/web/params"
)

type ApiController struct {
	Ctx iris.Context
}

func (c *ApiController) GetInit() *web.JsonResult {
	var (
		list       []dto.ApiRoute
		apiMethods = make(map[string][]string)
		exists     = make(map[string]bool)
		app        = c.Ctx.Application()
	)

	routes := app.GetRoutesReadOnly()
	for _, route := range routes {
		methods := apiMethods[route.Path()]
		methods = append(methods, route.Method())
		apiMethods[route.Path()] = methods
	}

	for _, route := range routes {
		if !strings.HasPrefix(route.Path(), "/api/admin") {
			continue
		}

		if _, found := exists[route.Path()]; found {
			continue
		}

		if !arrs.Contains([]string{http.MethodGet, http.MethodPost, http.MethodDelete, http.MethodPut}, route.Method()) {
			continue
		}

		method := route.Method()
		// 约定trace不会再项目中用，如果遇到这个说明接口设置的是Any
		if methods := apiMethods[route.Path()]; arrs.Contains(methods, http.MethodTrace) {
			method = "ANY"
		}

		list = append(list, dto.ApiRoute{
			Method: method,
			Path:   route.Path(),
			Name:   route.MainHandlerName(),
		})
		exists[route.Path()] = true

		slog.Info("Route: " + method + " " + route.Path() + " " + route.MainHandlerName())
	}

	services.ApiService.Init(list)
	return web.JsonSuccess()
}

func (c *ApiController) GetBy(id int64) *web.JsonResult {
	t := services.ApiService.Get(id)
	if t == nil {
		return web.JsonErrorMsg("Not found, id=" + strconv.FormatInt(id, 10))
	}
	return web.JsonData(t)
}

func (c *ApiController) AnyList() *web.JsonResult {
	list, paging := services.ApiService.FindPageByCnd(params.NewPagedSqlCnd(c.Ctx,
		params.QueryFilter{
			ParamName: "name",
			Op:        params.Like,
		},
		params.QueryFilter{
			ParamName: "path",
			Op:        params.Like,
		},
	).Desc("id"))
	return web.JsonData(&web.PageResult{Results: list, Page: paging})
}

func (c *ApiController) PostCreate() *web.JsonResult {
	t := &models.Api{}
	if err := params.ReadForm(c.Ctx, t); err != nil {
		return web.JsonErrorMsg(err.Error())
	}

	t.CreateTime = dates.NowTimestamp()
	t.UpdateTime = dates.NowTimestamp()
	if err := services.ApiService.Create(t); err != nil {
		return web.JsonErrorMsg(err.Error())
	}
	return web.JsonData(t)
}

func (c *ApiController) PostUpdate() *web.JsonResult {
	id, err := params.FormValueInt64(c.Ctx, "id")
	if err != nil {
		return web.JsonErrorMsg(err.Error())
	}
	t := services.ApiService.Get(id)
	if t == nil {
		return web.JsonErrorMsg("entity not found")
	}

	if err := params.ReadForm(c.Ctx, t); err != nil {
		return web.JsonErrorMsg(err.Error())
	}

	t.UpdateTime = dates.NowTimestamp()
	if err := services.ApiService.Update(t); err != nil {
		return web.JsonErrorMsg(err.Error())
	}
	return web.JsonData(t)
}

// func (c *ApiController) PostDelete() *web.JsonResult {
// 	ids := params.GetInt64Arr(c.Ctx, "ids")
// 	if len(ids) == 0 {
// 		return web.JsonErrorMsg("delete ids is empty")
// 	}
// 	for _, id := range ids {
// 		services.ApiService.Delete(id)
// 	}
// 	return web.JsonSuccess()
// }

func (c *ApiController) GetList_all() *web.JsonResult {
	list := services.ApiService.Find(sqls.NewCnd().Asc("id"))
	return web.JsonData(list)
}
