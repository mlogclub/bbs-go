package web

import (
	"github.com/mlogclub/mlog/services/cache"
	"github.com/mlogclub/mlog/utils/config"
	"github.com/mlogclub/mlog/utils/github"
	"github.com/mlogclub/mlog/utils/session"
	"strconv"
	"strings"

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
	Ctx iris.Context
}

// 用户主页
func (this *UserController) GetBy(userId int64) {
	user := cache.UserCache.Get(userId)
	if user == nil {
		this.Ctx.StatusCode(404)
		return
	}

	currentUser := session.GetCurrentUser(this.Ctx)

	tab := this.Ctx.FormValueDefault("tab", "articles")
	owner := currentUser != nil && currentUser.Id == user.Id // 是否是主人态

	var articles []model.Article
	var topics []model.Topic
	var tags []model.Tag
	var messages []model.Message
	var favorites []model.Favorite
	if tab == "articles" {
		articles, _ = services.ArticleService.QueryCnd(simple.NewQueryCnd("user_id = ? and status = ?", userId,
			model.ArticleStatusPublished).Order("id desc").Size(10))
	} else if tab == "topics" {
		topics, _ = services.TopicService.QueryCnd(simple.NewQueryCnd("user_id = ? and status = ?", userId, model.TopicStatusOk).
			Order("id desc").Size(10))
	} else if tab == "tags" && owner {
		tags = services.UserArticleTagService.GetUserTags(userId)
	} else if tab == "messages" && owner {
		messages, _ = services.MessageService.QueryCnd(simple.NewQueryCnd("user_id = ?", userId).Order("id desc").Size(10))
	} else if tab == "favorites" && owner {
		currentUser := session.GetCurrentUser(this.Ctx)
		if currentUser != nil && currentUser.Id == user.Id {
			favorites, _ = services.FavoriteService.QueryCnd(simple.NewQueryCnd("user_id = ?", userId).Order("id desc").Size(10))
		}
	}

	render.View(this.Ctx, "user/index.html", iris.Map{
		utils.GlobalFieldSiteTitle: user.Nickname,
		"User":                     render.BuildUser(user),
		"Tab":                      tab,
		"Articles":                 render.BuildArticles(articles),
		"Topics":                   render.BuildTopics(topics),
		"Tags":                     render.BuildTags(tags),
		"Messages":                 render.BuildMessages(messages),
		"Favorites":                render.BuildFavorites(favorites),
	})
}

// 编辑资料页面
func (this *UserController) GetEditBy(userId int64) {
	currentUser := session.GetCurrentUser(this.Ctx)
	if currentUser == nil {
		this.Ctx.Redirect("/user/signin", iris.StatusSeeOther)
		return
	}
	render.View(this.Ctx, "user/edit.html", iris.Map{
		utils.GlobalFieldSiteTitle: currentUser.Nickname + " - 编辑资料",
		"UserId":                   currentUser.Id,
		"Username":                 currentUser.Username,
		"Email":                    currentUser.Email,
		"Nickname":                 currentUser.Nickname,
		"Avatar":                   currentUser.Avatar,
		"Description":              currentUser.Description,
	})
}

// 提交编辑资料
func (this *UserController) PostEditBy(userId int64) {
	currentUser := session.GetCurrentUser(this.Ctx)
	if currentUser == nil {
		this.Ctx.Redirect("/user/signin", iris.StatusSeeOther)
		return
	}

	nickname := strings.TrimSpace(simple.FormValue(this.Ctx, "nickname"))
	avatar := strings.TrimSpace(simple.FormValue(this.Ctx, "avatar"))
	description := simple.FormValue(this.Ctx, "description")

	if len(nickname) == 0 {
		render.View(this.Ctx, "user/edit.html", iris.Map{
			"ErrMsg":                   "昵称不能为空",
			utils.GlobalFieldSiteTitle: currentUser.Nickname + " - 编辑资料",
			"UserId":                   currentUser.Id,
			"Username":                 currentUser.Username,
			"Email":                    currentUser.Email,
			"Nickname":                 nickname,
			"Avatar":                   avatar,
			"Description":              description,
		})
		return
	}
	if len(avatar) == 0 {
		render.View(this.Ctx, "user/edit.html", iris.Map{
			"ErrMsg":                   "头像不能为空",
			utils.GlobalFieldSiteTitle: currentUser.Nickname + " - 编辑资料",
			"UserId":                   currentUser.Id,
			"Username":                 currentUser.Username,
			"Email":                    currentUser.Email,
			"Nickname":                 nickname,
			"Avatar":                   avatar,
			"Description":              description,
		})
		return
	}

	_ = services.UserService.Updates(currentUser.Id, map[string]interface{}{
		"nickname":    nickname,
		"avatar":      avatar,
		"description": description,
	})
	this.Ctx.Redirect("/user/edit/"+strconv.FormatInt(currentUser.Id, 10), iris.StatusSeeOther)
}

// 当前登录用户
func (this *UserController) GetCurrent() *simple.JsonResult {
	user := session.GetCurrentUser(this.Ctx)
	if user == nil {
		return simple.Error(simple.ErrorNotLogin)
	}
	return simple.JsonData(render.BuildUserById(user.Id))
}

// 登录页面
func (this *UserController) GetSignin() {
	user := session.GetCurrentUser(this.Ctx)
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

	user, err := services.UserService.SignIn(username, password)
	if err != nil {
		render.View(this.Ctx, "user/signin.html", iris.Map{
			"ErrMsg":   err.Error(),
			"Username": username,
			"Password": password,
		})
		return
	}

	session.SetCurrentUser(this.Ctx, user.Id)

	redirectUrl := simple.FormValue(this.Ctx, "redirectUrl")
	if len(redirectUrl) > 0 {
		this.Ctx.Redirect(redirectUrl, iris.StatusSeeOther)
	} else {
		this.Ctx.Redirect("/", iris.StatusSeeOther)
	}
}

// 注册
func (this *UserController) GetSignup() {
	user := session.GetCurrentUser(this.Ctx)
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
	nickname := this.Ctx.PostValueTrim("nickname")

	user, err := services.UserService.SignUp(username, email, password, rePassword, nickname, "")
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
	session.SetCurrentUser(this.Ctx, user.Id)

	redirectUrl := simple.FormValue(this.Ctx, "redirectUrl")
	if len(redirectUrl) > 0 {
		this.Ctx.Redirect(redirectUrl, iris.StatusSeeOther)
	} else {
		this.Ctx.Redirect("/", iris.StatusSeeOther)
	}
}

// 退出登录
func (this *UserController) AnySignout() {
	session.DelCurrentUser(this.Ctx)
	this.Ctx.Redirect(config.Conf.BaseUrl, iris.StatusSeeOther)
}

// 跳转到Github登录页面
func (this *UserController) GetGithubLogin() {
	url := github.OauthConfig.AuthCodeURL(oauthStateString)
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
	githubUser, err := services.GithubUserService.GetGithubUser(code)

	if err != nil {
		logrus.Errorf("Code exchange failed with '%s'", err)
		this.Ctx.StatusCode(500)
		return
	}

	user, codeErr := services.UserService.SignInByGithub(githubUser)
	if codeErr != nil {
		if codeErr.Code == utils.ErrorCodeUserNameExists {
			this.Ctx.Redirect("/user/github/bind?id="+strconv.FormatInt(githubUser.Id, 10), iris.StatusSeeOther)
		} else if codeErr.Code == utils.ErrorCodeEmailExists {
			this.Ctx.Redirect("/user/github/bind?id="+strconv.FormatInt(githubUser.Id, 10), iris.StatusSeeOther)
		} else {
			logrus.Errorf("Code exchange failed with '%s'", codeErr)
			this.Ctx.StatusCode(500)
		}
	} else { // 直接登录
		session.SetCurrentUser(this.Ctx, user.Id)
		this.Ctx.Redirect("/", iris.StatusSeeOther)
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

	githubUser := services.GithubUserService.Get(githubUserId)
	if githubUser == nil {
		render.View(this.Ctx, "/user/github_bind.html", iris.Map{
			"ErrMsg": "数据错误",
		})
		return
	}

	render.View(this.Ctx, "/user/github_bind.html", iris.Map{
		"GithubId": githubUser.Id,
		"Username": githubUser.Login,
		"Email":    githubUser.Email,
		"Nickname": githubUser.Name,
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

	user, err := services.UserService.Bind(githubId, bindType, username, email, password, rePassword, nickname)
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
	session.SetCurrentUser(this.Ctx, user.Id)

	redirectUrl := simple.FormValue(this.Ctx, "redirectUrl")
	if len(redirectUrl) > 0 {
		this.Ctx.Redirect(redirectUrl, iris.StatusSeeOther)
	} else {
		this.Ctx.Redirect("/", iris.StatusSeeOther)
	}
}

// 未读消息数量
func (this *UserController) GetMsgcount() *simple.JsonResult {
	user := session.GetCurrentUser(this.Ctx)
	var count int64 = 0
	if user != nil {
		count = services.MessageService.GetUnReadCount(user.Id)
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

	articles, paging := services.ArticleService.Query(simple.NewParamQueries(ctx).
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
	user := session.GetCurrentUser(ctx)
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

	list, paging := services.UserArticleTagService.Query(simple.NewParamQueries(ctx).Eq("user_id", userId).
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
	user := session.GetCurrentUser(ctx)
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

	favorites, paging := services.FavoriteService.Query(simple.NewParamQueries(ctx).
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
	user := session.GetCurrentUser(ctx)
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
	messages, paging := services.MessageService.Query(simple.NewParamQueries(ctx).
		Eq("user_id", user.Id).Page(page, 20).Desc("id"))

	// 全部标记为已读
	services.MessageService.MarkReadAll(user.Id)

	render.View(ctx, "user/messages.html", iris.Map{
		"User":                     render.BuildUser(user),
		"Messages":                 render.BuildMessages(messages),
		"Page":                     paging,
		"PrePageUrl":               utils.BuildMessagesUrl(page - 1),
		"NextPageUrl":              utils.BuildMessagesUrl(page + 1),
		utils.GlobalFieldSiteTitle: user.Nickname + " - 消息",
	})
}

// 用户中心 - 主题列表
func GetUserTopics(ctx context.Context) {
	userId := ctx.Params().GetInt64Default("userId", 0)
	page := ctx.Params().GetIntDefault("page", 1)

	user := cache.UserCache.Get(userId)
	if user == nil {
		ctx.StatusCode(404)
		return
	}

	// 查询列表
	topics, paging := services.TopicService.Query(simple.NewParamQueries(ctx).
		Eq("user_id", userId).Page(page, 20).Desc("id"))

	render.View(ctx, "user/topics.html", iris.Map{
		"User":                     render.BuildUser(user),
		"Topics":                   render.BuildTopics(topics),
		"Page":                     paging,
		"PrePageUrl":               utils.BuildUserTopicsUrl(userId, page-1),
		"NextPageUrl":              utils.BuildUserTopicsUrl(userId, page+1),
		utils.GlobalFieldSiteTitle: user.Nickname + " - 话题",
	})
}
