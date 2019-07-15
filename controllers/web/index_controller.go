package web

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"github.com/mlogclub/mlog/services"
	"github.com/mlogclub/mlog/services/cache"
	"github.com/mlogclub/simple"

	"github.com/mlogclub/mlog/controllers/render"
	"github.com/mlogclub/mlog/model"
	"github.com/mlogclub/mlog/utils"
)

type IndexController struct {
	Ctx iris.Context
}

func (this *IndexController) Any() mvc.View {
	categories := cache.CategoryCache.GetAllCategories()
	activeUsers := cache.UserCache.GetActiveUsers()
	activeTags := cache.TagCache.GetActiveTags()

	// 存在缓存从缓存里面取
	articles := cache.ArticleCache.GetIndexList()

	topics, _ := services.TopicService.QueryCnd(simple.NewQueryCnd("status = ?", model.TopicStatusOk).Order("id desc").Size(10))

	return mvc.View{
		Name: "index.html",
		Data: iris.Map{
			utils.GlobalFieldSiteDescription: "M-LOG分享",
			utils.GlobalFieldSiteKeywords:    "Go中国，Golang, Golang学习,Golang教程,Golang社区,Go基金会,Go语言中文网,Go,Go语言,主题,资源,文章,图书,开源项目,M-LOG",
			"Categories":                     categories,
			"Articles":                       render.BuildArticles(articles),
			"Topics":                         render.BuildTopics(topics),
			"ActiveUsers":                    render.BuildUsers(activeUsers),
			"ActiveTags":                     render.BuildTags(activeTags),
		},
	}
}

func (this *IndexController) GetAbout() mvc.View {
	return mvc.View{
		Name: "about.html",
	}
}
