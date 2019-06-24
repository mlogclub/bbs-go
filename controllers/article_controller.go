package controllers

import (
	"github.com/mlogclub/mlog/services/cache"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/mlogclub/simple"

	"github.com/mlogclub/mlog/controllers/render"
	"github.com/mlogclub/mlog/model"
	"github.com/mlogclub/mlog/services"
	"github.com/mlogclub/mlog/services/collect"
	"github.com/mlogclub/mlog/utils"

	"github.com/kataras/iris"
)

type ArticleController struct {
	Ctx                   iris.Context
	TagService            *services.TagService
	ArticleService        *services.ArticleService
	UserArticleTagService *services.UserArticleTagService
	FavoriteService       *services.FavoriteService
}

// 文章详情页
func (this *ArticleController) GetBy(articleId int64) {
	article := this.ArticleService.Get(articleId)
	if article == nil || article.Status != model.ArticleStatusPublished {
		this.Ctx.StatusCode(404)
		return
	}

	relatedArticles := this.ArticleService.GetRelatedArticles(articleId)
	newestArticles := this.ArticleService.GetUserNewestArticles(article.UserId)

	tagIds := cache.ArticleTagCache.Get(articleId)
	tags := cache.TagCache.GetList(tagIds)
	var keywords []string
	if len(tags) > 0 {
		for _, tag := range tags {
			keywords = append(keywords, tag.Name)
		}
	}

	render.View(this.Ctx, "article/detail.html", iris.Map{
		utils.GlobalFieldSiteTitle:       article.Title,
		utils.GlobalFieldSiteKeywords:    strings.Join(keywords, ","),
		utils.GlobalFieldSiteDescription: article.Summary,
		"CurrentCategoryId":              0,
		"CurrentTagId":                   0,
		"Tags":                           tags,
		"CommentEntityType":              model.EntityTypeArticle,
		"CommentEntityId":                article.Id,
		"Article":                        render.BuildArticle(article),
		"RelatedArticles":                render.BuildArticles(relatedArticles),
		"NewestArticles":                 render.BuildArticles(newestArticles),
	})
}

// 创建文章页面
func (this *ArticleController) GetCreate() {
	user := utils.GetCurrentUser(this.Ctx)
	if user == nil {
		this.Ctx.Redirect("/user/signin?redirectUrl=/article/create", iris.StatusTemporaryRedirect)
		return
	}

	tags := this.UserArticleTagService.GetUserTags(user.Id)

	render.View(this.Ctx, "article/create.html", iris.Map{
		"Tags": render.BuildTags(tags),
	})
}

// 创建文章
func (this *ArticleController) PostCreate() *simple.JsonResult {
	currentUser := utils.GetCurrentUser(this.Ctx)
	if currentUser == nil {
		return simple.Error(simple.ErrorNotLogin)
	}

	var (
		tagId   = this.Ctx.PostValueInt64Default("tagId", 0)
		title   = this.Ctx.PostValue("title")
		summary = this.Ctx.PostValue("summary")
		content = this.Ctx.PostValue("content")
	)
	if tagId <= 0 {
		return simple.ErrorMsg("请选择标签")
	}
	if len(title) == 0 {
		return simple.ErrorMsg("请输入标题")
	}
	if len(content) == 0 {
		return simple.ErrorMsg("请填写文章内容")
	}

	tag := this.TagService.Get(tagId)
	if tag == nil {
		return simple.ErrorMsg("标签不存在")
	}

	article, err := this.ArticleService.Publish(currentUser.Id, title, summary, content, model.ArticleContentTypeMarkdown, 0, []int64{tagId}, "")
	if err != nil {
		return simple.ErrorMsg(err.Error())
	}
	return simple.NewEmptyRspBuilder().Put("articleId", article.Id).JsonResult()
}

// 编辑文章页面
func (this *ArticleController) GetEditBy(articleId int64) {
	currentUser := utils.GetCurrentUser(this.Ctx)
	if currentUser == nil {
		this.Ctx.Redirect("/user/signin?redirectUrl=/article/edit/"+strconv.FormatInt(articleId, 10), iris.StatusTemporaryRedirect)
		return
	}
	article := this.ArticleService.Get(articleId)
	if article == nil || article.Status == model.ArticleStatusDeleted {
		this.Ctx.StatusCode(404)
		return
	}

	if article.UserId != currentUser.Id {
		this.Ctx.StatusCode(403)
		return
	}

	tags := this.TagService.GetTags()

	render.View(this.Ctx, "article/edit.html", iris.Map{
		"Tags":    tags,
		"Article": article,
	})
}

// 编辑文章
func (this *ArticleController) PostEdit() *simple.JsonResult {
	currentUser := utils.GetCurrentUser(this.Ctx)
	if currentUser == nil {
		return simple.Error(simple.ErrorNotLogin)
	}

	var (
		articleId = this.Ctx.PostValueInt64Default("id", 0)
		tagId     = this.Ctx.PostValueInt64Default("tagId", 0)
		title     = this.Ctx.PostValue("title")
		content   = this.Ctx.PostValue("content")
	)
	if tagId <= 0 {
		return simple.ErrorMsg("请选择标签")
	}
	if len(title) == 0 {
		return simple.ErrorMsg("请输入标题")
	}
	if len(content) == 0 {
		return simple.ErrorMsg("请填写文章内容")
	}

	tag := this.TagService.Get(tagId)
	if tag == nil {
		return simple.ErrorMsg("标签不存在")
	}

	article := this.ArticleService.Get(articleId)
	if article == nil || article.Status == model.ArticleStatusDeleted {
		return simple.ErrorMsg("文章不存在")
	}

	if article.UserId != currentUser.Id {
		return simple.ErrorMsg("无权限")
	}

	article.Title = title
	article.Content = content
	err := this.ArticleService.Update(article)
	if err != nil {
		return simple.ErrorMsg(err.Error())
	}
	return simple.NewEmptyRspBuilder().Put("articleId", article.Id).JsonResult()
}

// 删除文章
func (this *ArticleController) PostDeleteBy(articleId int64) *simple.JsonResult {
	user := utils.GetCurrentUser(this.Ctx)
	if user == nil {
		return simple.Error(simple.ErrorNotLogin)
	}

	article := this.ArticleService.Get(articleId)
	if article == nil || article.Status == model.ArticleStatusDeleted {
		return simple.ErrorMsg("文章不存在")
	}

	if article.UserId != user.Id {
		return simple.ErrorMsg("无权限")
	}

	err := this.ArticleService.Delete(articleId)
	if err != nil {
		return simple.ErrorMsg(err.Error())
	}
	return simple.Success()
}

// 收藏文章
func (this *ArticleController) PostFavoriteBy(articleId int64) *simple.JsonResult {
	user := utils.GetCurrentUser(this.Ctx)
	if user == nil {
		return simple.Error(simple.ErrorNotLogin)
	}
	err := this.FavoriteService.AddArticleFavorite(user.Id, articleId)
	if err != nil {
		return simple.ErrorMsg(err.Error())
	}
	return simple.Success()
}

// 微信采集发布接口
func (this *ArticleController) PostWxpublish() *simple.JsonResult {
	token := this.Ctx.FormValue("token")
	data, err := ioutil.ReadFile("/data/publish_token")
	if err != nil {
		return simple.ErrorMsg("ReadToken error: " + err.Error())
	}
	if token != string(data) {
		return simple.ErrorMsg("Token invalidate")
	}
	article := &collect.WxArticle{}
	err := this.Ctx.ReadJSON(article)
	if err != nil {
		return simple.ErrorMsg(err.Error())
	}
	t, err := collect.NewWxbotApi().Publish(article)
	if err != nil {
		return simple.ErrorMsg(err.Error())
	}
	return simple.NewEmptyRspBuilder().Put("id", t.Id).JsonResult()
}

// 提交百度网址
func (this *ArticleController) GetBaidu() *simple.JsonResult {
	go func() {
		this.ArticleService.Scan(func(articles []model.Article) bool {
			if articles != nil {
				var urls []string
				for _, article := range articles {
					urls = append(urls, utils.BuildArticleUrl(article.Id))
				}
				utils.BaiduUrlPush(urls)
			}
			return true
		})
	}()
	return simple.Success()
}

// 生成每日分享
func (this *ArticleController) GetDaily() *simple.JsonResult {
	services.NewArticleService().CreateDailyShare("M-LOG码农分享", "", []int64{
		177, 79, 105, 115, 197, 88, 29, 171, 60, 53, 128, 143, 20,
	})
	return simple.Success()
}

// 文章列表
func GetArticles(ctx iris.Context) {
	page := ctx.Params().GetIntDefault("page", 1)

	categories := cache.CategoryCache.GetAllCategories()
	activeUsers := cache.UserCache.GetActiveUsers()

	articles, paging := services.ArticleServiceInstance.Query(simple.NewParamQueries(ctx).
		Eq("status", model.ArticleStatusPublished).
		NotEq("category_id", 4).
		Page(page, 20).Desc("id"))

	render.View(ctx, "article/index.html", iris.Map{
		utils.GlobalFieldSiteTitle: "文章",
		"Categories":               categories,
		"Articles":                 render.BuildArticles(articles),
		"ActiveUsers":              render.BuildUsers(activeUsers),
		"Page":                     paging,
		"PrePageUrl":               utils.BuildArticlesUrl(page - 1),
		"NextPageUrl":              utils.BuildArticlesUrl(page + 1),
	})
}

// 分类文章列表
func GetCategoryArticles(ctx iris.Context) {
	categoryId := ctx.Params().GetInt64Default("categoryId", 0)
	page := ctx.Params().GetIntDefault("page", 1)

	categories := cache.CategoryCache.GetAllCategories()
	activeUsers := cache.UserCache.GetActiveUsers()
	category := services.CategoryServiceInstance.Get(categoryId)

	title := "文章"
	if category != nil {
		title = category.Name + " - " + title
	}

	articles, paging := services.ArticleServiceInstance.Query(simple.NewParamQueries(ctx).
		Eq("category_id", categoryId).
		Eq("status", model.ArticleStatusPublished).
		Page(page, 20).Desc("id"))

	render.View(ctx, "article/index.html", iris.Map{
		utils.GlobalFieldSiteTitle: title,
		"Categories":               categories,
		"Articles":                 render.BuildArticles(articles),
		"ActiveUsers":              render.BuildUsers(activeUsers),
		"Page":                     paging,
		"PrePageUrl":               utils.BuildCategoryArticlesUrl(categoryId, page-1),
		"NextPageUrl":              utils.BuildCategoryArticlesUrl(categoryId, page+1),
	})
}

// 标签文章列表
func GetTagArticles(ctx iris.Context) {
	tagId := ctx.Params().GetInt64Default("tagId", 0)
	page := ctx.Params().GetIntDefault("page", 1)

	categories := cache.CategoryCache.GetAllCategories()
	activeUsers := cache.UserCache.GetActiveUsers()
	tag := services.TagServiceInstance.Get(tagId)

	title := "文章"
	if tag != nil {
		title = tag.Name + " - " + title
	}

	articles, paging := services.ArticleServiceInstance.GetTagArticles(tagId, page)

	render.View(ctx, "article/index.html", iris.Map{
		utils.GlobalFieldSiteTitle: title,
		"Categories":               categories,
		"Articles":                 render.BuildArticles(articles),
		"ActiveUsers":              render.BuildUsers(activeUsers),
		"Page":                     paging,
		"PrePageUrl":               utils.BuildTagArticlesUrl(tagId, page-1),
		"NextPageUrl":              utils.BuildTagArticlesUrl(tagId, page+1),
	})
}
