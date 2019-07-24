package app

import (
	"github.com/mlogclub/mlog/controllers/web"
	"github.com/mlogclub/mlog/model"
	"github.com/mlogclub/mlog/services/cache"
	"github.com/mlogclub/mlog/utils/config"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/iris-contrib/middleware/cors"
	"github.com/kataras/iris"
	"github.com/kataras/iris/context"
	"github.com/kataras/iris/middleware/logger"
	"github.com/kataras/iris/middleware/recover"
	"github.com/kataras/iris/mvc"
	"github.com/kataras/iris/view"
	"github.com/mlogclub/simple"
	"github.com/sirupsen/logrus"

	"github.com/mlogclub/mlog/controllers/admin"
	"github.com/mlogclub/mlog/controllers/render"
	"github.com/mlogclub/mlog/middleware"
	"github.com/mlogclub/mlog/utils"
)

func InitIris() {
	app := iris.New()
	app.Logger().SetLevel("warn")
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // allows everything, use that to change the hosts.
		AllowCredentials: true,
		MaxAge:           600,
		AllowedMethods:   []string{iris.MethodGet, iris.MethodPost, iris.MethodOptions, iris.MethodHead, iris.MethodDelete, iris.MethodPut},
		AllowedHeaders:   []string{"*"},
	}))
	app.AllowMethods(iris.MethodOptions)

	app.OnAnyErrorCode(func(ctx context.Context) {
		path := ctx.Path()
		var err error
		if strings.Contains(path, "/api/admin/") {
			_, err = ctx.JSON(simple.ErrorCode(ctx.GetStatusCode(), "Http error"))
		} else {
			if ctx.GetStatusCode() == 404 {
				render.View(ctx, "404.html", iris.Map{
					"StatusCode": ctx.GetStatusCode(),
				})
			} else {
				render.View(ctx, "5xx.html", iris.Map{
					"StatusCode": ctx.GetStatusCode(),
				})
			}
		}
		if err != nil {
			logrus.Error(err)
		}
	})

	handleViews(app)

	{
		mvc.Configure(app.Party("/"), func(m *mvc.Application) {
			m.Router.Use(middleware.NewGlobalMiddleware())

			m.Party("/upload").Handle(new(web.UploadController))

			m.Party("/").Handle(new(web.IndexController))

			m.Party("/article").Handle(new(web.ArticleController))
			m.Router.Get("/articles", web.GetArticles)
			m.Router.Get("/articles/{page:int}", web.GetArticles)
			m.Router.Get("/articles/tag/{tagId:int64}", web.GetTagArticles)
			m.Router.Get("/articles/tag/{tagId:int64}/{page:int}", web.GetTagArticles)
			m.Router.Get("/articles/cat/{categoryId:int64}", web.GetCategoryArticles)
			m.Router.Get("/articles/cat/{categoryId:int64}/{page:int}", web.GetCategoryArticles)

			m.Party("/topic").Handle(new(web.TopicController))
			m.Router.Get("/topics", web.GetTopics)
			m.Router.Get("/topics/{page:int}", web.GetTopics)
			m.Router.Get("/topics/tag/{tagId:int64}", web.GetTagTopics)
			m.Router.Get("/topics/tag/{tagId:int64}/{page:int}", web.GetTagTopics)

			m.Party("/comment").Handle(new(web.CommentController))

			m.Party("/user").Handle(new(web.UserController))
			m.Router.Get("/user/{userId:int64}/articles", web.GetUserArticles)
			m.Router.Get("/user/{userId:int64}/articles/{page:int}", web.GetUserArticles)
			m.Router.Get("/user/{userId:int64}/topics", web.GetUserTopics)
			m.Router.Get("/user/{userId:int64}/topics/{page:int}", web.GetUserTopics)
			m.Router.Get("/user/{userId:int64}/messages", web.GetUserMessages)
			m.Router.Get("/user/{userId:int64}/messages/{page:int}", web.GetUserMessages)
			m.Router.Get("/user/{userId:int64}/favorites", web.GetUserFavorites)
			m.Router.Get("/user/{userId:int64}/favorites/{page:int}", web.GetUserFavorites)

			// 标签
			m.Party("/tag").Handle(new(web.TagController))
			m.Router.Get("/tags", web.GetTags)
			m.Router.Get("/tags/{page:int}", web.GetTags)

			// 配置
			m.Party("/config").Handle(new(web.ConfigController))

			m.Router.Get("/redirect", func(ctx context.Context) {
				url := ctx.FormValue("url")
				render.View(ctx, "redirect.html", iris.Map{
					model.TplSiteTitle: "跳转中",
					"Url":              url,
				})
			})
		})

		mvc.Configure(app.Party("/api/admin"), func(m *mvc.Application) {
			m.Router.Use(middleware.AdminAuthHandler)
			m.Party("/user").Handle(new(admin.UserController))
			m.Party("/github-user").Handle(new(admin.GithubUserController))
			m.Party("/category").Handle(new(admin.CategoryController))
			m.Party("/tag").Handle(new(admin.TagController))
			m.Party("/article").Handle(new(admin.ArticleController))
			m.Party("/comment").Handle(new(admin.CommentController))
			m.Party("/favorite").Handle(new(admin.FavoriteController))
			m.Party("/article-tag").Handle(new(admin.ArticleTagController))
			m.Party("/topic").Handle(new(admin.TopicController))
			m.Party("/oauth-client").Handle(new(admin.OauthClientController))
			m.Party("/oauth-token").Handle(new(admin.OauthTokenController))
			m.Party("/sys-config").Handle(new(admin.SysConfigController))
		})

		mvc.Configure(app.Party("/oauth"), func(m *mvc.Application) {
			m.Party("/").Handle(new(web.OauthServerController))
		})

		mvc.Configure(app.Party("/oauth/client"), func(m *mvc.Application) {
			m.Party("/").Handle(new(web.OauthClientController))
		})
	}

	server := &http.Server{Addr: ":" + config.Conf.Port}
	handleSignal(server)
	err := app.Run(iris.Server(server), iris.WithConfiguration(iris.Configuration{
		DisableStartupLog:                 false,
		DisableInterruptHandler:           false,
		DisablePathCorrection:             false,
		EnablePathEscape:                  false,
		FireMethodNotAllowed:              false,
		DisableBodyConsumptionOnUnmarshal: false,
		DisableAutoFireStatusCode:         false,
		EnableOptimizations:               true,
		TimeFormat:                        "2006-01-02 15:04:05",
		Charset:                           "UTF-8",
	}))
	if err != nil {
		logrus.Error(err)
		os.Exit(-1)
	}
}

func handleViews(app *iris.Application) {
	if len(config.Conf.RootStaticPath) > 0 {
		app.StaticWeb("/", config.Conf.RootStaticPath)
	}
	if len(config.Conf.StaticPath) > 0 {
		app.StaticWeb("/static", config.Conf.StaticPath)
	}

	engine := iris.HTML(config.Conf.ViewsPath, ".html").Reload(true)
	viewFunctions(engine)
	app.RegisterView(engine)
}

func viewFunctions(engine *view.HTMLEngine) {
	engine.AddFunc("siteTitle", func(title string) string {
		siteTitle := cache.SysConfigCache.GetValue(model.SysConfigSiteTitle)
		if len(title) > 0 {
			return title + " | " + siteTitle
		}
		return siteTitle
	})
	engine.AddFunc("formatDate", func(timestamp int64) string {
		return simple.TimeFormat(simple.TimeFromTimestamp(timestamp), simple.FMT_DATE_TIME)
	})
	engine.AddFunc("baseUrl", func() string {
		return config.Conf.BaseUrl
	})
	engine.AddFunc("absUrl", func(path string) string {
		return utils.BuildAbsUrl(path)
	})
	engine.AddFunc("articleUrl", func(articleId int64) string {
		return utils.BuildArticleUrl(articleId)
	})
	engine.AddFunc("articlesUrl", func(page int) string {
		return utils.BuildArticlesUrl(page)
	})
	engine.AddFunc("tagArticlesUrl", func(tagId int64, page int) string {
		return utils.BuildTagArticlesUrl(tagId, page)
	})
	engine.AddFunc("categoryArticlesUrl", func(categoryId int64, page int) string {
		return utils.BuildCategoryArticlesUrl(categoryId, page)
	})
	engine.AddFunc("topicUrl", func(topicId int64) string {
		return utils.BuildTopicUrl(topicId)
	})
	engine.AddFunc("topicsUrl", func(page int) string {
		return utils.BuildTopicsUrl(page)
	})
	engine.AddFunc("tagTopicsUrl", func(tagId int64, page int) string {
		return utils.BuildTagTopicsUrl(tagId, page)
	})
	engine.AddFunc("userUrl", func(userId int64) string {
		return utils.BuildUserUrl(userId)
	})
	engine.AddFunc("prettyTime", func(timestamp int64) string {
		return simple.PrettyTime(timestamp)
	})
}

func handleSignal(server *http.Server) {
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)

	go func() {
		s := <-c
		logrus.Infof("got signal [%s], exiting now", s)
		if err := server.Close(); nil != err {
			logrus.Errorf("server close failed: " + err.Error())
		}

		simple.CloseDB()

		logrus.Infof("Exited")
		os.Exit(0)
	}()
}
