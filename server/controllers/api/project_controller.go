package api

import (
	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple"

	"github.com/mlogclub/bbs-go/controllers/render"
	"github.com/mlogclub/bbs-go/services"
)

type ProjectController struct {
	Ctx iris.Context
}

func (this *ProjectController) GetBy(projectId int64) *simple.JsonResult {
	project := services.ProjectService.Get(projectId)
	if project == nil {
		return simple.JsonErrorMsg("项目不存在")
	}
	return simple.JsonData(render.BuildProject(project))
}

func (this *ProjectController) GetProjects() *simple.JsonResult {
	page := simple.FormValueIntDefault(this.Ctx, "page", 1)

	projects, paging := services.ProjectService.FindPageByParams(simple.NewQueryParams(this.Ctx).
		Page(page, 20).Desc("id"))

	return simple.JsonPageData(render.BuildSimpleProjects(projects), paging)
}
