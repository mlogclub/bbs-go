package cache

// 包级别缓存变量，直接指向Redis实现
var (
	UserCache        = UserCacheRedis
	UserTokenCache   = UserTokenCacheRedis  
	SysConfigCache   = SysConfigCacheRedis
	TopicCache       = TopicCacheRedis
	TagCache         = TagCacheRedis
	ForbiddenWordCache = ForbiddenWordCacheRedis
	ArticleTagCache  = ArticleTagCacheRedis
)