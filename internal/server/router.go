package server

import (
	adminHandlers "bbs-go/internal/handlers/admin"
	apiHandlers "bbs-go/internal/handlers/api"
	"bbs-go/internal/middleware"
	"bbs-go/internal/pkg/config"
	"bbs-go/internal/pkg/ginx"
	"bbs-go/internal/pkg/respath"
	"bbs-go/internal/services"
	webspa "bbs-go/web"
	"log/slog"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/mlogclub/simple/common/strs"
	"github.com/mlogclub/simple/web"
	"github.com/spf13/cast"
)

func NewServer() {
	printBanner()
	if err := newRouter().Run(":" + cast.ToString(config.Instance.Port)); err != nil {
		slog.Error(err.Error(), slog.Any("err", err))
		os.Exit(-1)
	}
}

func newRouter() *gin.Engine {
	conf := config.Instance
	if conf == nil {
		conf = &config.Config{}
	}

	gin.SetMode(gin.ReleaseMode)
	app := gin.New()
	app.Use(gin.Recovery())
	app.Use(gin.Logger())
	corsConfig := cors.Config{
		AllowOrigins:     conf.AllowedOrigins,
		AllowCredentials: true,
		MaxAge:           600,
		AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodOptions, http.MethodHead, http.MethodDelete, http.MethodPut},
		AllowHeaders:     []string{"*"},
	}
	if len(corsConfig.AllowOrigins) == 0 {
		corsConfig.AllowAllOrigins = true
		corsConfig.AllowCredentials = false
	}
	app.Use(cors.New(corsConfig))
	app.Use(middleware.AttachmentMiddleware)

	registerAPIRoutes(app.Group("/api", middleware.InstallMiddleware, middleware.AuthMiddleware))
	registerAdminRoutes(app.Group("/api/admin", middleware.InstallMiddleware, middleware.AuthMiddleware, middleware.AdminMiddleware))

	app.StaticFS("/res", ginx.StaticFiles(respath.ResDir()))
	ginx.HandleSPA(app, ginx.SPAOptions{
		Root:         "./web/build/spa",
		EmbeddedFS:   webspa.SPA,
		EmbeddedRoot: "build/spa",
		DirOptions: ginx.DirOptions{
			ShowList:  false,
			SPA:       true,
			IndexName: "index.html",
		},
		NotFoundPrefixes: []string{"/api/", "/res/"},
		NotFoundHandler: func(ctx *gin.Context) {
			ginx.WriteHttpStatusJSON(ctx, http.StatusNotFound, web.JsonErrorCode(http.StatusNotFound, "Not found"))
		},
	})
	app.GET("/sitemap.xml", func(ctx *gin.Context) {
		redirectURL := services.SeoSitemapService.RedirectURL()
		if strs.IsBlank(redirectURL) {
			ginx.WriteHttpStatusJSON(ctx, http.StatusNotFound, web.JsonErrorCode(http.StatusNotFound, "Not found"))
			return
		}
		ctx.Redirect(http.StatusFound, redirectURL)
	})

	return app
}

func registerAPIRoutes(group *gin.RouterGroup) {
	installGroup := group.Group("/install")
	installGroup.GET("/status", apiHandlers.InstallStatus)
	installGroup.POST("/test_db_connection", apiHandlers.InstallTestDbConnection)
	installGroup.POST("/install", apiHandlers.InstallInstall)

	topicGroup := group.Group("/topic")
	topicGroup.GET("/category_navs", apiHandlers.CategoryNavs)
	topicGroup.GET("/categories", apiHandlers.Categories)
	topicGroup.GET("/category", apiHandlers.Category)
	topicGroup.POST("/create", apiHandlers.TopicCreate)
	topicGroup.GET("/edit/:id", apiHandlers.TopicEditForm)
	topicGroup.POST("/edit/:id", apiHandlers.TopicEdit)
	topicGroup.POST("/delete/:id", apiHandlers.TopicRemove)
	topicGroup.POST("/recommend/:id", apiHandlers.TopicRecommend)
	topicGroup.GET("/recentlikes/:id", apiHandlers.TopicRecentlikes)
	topicGroup.GET("/recent", apiHandlers.TopicRecent)
	topicGroup.GET("/user_topics", apiHandlers.TopicUserTopics)
	topicGroup.GET("/topics", apiHandlers.TopicTopics)
	topicGroup.POST("/accept_answer/:id", apiHandlers.TopicAcceptAnswer)
	topicGroup.POST("/unaccept_answer/:id", apiHandlers.TopicUnacceptAnswer)
	topicGroup.GET("/tag/topics", apiHandlers.TopicTagTopics)
	topicGroup.GET("/favorite/:id", apiHandlers.TopicFavorite)
	topicGroup.POST("/sticky/:id", apiHandlers.TopicSticky)
	topicGroup.GET("/hide_content", apiHandlers.TopicHideContent)
	topicGroup.GET("/:id", apiHandlers.TopicDetail)

	articleGroup := group.Group("/article")
	articleGroup.POST("/create", apiHandlers.ArticleCreate)
	articleGroup.GET("/edit/:id", apiHandlers.ArticleEditForm)
	articleGroup.POST("/edit/:id", apiHandlers.ArticleEdit)
	articleGroup.POST("/delete/:id", apiHandlers.ArticleRemove)
	articleGroup.POST("/favorite/:id", apiHandlers.ArticleFavorite)
	articleGroup.GET("/redirect/:id", apiHandlers.ArticleRedirect)
	articleGroup.GET("/user_articles", apiHandlers.ArticleUserArticles)
	articleGroup.GET("/articles", apiHandlers.ArticleArticles)
	articleGroup.GET("/tag/articles", apiHandlers.ArticleTagArticles)
	articleGroup.GET("/:id", apiHandlers.ArticleDetail)

	loginGroup := group.Group("/login")
	loginGroup.POST("/signup", apiHandlers.LoginSignup)
	loginGroup.POST("/signin", apiHandlers.LoginSignin)
	loginGroup.POST("/send_reset_password_email", apiHandlers.LoginSendResetPasswordEmail)
	loginGroup.POST("/reset_password", apiHandlers.LoginResetPassword)
	loginGroup.GET("/signout", apiHandlers.LoginSignout)
	loginGroup.POST("/login_sms_code", apiHandlers.LoginLoginSmsCode)
	loginGroup.POST("/login_sms", apiHandlers.LoginLoginSms)
	loginGroup.GET("/wx_login_config", apiHandlers.LoginWxLoginConfig)
	loginGroup.POST("/wx_login_submit", apiHandlers.LoginWxLoginSubmit)
	loginGroup.POST("/wx_bind", apiHandlers.LoginWxBind)
	loginGroup.POST("/wx_unbind", apiHandlers.LoginWxUnbind)
	loginGroup.GET("/google_login_config", apiHandlers.LoginGoogleLoginConfig)
	loginGroup.POST("/google_login_submit", apiHandlers.LoginGoogleLoginSubmit)
	loginGroup.POST("/google_bind", apiHandlers.LoginGoogleBind)
	loginGroup.POST("/google_one_tap", apiHandlers.LoginGoogleOneTap)
	loginGroup.POST("/google_unbind", apiHandlers.LoginGoogleUnbind)
	loginGroup.GET("/github_login_config", apiHandlers.LoginGithubLoginConfig)
	loginGroup.POST("/github_login_submit", apiHandlers.LoginGithubLoginSubmit)
	loginGroup.POST("/github_unbind", apiHandlers.LoginGithubUnbind)

	userGroup := group.Group("/user")
	userGroup.GET("/current", apiHandlers.UserCurrent)
	userGroup.POST("/update/:id", apiHandlers.UserUpdate)
	userGroup.POST("/update_avatar", apiHandlers.UserUpdateAvatar)
	userGroup.POST("/set_username", apiHandlers.UserSetUsername)
	userGroup.POST("/set_email", apiHandlers.UserSetEmail)
	userGroup.POST("/set_password", apiHandlers.UserSetPassword)
	userGroup.POST("/update_password", apiHandlers.UserUpdatePassword)
	userGroup.POST("/set_background_image", apiHandlers.UserSetBackgroundImage)
	userGroup.GET("/favorites", apiHandlers.UserFavorites)
	userGroup.GET("/msg_recent", apiHandlers.UserMsgRecent)
	userGroup.GET("/messages", apiHandlers.UserMessages)
	userGroup.GET("/score_logs", apiHandlers.UserScoreLogs)
	userGroup.GET("/score/rank", apiHandlers.UserScoreRank)
	userGroup.POST("/forbidden", apiHandlers.UserForbidden)
	userGroup.POST("/send_verify_email", apiHandlers.UserSendVerifyEmail)
	userGroup.POST("/verify_email", apiHandlers.UserVerifyEmail)
	userGroup.GET("/wx_bind_info", apiHandlers.UserWxBindInfo)
	userGroup.GET("/google_bind_info", apiHandlers.UserGoogleBindInfo)
	userGroup.GET("/github_bind_info", apiHandlers.UserGithubBindInfo)
	userGroup.GET("/:id", apiHandlers.UserDetail)

	tagGroup := group.Group("/tag")
	tagGroup.GET("/tags", apiHandlers.TagTags)
	tagGroup.POST("/autocomplete", apiHandlers.TagAutocompleteSubmit)
	tagGroup.GET("/:id", apiHandlers.TagDetail)

	commentGroup := group.Group("/comment")
	commentGroup.GET("/comments", apiHandlers.CommentComments)
	commentGroup.GET("/replies", apiHandlers.CommentReplies)
	commentGroup.POST("/create", apiHandlers.CommentCreate)
	commentGroup.POST("/delete/:id", apiHandlers.CommentRemove)

	favoriteGroup := group.Group("/favorite")
	favoriteGroup.POST("/add", apiHandlers.FavoriteAdd)
	favoriteGroup.POST("/delete", apiHandlers.FavoriteRemove)

	likeGroup := group.Group("/like")
	likeGroup.POST("/like", apiHandlers.LikeLike)
	likeGroup.POST("/unlike", apiHandlers.LikeUnlike)
	likeGroup.GET("/liked_ids", apiHandlers.LikeLikedIds)
	likeGroup.GET("/liked", apiHandlers.LikeLiked)

	checkinGroup := group.Group("/checkin")
	checkinGroup.POST("/checkin", apiHandlers.CheckinSubmit)
	checkinGroup.GET("/checkin", apiHandlers.CheckinStatus)
	checkinGroup.GET("/rank", apiHandlers.CheckinRank)

	configGroup := group.Group("/config")
	configGroup.GET("/configs", apiHandlers.ConfigConfigs)
	configGroup.GET("/about", apiHandlers.ConfigAbout)

	uploadGroup := group.Group("/upload")
	uploadGroup.POST("", apiHandlers.UploadHandle)

	attachmentGroup := group.Group("/attachment")
	attachmentGroup.POST("/upload", apiHandlers.AttachmentUpload)
	attachmentGroup.GET("/download/:id", apiHandlers.AttachmentDownload)
	attachmentGroup.POST("/update_download_score", apiHandlers.AttachmentUpdateDownloadScore)

	linkGroup := group.Group("/link")
	linkGroup.GET("/list", apiHandlers.LinkList)
	linkGroup.GET("/top_links", apiHandlers.LinkTopLinks)

	captchaGroup := group.Group("/captcha")
	captchaGroup.GET("/request", apiHandlers.CaptchaRequest)
	captchaGroup.GET("/verify", apiHandlers.CaptchaVerify)
	captchaGroup.GET("/request_angle", apiHandlers.CaptchaRequestAngle)

	searchGroup := group.Group("/search")
	searchGroup.GET("/topic", apiHandlers.SearchTopic)
	searchGroup.GET("/article", apiHandlers.SearchArticle)
	searchGroup.GET("/user", apiHandlers.SearchUser)

	fansGroup := group.Group("/fans")
	fansGroup.POST("/follow", apiHandlers.FansFollow)
	fansGroup.POST("/unfollow", apiHandlers.FansUnfollow)
	fansGroup.GET("/is_followed", apiHandlers.FansIsFollowed)
	fansGroup.GET("/fans", apiHandlers.FansFans)
	fansGroup.GET("/followed", apiHandlers.FansFollowed)
	fansGroup.GET("/recent/fans", apiHandlers.FansRecentFans)
	fansGroup.GET("/recent/follow", apiHandlers.FansRecentFollow)

	userReportGroup := group.Group("/user-report")
	userReportGroup.POST("/submit", apiHandlers.UserReportSubmit)

	taskGroup := group.Group("/task")
	taskGroup.GET("/tasks", apiHandlers.TaskTasks)
	taskGroup.GET("/groups", apiHandlers.TaskGroups)

	badgeGroup := group.Group("/badge")
	badgeGroup.GET("/badges", apiHandlers.BadgeBadges)

	voteGroup := group.Group("/vote")
	voteGroup.POST("/cast", apiHandlers.VoteCast)
	voteGroup.GET("/:id", apiHandlers.VoteDetail)

}

func registerAdminRoutes(group *gin.RouterGroup) {
	roleGroup := group.Group("/role")
	roleGroup.GET("/roles", adminHandlers.Roles)
	roleGroup.POST("/list", adminHandlers.RoleList)
	roleGroup.GET("/permissions", adminHandlers.RolePermissions)
	roleGroup.POST("/create", adminHandlers.RoleCreate)
	roleGroup.POST("/update", adminHandlers.RoleUpdate)
	roleGroup.POST("/update_permissions", adminHandlers.RoleUpdatePermissions)
	roleGroup.POST("/delete", adminHandlers.RoleRemove)
	roleGroup.POST("/update_sort", adminHandlers.RoleUpdateSort)
	roleGroup.GET("/:id", adminHandlers.RoleDetail)

	dictTypeGroup := group.Group("/dict-type")
	dictTypeGroup.GET("/list", adminHandlers.DictTypeList)
	dictTypeGroup.POST("/create", adminHandlers.DictTypeCreate)
	dictTypeGroup.POST("/update", adminHandlers.DictTypeUpdate)
	dictTypeGroup.POST("/delete", adminHandlers.DictTypeRemove)
	dictTypeGroup.GET("/:id", adminHandlers.DictTypeDetail)

	dictGroup := group.Group("/dict")
	dictGroup.GET("/list", adminHandlers.DictList)
	dictGroup.POST("/create", adminHandlers.DictCreate)
	dictGroup.POST("/update", adminHandlers.DictUpdate)
	dictGroup.POST("/delete", adminHandlers.DictRemove)
	dictGroup.POST("/update_sort", adminHandlers.DictUpdateSort)
	dictGroup.GET("/dicts", adminHandlers.DictDicts)
	dictGroup.GET("/:id", adminHandlers.DictDetail)

	emailLogGroup := group.Group("/email-log")
	emailLogGroup.POST("/list", adminHandlers.EmailLogList)
	emailLogGroup.GET("/:id", adminHandlers.EmailLogDetail)

	commonGroup := group.Group("/common")
	commonGroup.GET("/overview", adminHandlers.CommonOverview)
	commonGroup.GET("/task_event_types", adminHandlers.CommonTaskEventTypes)

	userGroup := group.Group("/user")
	userGroup.GET("/synccount", adminHandlers.UserSynccount)
	userGroup.POST("/list", adminHandlers.UserList)
	userGroup.POST("/create", adminHandlers.UserCreate)
	userGroup.POST("/update", adminHandlers.UserUpdate)
	userGroup.POST("/forbidden", adminHandlers.UserForbidden)
	userGroup.POST("/update_password", adminHandlers.UserUpdatePassword)
	userGroup.POST("/reset_password", adminHandlers.UserResetPassword)
	userGroup.GET("/:id", adminHandlers.UserDetail)

	tagGroup := group.Group("/tag")
	tagGroup.POST("/list", adminHandlers.TagList)
	tagGroup.POST("/create", adminHandlers.TagCreate)
	tagGroup.POST("/update", adminHandlers.TagUpdate)
	tagGroup.GET("/autocomplete", adminHandlers.TagAutocomplete)
	tagGroup.GET("/tags", adminHandlers.TagTags)
	tagGroup.GET("/:id", adminHandlers.TagDetail)

	articleGroup := group.Group("/article")
	articleGroup.POST("/list", adminHandlers.ArticleList)
	articleGroup.POST("/update", adminHandlers.ArticleUpdate)
	articleGroup.GET("/tags", adminHandlers.ArticleTags)
	articleGroup.POST("/tags", adminHandlers.ArticleSaveTags)
	articleGroup.POST("/delete", adminHandlers.ArticleRemove)
	articleGroup.POST("/audit", adminHandlers.ArticleAudit)
	articleGroup.GET("/:id", adminHandlers.ArticleDetail)

	favoriteGroup := group.Group("/favorite")
	favoriteGroup.POST("/list", adminHandlers.FavoriteList)
	favoriteGroup.POST("/create", adminHandlers.FavoriteCreate)
	favoriteGroup.POST("/update", adminHandlers.FavoriteUpdate)
	favoriteGroup.GET("/:id", adminHandlers.FavoriteDetail)

	articleTagGroup := group.Group("/article-tag")
	articleTagGroup.POST("/list", adminHandlers.ArticleTagList)
	articleTagGroup.POST("/create", adminHandlers.ArticleTagCreate)
	articleTagGroup.POST("/update", adminHandlers.ArticleTagUpdate)
	articleTagGroup.GET("/:id", adminHandlers.ArticleTagDetail)

	topicGroup := group.Group("/topic")
	topicGroup.POST("/list", adminHandlers.TopicList)
	topicGroup.POST("/recommend", adminHandlers.TopicRecommend)
	topicGroup.DELETE("/recommend", adminHandlers.TopicRemoveRecommend)
	topicGroup.POST("/delete", adminHandlers.TopicRemove)
	topicGroup.POST("/undelete", adminHandlers.TopicUndelete)
	topicGroup.POST("/audit", adminHandlers.TopicAudit)
	topicGroup.POST("/accept_answer", adminHandlers.TopicAcceptAnswer)
	topicGroup.POST("/unaccept_answer", adminHandlers.TopicUnacceptAnswer)
	topicGroup.POST("/mark_solved", adminHandlers.TopicMarkSolved)
	topicGroup.POST("/mark_unsolved", adminHandlers.TopicMarkUnsolved)
	topicGroup.GET("/:id", adminHandlers.TopicDetail)

	categoryGroup := group.Group("/category")
	categoryGroup.POST("/list", adminHandlers.CategoryList)
	categoryGroup.POST("/create", adminHandlers.CategoryCreate)
	categoryGroup.POST("/update", adminHandlers.CategoryUpdate)
	categoryGroup.GET("/options", adminHandlers.CategoryOptions)
	categoryGroup.POST("/update_sort", adminHandlers.CategoryUpdateSort)
	categoryGroup.POST("/delete", adminHandlers.CategoryRemove)
	categoryGroup.GET("/:id", adminHandlers.CategoryDetail)

	sysConfigGroup := group.Group("/sys-config")
	sysConfigGroup.POST("/list", adminHandlers.SysConfigList)
	sysConfigGroup.GET("/configs", adminHandlers.SysConfigConfigs)
	sysConfigGroup.POST("/save", adminHandlers.SysConfigSave)
	sysConfigGroup.GET("/:id", adminHandlers.SysConfigDetail)

	searchGroup := group.Group("/search")
	searchGroup.GET("/reindex/status", adminHandlers.SearchReindexStatus)
	searchGroup.POST("/reindex", adminHandlers.SearchReindex)

	seoGroup := group.Group("/seo")
	seoGroup.GET("/sitemap/status", adminHandlers.SeoSitemapStatus)
	seoGroup.POST("/sitemap/generate", adminHandlers.SeoSitemapGenerate)

	linkGroup := group.Group("/link")
	linkGroup.POST("/list", adminHandlers.LinkList)
	linkGroup.POST("/create", adminHandlers.LinkCreate)
	linkGroup.POST("/update", adminHandlers.LinkUpdate)
	linkGroup.POST("/delete", adminHandlers.LinkRemove)
	linkGroup.POST("/update_sort", adminHandlers.LinkUpdateSort)
	linkGroup.GET("/:id", adminHandlers.LinkDetail)

	userScoreLogGroup := group.Group("/user-score-log")
	userScoreLogGroup.POST("/list", adminHandlers.UserScoreLogList)
	userScoreLogGroup.GET("/:id", adminHandlers.UserScoreLogDetail)

	taskConfigGroup := group.Group("/task-config")
	taskConfigGroup.GET("/groups", adminHandlers.TaskConfigGroups)
	taskConfigGroup.POST("/list", adminHandlers.TaskConfigList)
	taskConfigGroup.POST("/create", adminHandlers.TaskConfigCreate)
	taskConfigGroup.POST("/update", adminHandlers.TaskConfigUpdate)
	taskConfigGroup.POST("/delete", adminHandlers.TaskConfigRemove)
	taskConfigGroup.POST("/update_sort", adminHandlers.TaskConfigUpdateSort)
	taskConfigGroup.GET("/:id", adminHandlers.TaskConfigDetail)

	badgeGroup := group.Group("/badge")
	badgeGroup.GET("/list", adminHandlers.BadgeList)
	badgeGroup.POST("/list", adminHandlers.BadgeList)
	badgeGroup.POST("/create", adminHandlers.BadgeCreate)
	badgeGroup.POST("/update", adminHandlers.BadgeUpdate)
	badgeGroup.POST("/delete", adminHandlers.BadgeRemove)
	badgeGroup.POST("/update_sort", adminHandlers.BadgeUpdateSort)
	badgeGroup.GET("/:id", adminHandlers.BadgeDetail)

	levelConfigGroup := group.Group("/level-config")
	levelConfigGroup.POST("/list", adminHandlers.LevelConfigList)
	levelConfigGroup.POST("/save_all", adminHandlers.LevelConfigSaveAll)
	levelConfigGroup.GET("/:id", adminHandlers.LevelConfigDetail)

	userTaskLogGroup := group.Group("/user-task-log")
	userTaskLogGroup.POST("/list", adminHandlers.UserTaskLogList)
	userTaskLogGroup.GET("/:id", adminHandlers.UserTaskLogDetail)

	userExpLogGroup := group.Group("/user-exp-log")
	userExpLogGroup.POST("/list", adminHandlers.UserExpLogList)
	userExpLogGroup.GET("/:id", adminHandlers.UserExpLogDetail)

	userBadgeGroup := group.Group("/user-badge")
	userBadgeGroup.POST("/list", adminHandlers.UserBadgeList)
	userBadgeGroup.GET("/:id", adminHandlers.UserBadgeDetail)

	operateLogGroup := group.Group("/operate-log")
	operateLogGroup.POST("/list", adminHandlers.OperateLogList)
	operateLogGroup.GET("/:id", adminHandlers.OperateLogDetail)

	userReportGroup := group.Group("/user-report")
	userReportGroup.POST("/list", adminHandlers.UserReportList)
	userReportGroup.POST("/create", adminHandlers.UserReportCreate)
	userReportGroup.POST("/update", adminHandlers.UserReportUpdate)
	userReportGroup.POST("/audit", adminHandlers.UserReportAudit)
	userReportGroup.GET("/:id", adminHandlers.UserReportDetail)

	forbiddenWordGroup := group.Group("/forbidden-word")
	forbiddenWordGroup.POST("/list", adminHandlers.ForbiddenWordList)
	forbiddenWordGroup.POST("/create", adminHandlers.ForbiddenWordCreate)
	forbiddenWordGroup.POST("/update", adminHandlers.ForbiddenWordUpdate)
	forbiddenWordGroup.POST("/delete", adminHandlers.ForbiddenWordRemove)
	forbiddenWordGroup.GET("/:id", adminHandlers.ForbiddenWordDetail)

	voteGroup := group.Group("/vote")
	voteGroup.POST("/list", adminHandlers.VoteList)
	voteGroup.POST("/create", adminHandlers.VoteCreate)
	voteGroup.POST("/update", adminHandlers.VoteUpdate)
	voteGroup.POST("/delete", adminHandlers.VoteRemove)
	voteGroup.GET("/:id", adminHandlers.VoteDetail)

	voteOptionGroup := group.Group("/vote-option")
	voteOptionGroup.POST("/list", adminHandlers.VoteOptionList)
	voteOptionGroup.POST("/create", adminHandlers.VoteOptionCreate)
	voteOptionGroup.POST("/update", adminHandlers.VoteOptionUpdate)
	voteOptionGroup.POST("/delete", adminHandlers.VoteOptionRemove)
	voteOptionGroup.GET("/:id", adminHandlers.VoteOptionDetail)

	voteRecordGroup := group.Group("/vote-record")
	voteRecordGroup.POST("/list", adminHandlers.VoteRecordList)
	voteRecordGroup.POST("/create", adminHandlers.VoteRecordCreate)
	voteRecordGroup.POST("/update", adminHandlers.VoteRecordUpdate)
	voteRecordGroup.POST("/delete", adminHandlers.VoteRecordRemove)
	voteRecordGroup.GET("/:id", adminHandlers.VoteRecordDetail)

}
