package model

import "database/sql"

var Models = []interface{}{
	&User{}, &UserToken{}, &Category{}, &Tag{}, &Article{}, &ArticleTag{}, &Comment{}, &Favorite{},
	&Topic{}, &TopicTag{}, &TopicLike{}, &Message{}, &SysConfig{}, &Project{}, &Subject{}, &SubjectContent{}, &Link{},
	&CollectRule{}, &CollectArticle{}, &ThirdAccount{},
}

type Model struct {
	Id int64 `gorm:"PRIMARY_KEY;AUTO_INCREMENT" json:"id" form:"id"`
}

const (
	UserStatusOk       = 0
	UserStatusDisabled = 1

	UserTokenStatusOk       = 0
	UserTokenStatusDisabled = 1

	UserTypeNormal = 0 // 普通用户
	UserTypeGzh    = 1 // 公众号用户

	CategoryStatusOk       = 0
	CategoryStatusDisabled = 1

	TagStatusOk       = 0
	TagStatusDisabled = 1

	ArticleStatusPublished = 0 // 已发布
	ArticleStatusDeleted   = 1 // 已删除
	ArticleStatusDraft     = 2 // 草稿

	ArticleTagStatusOk      = 0
	ArticleTagStatusDeleted = 1

	TopicStatusOk      = 0
	TopicStatusDeleted = 1

	TopicTagStatusOk      = 0
	TopicTagStatusDeleted = 1

	ContentTypeHtml     = "html"
	ContentTypeMarkdown = "markdown"

	CommentStatusOk      = 0
	CommentStatusDeleted = 1

	EntityTypeArticle = "article"
	EntityTypeTopic   = "topic"

	MsgStatusUnread = 0 // 消息未读
	MsgStatusReaded = 1 // 消息已读

	MsgTypeComment = 0 // 回复消息

	LinkStatusOk      = 0 // 正常
	LinkStatusDeleted = 1 // 删除
	LinkStatusPending = 2 // 待审核

	CollectRuleStatusOk       = 0 // 启用
	CollectRuleStatusDisabled = 1 // 禁用

	CollectArticleStatusPending   = 0 // 待审核
	CollectArticleStatusAuditPass = 1 // 审核通过
	CollectArticleStatusAuditFail = 2 // 审核失败
	CollectArticleStatusPublished = 3 // 已发布

	ThirdAccountTypeGithub = "github"
	ThirdAccountTypeQQ     = "qq"
)

type User struct {
	Model
	Username    sql.NullString `gorm:"size:32;unique;" json:"username" form:"username"`
	Email       sql.NullString `gorm:"size:128;unique;" json:"email" form:"email"`
	Nickname    string         `gorm:"size:16;" json:"nickname" form:"nickname"`
	Avatar      string         `gorm:"type:text" json:"avatar" form:"avatar"`
	Password    string         `gorm:"size:512" json:"password" form:"password"`
	Status      int            `gorm:"index:idx_status;not null" json:"status" form:"status"`
	Roles       string         `gorm:"type:text" json:"roles" form:"roles"`
	Type        int            `gorm:"not null" json:"type" form:"type"`
	Description string         `gorm:"type:text" json:"description" form:"description"`
	CreateTime  int64          `json:"createTime" form:"createTime"`
	UpdateTime  int64          `json:"updateTime" form:"updateTime"`
}

type UserToken struct {
	Model
	Token      string `gorm:"size:32;unique;not null" json:"token" form:"token"`
	UserId     int64  `gorm:"not null;index:idx_user_id;" json:"userId" form:"userId"`
	ExpiredAt  int64  `gorm:"not null" json:"expiredAt" form:"expiredAt"`
	Status     int    `gorm:"not null;index:idx_status" json:"status" form:"status"`
	CreateTime int64  `gorm:"not null" json:"createTime" form:"createTime"`
}

type ThirdAccount struct {
	Model
	UserId     sql.NullInt64 `gorm:"unique_index:idx_user_id_third_type;" json:"userId" form:"userId"`                                  // 用户编号
	Avatar     string        `gorm:"size:1024" json:"avatar" form:"avatar"`                                                             // 头像
	Nickname   string        `gorm:"size:32" json:"nickname" form:"nickname"`                                                           // 昵称
	ThirdType  string        `gorm:"size:32;not null;unique_index:idx_user_id_third_type,idx_third;" json:"thirdType" form:"thirdType"` // 第三方类型
	ThirdId    string        `gorm:"size:64;not null;unique_index:idx_third;" json:"thirdId" form:"thirdId"`                            // 第三方唯一标识，例如：openId,unionId
	ExtraData  string        `gorm:"type:longtext" json:"extraData" form:"extraData"`                                                   // 扩展数据
	CreateTime int64         `json:"createTime" form:"createTime"`                                                                      // 创建时间
	UpdateTime int64         `json:"updateTime" form:"updateTime"`                                                                      // 更新时间
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
	ArticleId  int64 `gorm:"not null;index:idx_article_id;" json:"articleId" form:"articleId"` // 文章编号
	TagId      int64 `gorm:"not null;index:idx_tag_id;" json:"tagId" form:"tagId"`             // 标签编号
	Status     int64 `gorm:"not null;index:idx_status" json:"status" form:"status"`            // 状态：正常、删除
	CreateTime int64 `json:"createTime" form:"createTime"`                                     // 创建时间
}

// 评论
type Comment struct {
	Model
	UserId     int64  `gorm:"index:idx_user_id;not null" json:"userId" form:"userId"`             // 用户编号
	EntityType string `gorm:"index:idx_entity_type;not null" json:"entityType" form:"entityType"` // 被评论实体类型
	EntityId   int64  `gorm:"index:idx_entity_id;not null" json:"entityId" form:"entityId"`       // 被评论实体编号
	Content    string `gorm:"type:text;not null" json:"content" form:"content"`                   // 内容
	QuoteId    int64  `gorm:"not null"  json:"quoteId" form:"quoteId"`                            // 引用的评论编号
	Status     int    `gorm:"int;index:idx_status" json:"status" form:"status"`                   // 状态：0：待审核、1：审核通过、2：审核失败、3：已发布
	CreateTime int64  `json:"createTime" form:"createTime"`                                       // 创建时间
}

// 收藏
type Favorite struct {
	Model
	UserId     int64  `gorm:"index:idx_user_id;not null" json:"userId" form:"userId"`                     // 用户编号
	EntityType string `gorm:"index:idx_entity_type;size:32;not null" json:"entityType" form:"entityType"` // 收藏实体类型
	EntityId   int64  `gorm:"index:idx_entity_id;not null" json:"entityId" form:"entityId"`               // 收藏实体编号
	CreateTime int64  `json:"createTime" form:"createTime"`                                               // 创建时间
}

// 主题
type Topic struct {
	Model
	UserId          int64  `gorm:"not null;index:idx_user_id;" json:"userId" form:"userId"` // 用户
	Title           string `gorm:"size:128" json:"title" form:"title"`                      // 标题
	Content         string `gorm:"type:longtext" json:"content" form:"content"`             // 内容
	ViewCount       int64  `gorm:"not null" json:"viewCount" form:"viewCount"`              // 查看数量
	CommentCount    int64  `gorm:"not null" json:"commentCount" form:"commentCount"`        // 跟帖数量
	LikeCount       int64  `gorm:"not null" json:"likeCount" form:"likeCount"`              // 点赞数量
	Status          int    `gorm:"index:idx_status;" json:"status" form:"status"`           // 状态：0：正常、1：删除
	LastCommentTime int64  `json:"lastCommentTime" form:"lastCommentTime"`                  // 最后回复时间
	CreateTime      int64  `json:"createTime" form:"createTime"`                            // 创建时间
	ExtraData       string `gorm:"type:text" json:"extraData" form:"extraData"`             // 扩展数据
}

// 主题标签
type TopicTag struct {
	Model
	TopicId    int64 `gorm:"not null;index:idx_topic_id;" json:"topicId" form:"topicId"` // 主题编号
	TagId      int64 `gorm:"not null;index:idx_tag_id;" json:"tagId" form:"tagId"`       // 标签编号
	Status     int64 `gorm:"not null;index:idx_status" json:"status" form:"status"`      // 状态：正常、删除
	CreateTime int64 `json:"createTime" form:"createTime"`                               // 创建时间
}

// 话题点赞
type TopicLike struct {
	Model
	UserId     int64 `gorm:"not null;index:idx_user_id;" json:"userId" form:"userId"`    // 用户
	TopicId    int64 `gorm:"not null;index:idx_topic_id;" json:"topicId" form:"topicId"` // 主题编号
	CreateTime int64 `json:"createTime" form:"createTime"`                               // 创建时间
}

// 消息
type Message struct {
	Model
	FromId       int64  `gorm:"not null" json:"fromId" form:"fromId"`                    // 消息发送人
	UserId       int64  `gorm:"not null;index:idx_user_id;" json:"userId" form:"userId"` // 用户编号(消息接收人)
	Content      string `gorm:"type:text;not null" json:"content" form:"content"`        // 消息内容
	QuoteContent string `gorm:"type:text" json:"quoteContent" form:"quoteContent"`       // 引用内容
	Type         int    `gorm:"not null" json:"type" form:"type"`                        // 消息类型
	ExtraData    string `gorm:"type:text" json:"extraData" form:"extraData"`             // 扩展数据
	Status       int    `gorm:"not null" json:"status" form:"status"`                    // 状态：0：未读、1：已读
	CreateTime   int64  `json:"createTime" form:"createTime"`                            // 创建时间
}

// 系统配置
type SysConfig struct {
	Model
	Key         string `gorm:"not null;size:128;unique" json:"key" form:"key"` // 配置key
	Value       string `gorm:"type:text" json:"value" form:"value"`            // 配置值
	Name        string `gorm:"not null;size:32" json:"name" form:"name"`       // 配置名称
	Description string `gorm:"size:128" json:"description" form:"description"` // 配置描述
	CreateTime  int64  `gorm:"not null" json:"createTime" form:"createTime"`   // 创建时间
	UpdateTime  int64  `gorm:"not null" json:"updateTime" form:"updateTime"`   // 更新时间
}

// 开源项目
type Project struct {
	Model
	UserId      int64  `gorm:"not null" json:"userId" form:"userId"`
	Name        string `gorm:"type:varchar(1024)" json:"name" form:"name"`
	Title       string `gorm:"type:text" json:"title" form:"title"`
	Logo        string `gorm:"type:varchar(1024)" json:"logo" form:"logo"`
	Url         string `gorm:"type:varchar(1024)" json:"url" form:"url"`
	DocUrl      string `gorm:"type:varchar(1024)" json:"docUrl" form:"docUrl"`
	DownloadUrl string `gorm:"type:varchar(1024)" json:"downloadUrl" form:"downloadUrl"`
	ContentType string `gorm:"type:varchar(32);" json:"contentType" form:"contentType"`
	Content     string `gorm:"type:longtext" json:"content" form:"content"`
	CreateTime  int64  `json:"createTime" form:"createTime"`
}

// 专栏
type Subject struct {
	Model
	Title       string `gorm:"not null" json:"title" form:"title"`             // 标题
	Description string `gorm:"not null" json:"description" form:"description"` // 描述
	Logo        string `json:"logo" form:"logo"`                               // LOGO
	Status      int    `gorm:"not null" json:"status" form:"status"`           // 状态，0：正常，1：删除
	CreateTime  int64  `gorm:"not null" json:"createTime" form:"createTime"`   // 创建时间
}

// 专栏内容
type SubjectContent struct {
	Model
	SubjectId  int64  `gorm:"not null;index:idx_subject_id" json:"subjectId" form:"subjectId"`        // 专栏编号
	EntityType string `gorm:"not null;size:32;index:idx_entity;" json:"entityType" form:"entityType"` // 实体类型
	EntityId   int64  `gorm:"not null;index:idx_entity;" json:"entityId" form:"entityId"`             // 实体编号
	Title      string `gorm:"not null" json:"title" form:"title"`                                     // 标题
	Summary    string `gorm:"not null;size:1024" json:"summary" form:"summary"`                       // 描述
	Deleted    bool   `gorm:"not null;index:idx_deleted" json:"deleted" form:"deleted"`               // 是否删除
	CreateTime int64  `gorm:"not null" json:"createTime" form:"createTime"`                           // 创建时间
}

type Link struct {
	Model
	Url        string `gorm:"not null;type:text" json:"url" form:"url"`     // 链接
	Title      string `gorm:"not null;size:128" json:"title" form:"title"`  // 标题
	Summary    string `gorm:"size:1024" json:"summary" form:"summary"`      // 站点描述
	Logo       string `gorm:"type:text" json:"logo" form:"logo"`            // LOGO
	Category   string `gorm:"type:text" json:"category" form:"category"`    // 分类
	Status     int    `gorm:"not null" json:"status" form:"status"`         // 状态
	Score      int    `gorm:"not null" json:"score" form:"score"`           // 评分，0-100分，分数越高越优质
	Remark     string `gorm:"size:1024" json:"remark" form:"remark"`        // 备注，后台填写的
	CreateTime int64  `gorm:"not null" json:"createTime" form:"createTime"` // 创建时间
}

// 采集规则
type CollectRule struct {
	Model
	Title       string `gorm:"not null;size:128" json:"title" form:"title"`     // 标题
	Rule        string `gorm:"type:text;not null" json:"rule" form:"rule"`      // 规则详情，JSON文件
	Status      int    `gorm:"not null" json:"status" form:"status"`            // 状态
	Description string `gorm:"size:1024" json:"description" form:"description"` // 规则描述
	CreateTime  int64  `json:"createTime" form:"createTime"`                    // 创建时间
	UpdateTime  int64  `json:"updateTime" form:"updateTime"`                    // 更新时间
}

// 采集文章
type CollectArticle struct {
	Model
	UserId         int64  `gorm:"index:idx_user_id;not null" json:"userId" form:"userId"`          // 用户编号
	RuleId         int64  `gorm:"index:idx_rule_id;not null" json:"ruleId" form:"ruleId"`          // 采集规则编号
	LinkId         int64  `gorm:"not null;index:idx_link_id" json:"linkId" form:"linkId"`          // CollectLink外键
	Title          string `gorm:"size:128;not null;" json:"title" form:"title"`                    // 标题
	Summary        string `gorm:"type:text" json:"summary" form:"summary"`                         // 摘要
	Content        string `gorm:"type:longtext;not null;" json:"content" form:"content"`           // 内容
	ContentType    string `gorm:"type:varchar(32);not null" json:"contentType" form:"contentType"` // 内容类型：markdown、html
	SourceUrl      string `gorm:"type:text" json:"sourceUrl" form:"sourceUrl"`                     // 原文链接
	SourceId       string `gorm:"size:1024" json:"sourceId" form:"sourceId"`                       // 原id
	SourceUrlMd5   string `gorm:"index:idx_url_md5" json:"sourceUrlMd5" form:"sourceUrlMd5"`       // url md5
	SourceTitleMd5 string `gorm:"index:idx_title_md5" json:"sourceTitleMd5" form:"sourceTitleMd5"` // 标题 md5
	Status         int    `gorm:"int" json:"status" form:"status"`                                 // 状态：0：待审核、1：审核通过、2：审核失败、3：已发布
	ArticleId      int64  `gorm:"index:idx_article_id;not null" json:"articleId" form:"articleId"` // 发布之后的文章id
	CreateTime     int64  `json:"createTime" form:"createTime"`                                    // 创建时间
}
