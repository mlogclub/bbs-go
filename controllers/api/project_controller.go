package api

import (
	"github.com/kataras/iris/context"
	"github.com/mlogclub/simple"

	"github.com/mlogclub/mlog/controllers/render"
	"github.com/mlogclub/mlog/model"
	"github.com/mlogclub/mlog/services"
	"github.com/mlogclub/mlog/services/collect"
)

type ProjectController struct {
	Ctx context.Context
}

func (this *ProjectController) GetCollect1() *simple.JsonResult {
	user := services.UserTokenService.GetCurrent(this.Ctx)
	if user == nil || user.Id != 1 {
		return simple.JsonErrorMsg("无权限")
	}
	go func() {
		services.ProjectService.StartCollect()
	}()
	return simple.JsonSuccess()
}

func (this *ProjectController) GetCollect2() *simple.JsonResult {
	user := services.UserTokenService.GetCurrent(this.Ctx)
	if user == nil || user.Id != 1 {
		return simple.JsonErrorMsg("无权限")
	}
	go func() {
		for i := 48; i >= 1; i-- {
			collect.Page(i, func(url string) {
				p := collect.CollectProject(url)
				if p != nil {
					temp := services.ProjectService.Take("name = ?", p.Name)
					if temp == nil {
						_ = services.ProjectService.Create(&model.Project{
							UserId:      user.Id,
							Name:        p.Name,
							Url:         p.Url,
							DocUrl:      p.DocUrl,
							DownloadUrl: p.DownloadUrl,
							Description: p.Description,
							Content:     p.Content,
							CreateTime:  simple.NowTimestamp(),
						})
					}
				}
			})
		}
	}()
	return simple.JsonSuccess()
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
