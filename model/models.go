package model

var Models = []interface{}{
	&User{}, &GithubUser{}, &Category{}, &Tag{}, &Article{}, &ArticleTag{}, &Comment{}, &Favorite{},
	&UserArticleTag{}, &Topic{}, &TopicTag{}, &Message{}, &OauthClient{}, &OauthToken{},
}

type Model struct {
	Id int64 `gorm:"PRIMARY_KEY;AUTO_INCREMENT" json:"id" form:"id"`
}

const (
	UserStatusOk       = 0
	UserStatusDisabled = 1

	UserTypeNormal = 0 // 普通用户
	UserTypeGzh    = 1 // 公众号用户

	CategoryStatusOk       = 0
	CategoryStatusDisabled = 1

	TagStatusOk       = 0
	TagStatusDisabled = 1

	ArticleStatusPublished = 0 // 已发布
	ArticleStatusDeleted   = 1 // 已删除
	ArticleStatusDraft     = 2 // 草稿

	TopicStatusOk      = 0
	TopicStatusDeleted = 1

	ArticleContentTypeHtml     = "html"
	ArticleContentTypeMarkdown = "markdown"

	CommentStatusOk      = 0
	CommentStatusDeleted = 1

	EntityTypeArticle = "article"
	EntityTypeTopic   = "topic"

	MsgStatusUnread = 0 // 消息未读
	MsgStatusReaded = 1 // 消息已读

	MsgTypeComment = 0 // 回复消息
)

type User struct {
	Model
	Username    string `gorm:"size:32;unique" json:"username" form:"username"`
	Nickname    string `gorm:"size:16" json:"nickname" form:"nickname"`
	Avatar      string `gorm:"type:text" json:"avatar" form:"avatar"`
	Email       string `gorm:"size:512;" json:"email" form:"email"`
	Password    string `gorm:"size:512" json:"password" form:"password"`
	Status      int    `gorm:"index:idx_status;not null" json:"status" form:"status"`
	Roles       string `gorm:"type:text" json:"roles" form:"roles"`
	Type        int    `gorm:"not null" json:"type" form:"type"`
	Description string `gorm:"type:text" json:"description" form:"description"`
	CreateTime  int64  `json:"createTime" form:"createTime"`
	UpdateTime  int64  `json:"updateTime" form:"updateTime"`
}

type GithubUser struct {
	Model
	UserId     int64  `json:"userId" form:"userId"`
	GithubId   int64  `gorm:"unique" json:"githubId" form:"githubId"`
	Login      string `gorm:"size:512" json:"login" form:"login"`
	NodeId     string `gorm:"size:512" json:"nodeId" form:"nodeId"`
	AvatarUrl  string `gorm:"size:1024" json:"avatarUrl" form:"avatarUrl"`
	Url        string `gorm:"size:1024" json:"url" form:"url"`
	HtmlUrl    string `gorm:"size:1024" json:"htmlUrl" form:"htmlUrl"`
	Email      string `gorm:"size:512" json:"email" form:"email"`
	Name       string `gorm:"size:32" json:"name" form:"name"`
	Bio        string `gorm:"type:text" json:"bio" form:"bio"`
	Company    string `gorm:"size:128" json:"company" form:"company"`
	Blog       string `gorm:"type:text" json:"blog" form:"blog"`
	Location   string `gorm:"size:128" json:"location" form:"location"`
	CreateTime int64  `json:"createTime" form:"createTime"`
	UpdateTime int64  `json:"updateTime" form:"updateTime"`
}

// 分类
type Category struct {
	Model
	Name        string `gorm:"size:32;unique;not null" json:"name" form:"name"`
	Description string `gorm:"size:1024" json:"description" form:"description"`
	Status      int    `gorm:"index:idx_status;not null" json:"status" form:"status"`
	CreateTime  int64  `json:"createTime" form:"createTime"`
	UpdateTime  int64  `json:"updateTime" form:"updateTime"`
}

// 标签
type Tag struct {
	Model
	CategoryId  int64  `gorm:"index:idx_category_id;not null" json:"categoryId" form:"categoryId"`
	Name        string `gorm:"size:32;unique;not null" json:"name" form:"name"`
	Description string `gorm:"size:1024" json:"description" form:"description"`
	Status      int    `gorm:"index:idx_status;not null" json:"status" form:"status"`
	CreateTime  int64  `json:"createTime" form:"createTime"`
	UpdateTime  int64  `json:"updateTime" form:"updateTime"`
}

// 文章
type Article struct {
	Model
	CategoryId  int64  `gorm:"index:idx_category_id;not null" json:"categoryId" form:"categoryId"` // 分类编号
	UserId      int64  `gorm:"index:idx_user_id" json:"userId" form:"userId"`                      // 所属用户编号
	Title       string `gorm:"size:128;not null;" json:"title" form:"title"`                       // 标题
	Summary     string `gorm:"type:text" json:"summary" form:"summary"`                            // 摘要
	Content     string `gorm:"type:longtext;not null;" json:"content" form:"content"`              // 内容
	ContentType string `gorm:"type:varchar(32);not null" json:"contentType" form:"contentType"`    // 内容类型：markdown、html
	Status      int    `gorm:"int;not null" json:"status" form:"status"`                           // 状态
	Share       bool   `gorm:"not null" json:"share" form:"share"`                                 // 是否是分享的文章，如果是这里只会显示文章摘要，原文需要跳往原链接查看
	SourceUrl   string `gorm:"type:text" json:"sourceUrl" form:"sourceUrl"`                        // 原文链接
	CreateTime  int64  `json:"createTime" form:"createTime"`                                       // 创建时间
	UpdateTime  int64  `json:"updateTime" form:"updateTime"`                                       // 更新时间
}

// 文章标签
type ArticleTag struct {
	Model
	ArticleId  int64 `gorm:"not null" json:"articleId" form:"articleId"` // 文章编号
	TagId      int64 `gorm:"not null" json:"tagId" form:"tagId"`         // 标签编号
	CreateTime int64 `json:"createTime" form:"createTime"`               // 创建时间
}

// 用户文章标签
type UserArticleTag struct {
	Model
	UserId int64 `gorm:"not null" json:"userId" form:"userId"` // 用户编号
	TagId  int64 `gorm:"not null" json:"tagId" form:"tagId"`   // 标签编号
}

// 评论
type Comment struct {
	Model
	UserId     int64  `gorm:"index:idx_user_id;not null" json:"userId" form:"userId"`             // 用户编号
	EntityType string `gorm:"index:idx_entity_type;not null" json:"entityType" form:"entityType"` // 被评论实体类型
	EntityId   int64  `gorm:"index:idx_entity_id;not null" json:"entityId" form:"entityId"`       // 被评论实体编号
	Content    string `gorm:"type:text;not null" json:"content" form:"content"`                   // 内容
	QuoteId    int64  `gorm:"not null"  json:"quoteId" form:"quoteId"`                            // 引用的评论编号
	Status     int    `gorm:"int" json:"status" form:"status"`                                    // 状态：0：待审核、1：审核通过、2：审核失败、3：已发布
	CreateTime int64  `json:"createTime" form:"createTime"`                                       // 创建时间
}

// 收藏
type Favorite struct {
	Model
	UserId     int64  `gorm:"index:idx_user_id;not null" json:"userId" form:"userId"`             // 用户编号
	EntityType string `gorm:"index:idx_entity_type;not null" json:"entityType" form:"entityType"` // 收藏实体类型
	EntityId   int64  `gorm:"index:idx_entity_id;not null" json:"entityId" form:"entityId"`       // 收藏实体编号
	CreateTime int64  `json:"createTime" form:"createTime"`                                       // 创建时间
}

// 主题
type Topic struct {
	Model
	UserId          int64  `gorm:"not null" json:"userId" form:"userId"`        // 用户
	Title           string `gorm:"size:128" json:"title" form:"title"`          // 标题
	Content         string `gorm:"type:longtext" json:"content" form:"content"` // 内容
	ViewCount       int64  `gorm:"not null" json:"viewCount" form:"viewCount"`  // 查看数量
	Status          int    `gorm:"int" json:"status" form:"status"`             // 状态：0：正常、1：删除
	LastCommentTime int64  `json:"lastCommentTime" form:"lastCommentTime"`      // 最后回复时间
	CreateTime      int64  `json:"createTime" form:"createTime"`                // 创建时间
}

// 主题标签
type TopicTag struct {
	Model
	TopicId    int64 `gorm:"not null" json:"topicId" form:"topicId"` // 主题编号
	TagId      int64 `gorm:"not null" json:"tagId" form:"tagId"`     // 标签编号
	CreateTime int64 `json:"createTime" form:"createTime"`           // 创建时间
}

// 消息
type Message struct {
	Model
	FromId       int64  `gorm:"not null" json:"fromId" form:"fromId"`              // 消息发送人
	UserId       int64  `gorm:"not null" json:"userId" form:"userId"`              // 用户编号(消息接收人)
	Content      string `gorm:"type:text;not null" json:"content" form:"content"`  // 消息内容
	QuoteContent string `gorm:"type:text" json:"quoteContent" form:"quoteContent"` // 引用内容
	Type         int    `gorm:"not null" json:"type" form:"type"`                  // 消息类型
	ExtraData    string `gorm:"type:text" json:"extraData" form:"extraData"`       // 扩展数据
	Status       int    `gorm:"not null" json:"status" form:"status"`              // 状态：0：未读、1：已读
	CreateTime   int64  `json:"createTime" form:"createTime"`                      // 创建时间
}
