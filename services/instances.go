package services

var (
	// services
	UserServiceInstance           = NewUserService()
	GithubUserServiceInstance     = NewGithubUserService()
	CategoryServiceInstance       = NewCategoryService()
	TagServiceInstance            = NewTagService()
	ArticleServiceInstance        = NewArticleService()
	CommentServiceInstance        = NewCommentService()
	FavoriteServiceInstance       = NewFavoriteService()
	ArticleTagServiceInstance     = NewArticleTagService()
	UserArticleTagServiceInstance = NewUserArticleTagService()
	TopicServiceInstance          = NewTopicService()
	TopicTagServiceInstance       = NewTopicTagService()
	MessageServiceInstance        = NewMessageService()
	OauthClientServiceInstance    = NewOauthClientService()
	OauthTokenServiceInstance     = NewOauthTokenService()

	Instances = []interface{}{
		UserServiceInstance, GithubUserServiceInstance, CategoryServiceInstance, TagServiceInstance,
		ArticleServiceInstance, CommentServiceInstance, FavoriteServiceInstance, ArticleTagServiceInstance,
		UserArticleTagServiceInstance, TopicServiceInstance, TopicTagServiceInstance, MessageServiceInstance,
		OauthClientServiceInstance, OauthTokenServiceInstance,
	}
)
