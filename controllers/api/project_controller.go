package api

import (
	"github.com/kataras/iris/context"
	"github.com/mlogclub/simple"

	"github.com/mlogclub/mlog/controllers/render"
	"github.com/mlogclub/mlog/services"
)

type ProjectController struct {
	Ctx context.Context
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

	projects, paging := services.ProjectService.Query(simple.NewParamQueries(this.Ctx).
		Page(page, 20).Desc("id"))

	return simple.JsonPageData(render.BuildSimpleProjects(projects), paging)
}
