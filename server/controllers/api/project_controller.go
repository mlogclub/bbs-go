package api

import (
	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple/web"
	"github.com/mlogclub/simple/web/params"

	"bbs-go/controllers/render"
	"bbs-go/services"
)

type ProjectController struct {
	Ctx iris.Context
}

func (c *ProjectController) GetBy(projectId int64) *web.JsonResult {
	project := services.ProjectService.Get(projectId)
	if project == nil {
		return web.JsonErrorMsg("项目不存在")
	}
	return web.JsonData(render.BuildProject(project))
}

func (c *ProjectController) GetProjects() *web.JsonResult {
	page := params.FormValueIntDefault(c.Ctx, "page", 1)

	projects, paging := services.ProjectService.FindPageByParams(params.NewQueryParams(c.Ctx).
		Page(page, 20).Desc("id"))

	return web.JsonPageData(render.BuildSimpleProjects(projects), paging)
}
