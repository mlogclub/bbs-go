package web

import (
	"github.com/mlogclub/mlog/services/cache"
	"github.com/mlogclub/mlog/utils/session"
	"github.com/sirupsen/logrus"
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
	Ctx iris.Context
}

// 文章详情页
func (this *ArticleController) GetBy(articleId int64) {
	article := services.ArticleService.Get(articleId)
	if article == nil || article.Status != model.ArticleStatusPublished {
		this.Ctx.StatusCode(404)
		return
	}

	relatedArticles := services.ArticleService.GetRelatedArticles(articleId)
	newestArticles := services.ArticleService.GetUserNewestArticles(article.UserId)

	tagIds := cache.ArticleTagCache.Get(articleId)
	tags := cache.TagCache.GetList(tagIds)
	var keywords []string
	if len(tags) > 0 {
		for _, tag := range tags {
			keywords = append(keywords, tag.Name)
		}
	}

	render.View(this.Ctx, "article/detail.html", iris.Map{
		model.TplSiteTitle:       article.Title,
		model.TplSiteKeywords:    strings.Join(keywords, ","),
		model.TplSiteDescription: article.Summary,
		"CurrentCategoryId":      0,
		"CurrentTagId":           0,
		"Tags":                   tags,
		"CommentEntityType":      model.EntityTypeArticle,
		"CommentEntityId":        article.Id,
		"Article":                render.BuildArticle(article),
		"RelatedArticles":        render.BuildArticles(relatedArticles),
		"NewestArticles":         render.BuildArticles(newestArticles),
	})
}

// 创建文章页面
func (this *ArticleController) GetCreate() {
	user := session.GetCurrentUser(this.Ctx)
	if user == nil {
		this.Ctx.Redirect("/user/signin?redirectUrl=/article/create", iris.StatusTemporaryRedirect)
		return
	}
	render.View(this.Ctx, "article/create.html", iris.Map{})
}

// 创建文章
func (this *ArticleController) PostCreate() *simple.JsonResult {
	currentUser := session.GetCurrentUser(this.Ctx)
	if currentUser == nil {
		return simple.Error(simple.ErrorNotLogin)
	}

	var (
		tags    = simple.FormValueStringArray(this.Ctx, "tags")
		title   = this.Ctx.PostValue("title")
		summary = this.Ctx.PostValue("summary")
		content = this.Ctx.PostValue("content")
	)

	article, err := services.ArticleService.Publish(currentUser.Id, title, summary, content,
		model.ArticleContentTypeMarkdown, 0, tags, "", false)
	if err != nil {
		return simple.ErrorMsg(err.Error())
	}
	return simple.NewEmptyRspBuilder().Put("articleId", article.Id).JsonResult()
}

// 编辑文章页面
func (this *ArticleController) GetEditBy(articleId int64) {
	currentUser := session.GetCurrentUser(this.Ctx)
	if currentUser == nil {
		this.Ctx.Redirect("/user/signin?redirectUrl=/article/edit/"+strconv.FormatInt(articleId, 10), iris.StatusTemporaryRedirect)
		return
	}
	article := services.ArticleService.Get(articleId)
	if article == nil || article.Status == model.ArticleStatusDeleted {
		this.Ctx.StatusCode(404)
		return
	}

	if article.UserId != currentUser.Id {
		this.Ctx.StatusCode(403)
		return
	}

	tags := services.ArticleService.GetArticleTags(articleId)
	var tagNames []string
	if len(tags) > 0 {
		for _, tag := range tags {
			tagNames = append(tagNames, tag.Name)
		}
	}

	render.View(this.Ctx, "article/edit.html", iris.Map{
		"Article": iris.Map{
			"ArticleId": article.Id,
			"Title":     article.Title,
			"Content":   article.Content,
			"Tags":      tagNames,
		},
	})
}

// 编辑文章
func (this *ArticleController) PostEdit() *simple.JsonResult {
	currentUser := session.GetCurrentUser(this.Ctx)
	if currentUser == nil {
		return simple.Error(simple.ErrorNotLogin)
	}

	var (
		articleId = this.Ctx.PostValueInt64Default("id", 0)
		tags      = simple.FormValueStringArray(this.Ctx, "tags")
		title     = this.Ctx.PostValue("title")
		content   = this.Ctx.PostValue("content")
	)

	article := services.ArticleService.Get(articleId)
	if article == nil || article.Status == model.ArticleStatusDeleted {
		return simple.ErrorMsg("文章不存在")
	}

	if article.UserId != currentUser.Id {
		return simple.ErrorMsg("无权限")
	}

	err := services.ArticleService.Edit(articleId, tags, title, content)
	if err != nil {
		return simple.Error(err)
	}
	return simple.NewEmptyRspBuilder().Put("articleId", article.Id).JsonResult()
}

// 删除文章
func (this *ArticleController) PostDeleteBy(articleId int64) *simple.JsonResult {
	user := session.GetCurrentUser(this.Ctx)
	if user == nil {
		return simple.Error(simple.ErrorNotLogin)
	}

	article := services.ArticleService.Get(articleId)
	if article == nil || article.Status == model.ArticleStatusDeleted {
		return simple.ErrorMsg("文章不存在")
	}

	if article.UserId != user.Id {
		return simple.ErrorMsg("无权限")
	}

	err := services.ArticleService.Delete(articleId)
	if err != nil {
		return simple.ErrorMsg(err.Error())
	}
	return simple.Success()
}

// 收藏文章
func (this *ArticleController) PostFavoriteBy(articleId int64) *simple.JsonResult {
	user := session.GetCurrentUser(this.Ctx)
	if user == nil {
		return simple.Error(simple.ErrorNotLogin)
	}
	err := services.FavoriteService.AddArticleFavorite(user.Id, articleId)
	if err != nil {
		return simple.ErrorMsg(err.Error())
	}
	return simple.Success()
}

// 跳转到文章的原始链接
func (this *ArticleController) GetRedirectBy(articleId int64) {
	article := services.ArticleService.Get(articleId)
	if article == nil || article.Status != model.ArticleStatusPublished {
		this.Ctx.StatusCode(404)
		return
	}
	if article.Share && len(article.SourceUrl) > 0 {
		this.Ctx.Redirect(article.SourceUrl, iris.StatusFound)
	} else {
		this.Ctx.Redirect("/article/"+strconv.FormatInt(articleId, 10), iris.StatusFound)
	}
}

// 微信采集发布接口
func (this *ArticleController) PostWxpublish() *simple.JsonResult {
	token := this.Ctx.FormValue("token")
	data, err := ioutil.ReadFile("/data/publish_token")
	if err != nil {
		return simple.ErrorMsg("ReadToken error: " + err.Error())
	}
	token2 := strings.TrimSpace(string(data))
	logrus.Info("token: " + token + ", token2: " + token2)
	if token != token2 {
		return simple.ErrorMsg("Token invalidate")
	}
	article := &collect.WxArticle{}
	err = this.Ctx.ReadJSON(article)
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
		services.ArticleService.Scan(func(articles []model.Article) bool {
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

// 文章列表
func GetArticles(ctx iris.Context) {
	page := ctx.Params().GetIntDefault("page", 1)

	categories := cache.CategoryCache.GetAllCategories()
	activeUsers := cache.UserCache.GetActiveUsers()
	activeTags := cache.TagCache.GetActiveTags()

	articles, paging := services.ArticleService.Query(simple.NewParamQueries(ctx).
		Eq("status", model.ArticleStatusPublished).
		NotEq("category_id", 4).
		Page(page, 20).Desc("id"))

	render.View(ctx, "article/index.html", iris.Map{
		model.TplSiteTitle: "文章",
		"Categories":       categories,
		"Articles":         render.BuildArticles(articles),
		"ActiveUsers":      render.BuildUsers(activeUsers),
		"ActiveTags":       render.BuildTags(activeTags),
		"Page":             paging,
		"PrePageUrl":       utils.BuildArticlesUrl(page - 1),
		"NextPageUrl":      utils.BuildArticlesUrl(page + 1),
	})
}

// 分类文章列表
func GetCategoryArticles(ctx iris.Context) {
	categoryId := ctx.Params().GetInt64Default("categoryId", 0)
	page := ctx.Params().GetIntDefault("page", 1)

	categories := cache.CategoryCache.GetAllCategories()
	activeUsers := cache.UserCache.GetActiveUsers()
	activeTags := cache.TagCache.GetActiveTags()
	category := services.CategoryService.Get(categoryId)

	title := "文章"
	if category != nil {
		title = category.Name + " - " + title
	}

	articles, paging := services.ArticleService.Query(simple.NewParamQueries(ctx).
		Eq("category_id", categoryId).
		Eq("status", model.ArticleStatusPublished).
		Page(page, 20).Desc("id"))

	render.View(ctx, "article/index.html", iris.Map{
		model.TplSiteTitle: title,
		"Categories":       categories,
		"Articles":         render.BuildArticles(articles),
		"ActiveUsers":      render.BuildUsers(activeUsers),
		"ActiveTags":       render.BuildTags(activeTags),
		"Page":             paging,
		"PrePageUrl":       utils.BuildCategoryArticlesUrl(categoryId, page-1),
		"NextPageUrl":      utils.BuildCategoryArticlesUrl(categoryId, page+1),
	})
}

// 标签文章列表
func GetTagArticles(ctx iris.Context) {
	tagId := ctx.Params().GetInt64Default("tagId", 0)
	page := ctx.Params().GetIntDefault("page", 1)

	categories := cache.CategoryCache.GetAllCategories()
	activeUsers := cache.UserCache.GetActiveUsers()
	activeTags := cache.TagCache.GetActiveTags()
	tag := services.TagService.Get(tagId)

	title := "文章"
	if tag != nil {
		title = tag.Name + " - " + title
	}

	articles, paging := services.ArticleService.GetTagArticles(tagId, page)

	render.View(ctx, "article/index.html", iris.Map{
		model.TplSiteTitle: title,
		"Categories":       categories,
		"Articles":         render.BuildArticles(articles),
		"ActiveUsers":      render.BuildUsers(activeUsers),
		"ActiveTags":       render.BuildTags(activeTags),
		"Page":             paging,
		"PrePageUrl":       utils.BuildTagArticlesUrl(tagId, page-1),
		"NextPageUrl":      utils.BuildTagArticlesUrl(tagId, page+1),
	})
}
