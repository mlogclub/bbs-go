package app

import (
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

	"github.com/mlogclub/mlog/controllers"
	"github.com/mlogclub/mlog/controllers/admin"
	"github.com/mlogclub/mlog/controllers/render"
	"github.com/mlogclub/mlog/middleware"
	"github.com/mlogclub/mlog/services"
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
			m.Register(services.Instances...)

			m.Router.Use(middleware.NewGlobalMiddleware())

			m.Party("/upload").Handle(new(controllers.UploadController))

			m.Party("/").Handle(new(controllers.IndexController))

			m.Party("/article").Handle(new(controllers.ArticleController))
			m.Router.Get("/articles", controllers.GetArticles)
			m.Router.Get("/articles/{page:int}", controllers.GetArticles)
			m.Router.Get("/articles/tag/{tagId:int64}", controllers.GetTagArticles)
			m.Router.Get("/articles/tag/{tagId:int64}/{page:int}", controllers.GetTagArticles)
			m.Router.Get("/articles/cat/{categoryId:int64}", controllers.GetCategoryArticles)
			m.Router.Get("/articles/cat/{categoryId:int64}/{page:int}", controllers.GetCategoryArticles)

			m.Party("/topic").Handle(new(controllers.TopicController))
			m.Router.Get("/topics", controllers.GetTopics)
			m.Router.Get("/topics/{page:int}", controllers.GetTopics)

			m.Party("/comment").Handle(new(controllers.CommentController))

			m.Party("/user").Handle(new(controllers.UserController))
			m.Router.Get("/user/{userId:int64}/articles", controllers.GetUserArticles)
			m.Router.Get("/user/{userId:int64}/articles/{page:int}", controllers.GetUserArticles)
			m.Router.Get("/user/{userId:int64}/topics", controllers.GetUserTopics)
			m.Router.Get("/user/{userId:int64}/topics/{page:int}", controllers.GetUserTopics)
			m.Router.Get("/user/{userId:int64}/tags", controllers.GetUserTags)
			m.Router.Get("/user/{userId:int64}/tags/{page:int}", controllers.GetUserTags)
			m.Router.Get("/user/{userId:int64}/messages", controllers.GetUserMessages)
			m.Router.Get("/user/{userId:int64}/messages/{page:int}", controllers.GetUserMessages)
			m.Router.Get("/user/{userId:int64}/favorites", controllers.GetUserFavorites)
			m.Router.Get("/user/{userId:int64}/favorites/{page:int}", controllers.GetUserFavorites)

			// 标签
			m.Party("/tag").Handle(new(controllers.TagController))
			m.Router.Get("/tags", controllers.GetTags)
			m.Router.Get("/tags/{page:int}", controllers.GetTags)

			m.Party("/share").Handle(new(controllers.ArticleShareController))

			m.Router.Get("/redirect", func(ctx context.Context) {
				url := ctx.FormValue("url")
				render.View(ctx, "redirect.html", iris.Map{
					utils.GlobalFieldSiteTitle: "跳转中",
					"Url":                      url,
				})
			})
		})

		mvc.Configure(app.Party("/api/admin"), func(m *mvc.Application) {
			m.Register(services.Instances...)

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
		})

		mvc.Configure(app.Party("/oauth"), func(m *mvc.Application) {
			m.Register(services.Instances...)
			m.Party("/").Handle(new(controllers.OauthServerController))
		})

		mvc.Configure(app.Party("/oauth/client"), func(m *mvc.Application) {
			m.Party("/").Handle(new(controllers.OauthClientController))
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
	app.StaticWeb("/static", config.Conf.StaticPath)
	app.StaticWeb("/", config.Conf.StaticPath)

	engine := iris.HTML(config.Conf.ViewsPath, ".html").Reload(true)
	viewFunctions(engine)

	app.RegisterView(engine)
}

func viewFunctions(engine *view.HTMLEngine) {
	engine.AddFunc("siteTitle", func(siteTitle string) string {
		if len(siteTitle) > 0 {
			return siteTitle + " | " + config.Conf.SiteTitle
		}
		return config.Conf.SiteTitle
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
