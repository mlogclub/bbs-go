package server

import (
	"bbs-go/internal/controllers/admin"
	"bbs-go/internal/controllers/api"
	"bbs-go/internal/middleware"
	"bbs-go/internal/pkg/config"
	"log/slog"
	"os"
	"strings"

	"github.com/iris-contrib/middleware/cors"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/mlogclub/simple/web"

	"github.com/kataras/iris/v12/middleware/logger"
	"github.com/kataras/iris/v12/middleware/recover"
)

func NewServer() {
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

	app.OnAnyErrorCode(func(ctx iris.Context) {
		path := ctx.Path()
		var err error
		if strings.Contains(path, "/api/admin/") {
			err = ctx.JSON(web.JsonErrorCode(ctx.GetStatusCode(), "Http error"))
		}
		if err != nil {
			slog.Error(err.Error(), slog.Any("err", err))
		}
	})

	app.Any("/", func(i iris.Context) {
		_ = i.JSON(map[string]interface{}{
			"engine": "bbs-go",
		})
	})

	conf := config.Instance

	// api
	mvc.Configure(app.Party("/api"), func(m *mvc.Application) {
		m.Party("/topic").Handle(new(api.TopicController))
		m.Party("/article").Handle(new(api.ArticleController))
		m.Party("/login").Handle(new(api.LoginController))
		m.Party("/user").Handle(new(api.UserController))
		m.Party("/tag").Handle(new(api.TagController))
		m.Party("/comment").Handle(new(api.CommentController))
		m.Party("/favorite").Handle(new(api.FavoriteController))
		m.Party("/like").Handle(new(api.LikeController))
		m.Party("/checkin").Handle(new(api.CheckinController))
		m.Party("/config").Handle(new(api.ConfigController))
		m.Party("/upload").Handle(new(api.UploadController))
		m.Party("/link").Handle(new(api.LinkController))
		m.Party("/captcha").Handle(new(api.CaptchaController))
		m.Party("/search").Handle(new(api.SearchController))
		m.Party("/fans").Handle(new(api.FansController))
		m.Party("/user-report").Handle(new(api.UserReportController))
	})

	// admin
	mvc.Configure(app.Party("/api/admin"), func(m *mvc.Application) {
		m.Router.Use(middleware.AdminAuth)
		m.Party("/common").Handle(new(admin.CommonController))
		m.Party("/user").Handle(new(admin.UserController))
		m.Party("/tag").Handle(new(admin.TagController))
		m.Party("/article").Handle(new(admin.ArticleController))
		m.Party("/comment").Handle(new(admin.CommentController))
		m.Party("/favorite").Handle(new(admin.FavoriteController))
		m.Party("/article-tag").Handle(new(admin.ArticleTagController))
		m.Party("/topic").Handle(new(admin.TopicController))
		m.Party("/topic-node").Handle(new(admin.TopicNodeController))
		m.Party("/sys-config").Handle(new(admin.SysConfigController))
		m.Party("/link").Handle(new(admin.LinkController))
		m.Party("/user-score-log").Handle(new(admin.UserScoreLogController))
		m.Party("/operate-log").Handle(new(admin.OperateLogController))
		m.Party("/user-report").Handle(new(admin.UserReportController))
		m.Party("/forbidden-word").Handle(new(admin.ForbiddenWordController))

		m.Party("/role").Handle(new(admin.RoleController))
		m.Party("/menu").Handle(new(admin.MenuController))
	})

	if err := app.Listen(":"+conf.Port,
		iris.WithConfiguration(iris.Configuration{
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
		}),
	); err != nil {
		slog.Error(err.Error(), slog.Any("err", err))
		os.Exit(-1)
	}
}
