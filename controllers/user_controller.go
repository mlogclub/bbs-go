package controllers

import (
	context2 "context"
	"github.com/mlogclub/mlog/services/cache"
	"strconv"

	"github.com/kataras/iris"
	"github.com/kataras/iris/context"
	"github.com/mlogclub/simple"
	"github.com/sirupsen/logrus"

	"github.com/mlogclub/mlog/controllers/render"
	"github.com/mlogclub/mlog/model"
	"github.com/mlogclub/mlog/services"
	"github.com/mlogclub/mlog/utils"
)

const oauthStateString = "random"

type UserController struct {
	Ctx                   iris.Context
	UserService           *services.UserService
	GithubUserService     *services.GithubUserService
	ArticleService        *services.ArticleService
	UserArticleTagService *services.UserArticleTagService
	FavoriteService       *services.FavoriteService
	MessageService        *services.MessageService
}

// 用户主页
func (this *UserController) GetBy(userId int64) {
	user := cache.UserCache.Get(userId)
	if user == nil {
		this.Ctx.StatusCode(404)
		return
	}

	tab := this.Ctx.FormValueDefault("tab", "articles")

	var articles []model.Article
	var tags []model.Tag
	var messages []model.Message
	var favorites []model.Favorite
	if tab == "articles" {
		articles, _ = this.ArticleService.QueryCnd(simple.NewQueryCnd("user_id = ? and status = ?", userId,
			model.ArticleStatusPublished).Order("id desc").Size(10))
	} else if tab == "tags" {
		tags = this.UserArticleTagService.GetUserTags(userId)
	} else if tab == "messages" {
		messages, _ = this.MessageService.QueryCnd(simple.NewQueryCnd("user_id = ?", userId).Order("id desc").Size(10))
	} else if tab == "favorites" {
		currentUser := utils.GetCurrentUser(this.Ctx)
		if currentUser != nil && currentUser.Id == user.Id {
			favorites, _ = this.FavoriteService.QueryCnd(simple.NewQueryCnd("user_id = ?", userId).Order("id desc").Size(10))
		}
	}

	render.View(this.Ctx, "user/index.html", iris.Map{
		utils.GlobalFieldSiteTitle: user.Nickname,
		"User":                     render.BuildUser(user),
		"Tab":                      tab,
		"Articles":                 render.BuildArticles(articles),
		"Tags":                     render.BuildTags(tags),
		"Messages":                 render.BuildMessages(messages),
		"Favorites":                render.BuildFavorites(favorites),
	})
}

// 编辑资料页面
func (this *UserController) GetEditBy(userId int64) {
	currentUser := utils.GetCurrentUser(this.Ctx)
	if currentUser == nil {
		this.Ctx.Redirect("/user/signin", iris.StatusSeeOther)
		return
	}
	render.View(this.Ctx, "user/edit.html", iris.Map{
		utils.GlobalFieldSiteTitle: currentUser.Nickname + " - 编辑资料",
	})
}

// 提交编辑资料
func (this *UserController) PostEditBy(userId int64) {
	currentUser := utils.GetCurrentUser(this.Ctx)
	if currentUser == nil {
		this.Ctx.Redirect("/user/signin", iris.StatusSeeOther)
		return
	}

	nickname, err := simple.FormValueRequired(this.Ctx, "nickname")
	if err != nil {
		render.View(this.Ctx, "user/edit.html", iris.Map{
			"ErrMsg":                   "昵称不能为空",
			utils.GlobalFieldSiteTitle: currentUser.Nickname + " - 编辑资料",
		})
		return
	}
	avatar, err := simple.FormValueRequired(this.Ctx, "avatar")
	if err != nil {
		render.View(this.Ctx, "user/edit.html", iris.Map{
			"ErrMsg":                   "昵称不能为空",
			utils.GlobalFieldSiteTitle: currentUser.Nickname + " - 编辑资料",
		})
		return
	}
	description := simple.FormValue(this.Ctx, "description")

	_ = this.UserService.Updates(currentUser.Id, map[string]interface{}{
		"nickname":    nickname,
		"avatar":      avatar,
		"description": description,
	})
	cache.UserCache.Invalidate(currentUser.Id)
	this.Ctx.Redirect(utils.BuildUserUrl(currentUser.Id), iris.StatusSeeOther)
}

// 当前登录用户
func (this *UserController) GetCurrent() *simple.JsonResult {
	user := utils.GetCurrentUser(this.Ctx)
	if user == nil {
		return simple.Error(simple.ErrorNotLogin)
	}
	return simple.JsonData(render.BuildUserById(user.Id))
}

// 登录页面
func (this *UserController) GetSignin() {
	user := utils.GetCurrentUser(this.Ctx)
	if user != nil {
		this.Ctx.Redirect("/", iris.StatusSeeOther)
		return
	}
	redirectUrl := simple.FormValue(this.Ctx, "redirectUrl")
	render.View(this.Ctx, "user/signin.html", iris.Map{
		utils.GlobalFieldSiteTitle: "登录",
		"RedirectUrl":              redirectUrl,
	})
}

// 提交登录
func (this *UserController) PostSignin() {
	username := this.Ctx.PostValueTrim("username")
	password := this.Ctx.PostValueTrim("password")

	user, err := this.UserService.SignIn(username, password)
	if err != nil {
		render.View(this.Ctx, "user/signin.html", iris.Map{
			"ErrMsg":   err.Error(),
			"Username": username,
			"Password": password,
		})
		return
	}

	utils.SetCurrentUser(this.Ctx, user)

	redirectUrl := simple.FormValue(this.Ctx, "redirectUrl")
	if len(redirectUrl) > 0 {
		this.Ctx.Redirect(redirectUrl, iris.StatusSeeOther)
	} else {
		this.Ctx.Redirect("/", iris.StatusSeeOther)
	}
}

// 注册
func (this *UserController) GetSignup() {
	user := utils.GetCurrentUser(this.Ctx)
	if user != nil {
		this.Ctx.Redirect("/", iris.StatusSeeOther)
		return
	}
	render.View(this.Ctx, "user/signup.html", iris.Map{
		utils.GlobalFieldSiteTitle: "注册",
	})
}

// 注册
func (this *UserController) PostSignup() {
	username := this.Ctx.PostValueTrim("username")
	email := this.Ctx.PostValueTrim("email")
	password := this.Ctx.PostValueTrim("password")
	rePassword := this.Ctx.PostValueTrim("rePassword")
	nickname := this.Ctx.PostValueTrim("password")

	user, err := this.UserService.SignUp(username, email, password, rePassword, nickname, "")
	if err != nil {
		render.View(this.Ctx, "user/signup.html", iris.Map{
			"ErrMsg":     err.Error(),
			"Username":   username,
			"Email":      email,
			"Password":   password,
			"RePassword": rePassword,
			"Nickname":   nickname,
		})
		return
	}
	utils.SetCurrentUser(this.Ctx, user)

	redirectUrl := simple.FormValue(this.Ctx, "redirectUrl")
	if len(redirectUrl) > 0 {
		this.Ctx.Redirect(redirectUrl, iris.StatusSeeOther)
	} else {
		this.Ctx.Redirect("/", iris.StatusSeeOther)
	}
}

// 退出登录
func (this *UserController) AnySignout() {
	utils.DelCurrentUser(this.Ctx)
	this.Ctx.Redirect(utils.Conf.BaseUrl, iris.StatusSeeOther)
}

// 跳转到Github登录页面
func (this *UserController) GetGithubLogin() {
	url := utils.GithubOauthConfig.AuthCodeURL(oauthStateString)
	this.Ctx.Redirect(url, iris.StatusSeeOther)
}

// Github回调地址
func (this *UserController) GetGithubCallback() {
	state := this.Ctx.FormValue("state")
	if state != oauthStateString {
		logrus.Errorf("invalid oauth state, expected '%s', got '%s'", oauthStateString, state)
		this.Ctx.Redirect("/", iris.StatusSeeOther)
		return
	}

	code := this.Ctx.FormValue("code")
	token, err := utils.GithubOauthConfig.Exchange(context2.TODO(), code)
	if err != nil {
		logrus.Errorf("Code exchange failed with '%s'", err)
		this.Ctx.Redirect("/", iris.StatusSeeOther)
		return
	}

	third := utils.GetGithubUserinfo(token.AccessToken)
	githubUser := this.GithubUserService.GetByGithubId(third.Id)

	var user *model.User
	if githubUser != nil {
		user = cache.UserCache.Get(githubUser.UserId)
	} else {
		githubUser = &model.GithubUser{
			GithubId:   third.Id,
			Login:      third.Login,
			NodeId:     third.NodeId,
			AvatarUrl:  third.AvatarUrl,
			Url:        third.Url,
			HtmlUrl:    third.HtmlUrl,
			Email:      third.Email,
			Name:       third.Name,
			CreateTime: simple.NowTimestamp(),
			UpdateTime: simple.NowTimestamp(),
		}
		err := this.GithubUserService.Create(githubUser)
		if err != nil {
			logrus.Error(err)
		}
	}

	if user != nil {
		utils.SetCurrentUser(this.Ctx, user)
		this.Ctx.Redirect("/", iris.StatusSeeOther)
	} else {
		this.Ctx.Redirect("/user/github/bind?id="+strconv.FormatInt(githubUser.Id, 10), iris.StatusSeeOther)
	}
}

// Github绑定页面
func (this *UserController) GetGithubBind() {
	githubUserId, err := simple.FormValueInt64(this.Ctx, "id")
	if err != nil {
		render.View(this.Ctx, "/user/github_bind.html", iris.Map{
			"ErrMsg": err.Error(),
		})
		return
	}

	githubUser := this.GithubUserService.Get(githubUserId)
	if githubUser == nil {
		render.View(this.Ctx, "/user/github_bind.html", iris.Map{
			"ErrMsg": "数据错误",
		})
		return
	}

	render.View(this.Ctx, "/user/github_bind.html", iris.Map{
		"GithubId": githubUser.Id,
		"Email":    githubUser.Email,
		"BindType": "login",
	})
}

// Github提交绑定
func (this *UserController) PostGithubBind() {
	bindType := this.Ctx.PostValueTrim("bindType")
	githubId, err := this.Ctx.PostValueInt64("githubId")
	username := this.Ctx.PostValueTrim("username")
	email := this.Ctx.PostValueTrim("email")
	password := this.Ctx.PostValueTrim("password")
	rePassword := this.Ctx.PostValueTrim("rePassword")
	nickname := this.Ctx.PostValueTrim("nickname")

	if err != nil {
		render.View(this.Ctx, "/user/github_bind.html", iris.Map{
			"ErrMsg":     err.Error(),
			"BindType":   bindType,
			"GithubId":   githubId,
			"Username":   username,
			"Email":      email,
			"Password":   password,
			"RePassword": rePassword,
			"Nickname":   nickname,
		})
		return
	}

	user, err := this.UserService.Bind(githubId, bindType, username, email, password, rePassword, nickname)
	if err != nil {
		render.View(this.Ctx, "/user/github_bind.html", iris.Map{
			"ErrMsg":     err.Error(),
			"BindType":   bindType,
			"GithubId":   githubId,
			"Username":   username,
			"Email":      email,
			"Password":   password,
			"RePassword": rePassword,
			"Nickname":   nickname,
		})
		return
	}
	utils.SetCurrentUser(this.Ctx, user)

	redirectUrl := simple.FormValue(this.Ctx, "redirectUrl")
	if len(redirectUrl) > 0 {
		this.Ctx.Redirect(redirectUrl, iris.StatusSeeOther)
	} else {
		this.Ctx.Redirect("/", iris.StatusSeeOther)
	}
}

// 未读消息数量
func (this *UserController) GetMsgcount() *simple.JsonResult {
	user := utils.GetCurrentUser(this.Ctx)
	var count int64 = 0
	if user != nil {
		count = this.MessageService.GetUnReadCount(user.Id)
	}
	return simple.NewEmptyRspBuilder().Put("count", count).JsonResult()
}

// 用户中心 - 用户所有的文章
func GetUserArticles(ctx context.Context) {
	userId := ctx.Params().GetInt64Default("userId", 0)
	page := ctx.Params().GetIntDefault("page", 1)
	user := cache.UserCache.Get(userId)
	if user == nil {
		ctx.StatusCode(404)
		return
	}

	articles, paging := services.ArticleServiceInstance.Query(simple.NewParamQueries(ctx).
		Eq("user_id", userId).
		Eq("status", model.ArticleStatusPublished).
		Page(page, 20).Desc("id"))

	render.View(ctx, "user/articles.html", iris.Map{
		"User":                     render.BuildUser(user),
		"Articles":                 render.BuildArticles(articles),
		"Page":                     paging,
		"PrePageUrl":               utils.BuildUserArticlesUrl(userId, page-1),
		"NextPageUrl":              utils.BuildUserArticlesUrl(userId, page+1),
		utils.GlobalFieldSiteTitle: user.Nickname + " - 文章列表",
	})
}

// 用户中心 - 用户所有的标签
func GetUserTags(ctx context.Context) {
	userId := ctx.Params().GetInt64Default("userId", 0)
	page := ctx.Params().GetIntDefault("page", 1)
	user := cache.UserCache.Get(userId)

	list, paging := services.UserArticleTagServiceInstance.Query(simple.NewParamQueries(ctx).Eq("user_id", userId).
		Page(page, 20).Desc("id"))

	var tags []model.Tag
	if len(list) > 0 {
		for _, v := range list {
			tags = append(tags, *cache.TagCache.Get(v.TagId))
		}
	}

	render.View(ctx, "user/tags.html", iris.Map{
		"User":                     render.BuildUser(user),
		"Tags":                     render.BuildTags(tags),
		"Page":                     paging,
		utils.GlobalFieldSiteTitle: user.Nickname + " - 标签列表",
	})
}

// 用户中心 - 用户所有的搜藏列表
func GetUserFavorites(ctx context.Context) {
	user := utils.GetCurrentUser(ctx)
	userId := ctx.Params().GetInt64Default("userId", 0)
	page := ctx.Params().GetIntDefault("page", 1)

	// 查看收藏必须登录
	if user == nil {
		ctx.Redirect("/user/signin", iris.StatusSeeOther)
		return
	}
	// 只能查看自己的收藏
	if userId != user.Id {
		ctx.StatusCode(403)
		return
	}

	favorites, paging := services.FavoriteServiceInstance.Query(simple.NewParamQueries(ctx).
		Eq("user_id", userId).
		Page(page, 20).Desc("id"))

	render.View(ctx, "user/favorites.html", iris.Map{
		"User":                     render.BuildUser(user),
		"Favorites":                render.BuildFavorites(favorites),
		"Page":                     paging,
		"PrePageUrl":               utils.BuildUserFavoritesUrl(userId, page-1),
		"NextPageUrl":              utils.BuildUserFavoritesUrl(userId, page+1),
		utils.GlobalFieldSiteTitle: user.Nickname + " - 收藏列表",
	})
}

// 用户中心 - 消息列表
func GetUserMessages(ctx context.Context) {
	user := utils.GetCurrentUser(ctx)
	userId := ctx.Params().GetInt64Default("userId", 0)
	page := ctx.Params().GetIntDefault("page", 1)

	// 用户必须登录
	if user == nil {
		ctx.Redirect("/user/signin", iris.StatusSeeOther)
		return
	}

	// 只能查看自己的收藏
	if userId != user.Id {
		ctx.StatusCode(403)
		return
	}

	// 查询列表
	messages, paging := services.MessageServiceInstance.Query(simple.NewParamQueries(ctx).
		Eq("user_id", user.Id).Page(page, 20).Desc("id"))

	// 全部标记为已读
	services.MessageServiceInstance.MarkReadAll(user.Id)

	render.View(ctx, "user/messages.html", iris.Map{
		"User":                     render.BuildUser(user),
		"Messages":                 render.BuildMessages(messages),
		"Page":                     paging,
		"PrePageUrl":               utils.BuildMessagesUrl(page - 1),
		"NextPageUrl":              utils.BuildMessagesUrl(page + 1),
		utils.GlobalFieldSiteTitle: user.Nickname + " - 消息",
	})
}
