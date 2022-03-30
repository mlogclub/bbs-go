package api

import (
	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple/mvc"
	"github.com/mlogclub/simple/mvc/params"

	"bbs-go/controllers/render"
	"bbs-go/services"
)

type ProjectController struct {
	Ctx iris.Context
}

func (c *ProjectController) GetBy(projectId int64) *mvc.JsonResult {
	project := services.ProjectService.Get(projectId)
	if project == nil {
		return mvc.JsonErrorMsg("项目不存在")
	}
	return mvc.JsonData(render.BuildProject(project))
}

func (c *ProjectController) GetProjects() *mvc.JsonResult {
	page := params.FormValueIntDefault(c.Ctx, "page", 1)

	projects, paging := services.ProjectService.FindPageByParams(params.NewQueryParams(c.Ctx).
		Page(page, 20).Desc("id"))

	return mvc.JsonPageData(render.BuildSimpleProjects(projects), paging)
}
