package models

import (
	"bbs-go/internal/models/constants"
	"database/sql"
	"time"
)

var Models = []interface{}{
	&Migration{},
	&UserRole{}, &Role{}, &Menu{}, &RoleMenu{}, &Api{}, &MenuApi{}, &DictType{}, &Dict{},

	&User{}, &UserToken{}, &ThirdUser{}, &Tag{}, &Article{}, &ArticleTag{}, &Comment{}, &Favorite{}, &Topic{}, &TopicNode{},
	&TopicTag{}, &UserLike{}, &Message{}, &SysConfig{}, &Link{},
	&TaskConfig{}, &UserTaskEvent{}, &UserTaskLog{},
	&Badge{}, &UserBadge{},
	&LevelConfig{},
	&Vote{}, &VoteOption{}, &VoteRecord{},
	&UserScoreLog{}, &UserExpLog{},
	&OperateLog{}, &EmailCode{}, &SmsCode{}, &CheckIn{}, &UserFollow{}, &UserFeed{}, &UserReport{},
	&ForbiddenWord{},
}

type Model struct {
	Id int64 `gorm:"primaryKey;autoIncrement" json:"id" form:"id"`
}

type Migration struct {
	Model
	Version    int64  `json:"version" form:"version" gorm:"unique"`
	Remark     string `json:"remark" form:"remark" gorm:"type:text"`
	Success    bool   `json:"success" form:"success" gorm:"default:false"`
	ErrorInfo  string `json:"errorInfo" form:"errorInfo" gorm:"type:text"`
	RetryCount int    `json:"retryCount" form:"retryCount"`
	CreateTime int64  `json:"createTime" form:"createTime"`
	UpdateTime int64  `json:"updateTime" form:"updateTime"`
}

type Role struct {
	Model
	Type       int    `gorm:"not null;default:1" json:"type" form:"type"`             // 角色类型（0：系统角色、1：自定义角色）
	Name       string `gorm:"size:64" json:"name" form:"name"`                        // 角色名称
	Code       string `gorm:"unique;size:64" json:"code" form:"code"`                 // 角色编码
	SortNo     int    `json:"sortNo" form:"sortNo"`                                   // 排序
	Remark     string `gorm:"size:256" json:"remark" form:"remark"`                   // 备注
	Status     int    `json:"status" form:"status"`                                   // 状态
	CreateTime int64  `gorm:"not null;default:0" json:"createTime" form:"createTime"` // 创建时间
	UpdateTime int64  `gorm:"not null;default:0" json:"updateTime" form:"updateTime"` // 更新时间
}

type Menu struct {
	Model
	ParentId   int64  `json:"parentId" form:"parentId"`                               // 上级菜单
	Type       string `gorm:"size:32" json:"type" form:"type"`                        // 类型（menu/func）
	Name       string `gorm:"size:64" json:"name" form:"name"`                        // 名称
	Title      string `gorm:"size:64" json:"title" form:"title"`                      // 标题
	Icon       string `gorm:"size:1024" json:"icon" form:"icon"`                      // ICON
	Path       string `gorm:"size:1024" json:"path" form:"path"`                      // 路径
	Component  string `gorm:"size:256" json:"component" form:"component"`             // 组件
	SortNo     int    `gorm:"not null;default:0" json:"sortNo" form:"sortNo"`         // 排序
	Status     int    `json:"status" form:"status"`                                   // 状态
	CreateTime int64  `gorm:"not null;default:0" json:"createTime" form:"createTime"` // 创建时间
	UpdateTime int64  `gorm:"not null;default:0" json:"updateTime" form:"updateTime"` // 更新时间
}

// MenuApi 菜单和接口的权限关联
type MenuApi struct {
	Model
	MenuId     int64 `gorm:"not null;default:0;uniqueIndex:idx_menu_api" json:"menuId" form:"menuId"` // 菜单ID
	ApiId      int64 `gorm:"not null;default:0;uniqueIndex:idx_menu_api" json:"apiId" form:"apiId"`   // 接口ID
	CreateTime int64 `gorm:"not null;default:0" json:"createTime" form:"createTime"`                  // 创建时间
}

type UserRole struct {
	Model
	UserId     int64 `gorm:"uniqueIndex:idx_user_role" json:"userId" form:"userId"`
	RoleId     int64 `gorm:"uniqueIndex:idx_user_role" json:"roleId" form:"roleId"`
	CreateTime int64 `gorm:"not null;default:0" json:"createTime" form:"createTime"` // 创建时间
}

type RoleMenu struct {
	Model
	RoleId     int64 `gorm:"uniqueIndex:idx_role_menu" json:"roleId" form:"roleId"`
	MenuId     int64 `gorm:"uniqueIndex:idx_role_menu" json:"menuId" form:"menuId"`
	CreateTime int64 `gorm:"not null;default:0" json:"createTime" form:"createTime"` // 创建时间
}

type Api struct {
	Model
	Name       string `gorm:"size:512;unique" json:"name" form:"name"`                // 名称
	Method     string `gorm:"size:16" json:"method" form:"method"`                    // 方法
	Path       string `gorm:"size:512;unique" json:"path" form:"path"`                // 路径
	CreateTime int64  `gorm:"not null;default:0" json:"createTime" form:"createTime"` // 创建时间
	UpdateTime int64  `gorm:"not null;default:0" json:"updateTime" form:"updateTime"` // 更新时间
}

type DictType struct {
	Model
	Name       string `gorm:"size:32" json:"name" form:"name"`
	Code       string `gorm:"size:64;unique" json:"code" form:"code"`
	Status     int    `gorm:"not null;default:0" json:"status" form:"status"`
	Remark     string `gorm:"size:512" json:"remark" form:"remark"`
	CreateTime int64  `gorm:"not null;default:0" json:"createTime" form:"createTime"` // 创建时间
	UpdateTime int64  `gorm:"not null;default:0" json:"updateTime" form:"updateTime"` // 更新时间
}

type Dict struct {
	Model
	TypeId     int64  `gorm:"uniqueIndex:idx_dict_name" json:"typeId" form:"typeId"`     // 分类
	ParentId   int64  `gorm:"default:0" json:"parentId" form:"parentId"`                 // 上级
	Name       string `gorm:"size:64;uniqueIndex:idx_dict_name" json:"name" form:"name"` // 名称
	Label      string `gorm:"size:64" json:"label" form:"label"`                         // Label
	Value      string `gorm:"type:text" json:"value" form:"value"`                       // Value
	SortNo     int    `gorm:"not null;default:0" json:"sortNo" form:"sortNo"`            // 排序
	Status     int    `gorm:"not null;default:0" json:"status" form:"status"`            // 状态
	CreateTime int64  `gorm:"not null;default:0" json:"createTime" form:"createTime"`    // 创建时间
	UpdateTime int64  `gorm:"not null;default:0" json:"updateTime" form:"updateTime"`    // 更新时间
}

type User struct {
	Model
	Type             int              `gorm:"not null;default:0" json:"type" form:"type"`                              // 用户类型（0：用户、1：员工）
	Phone            sql.NullString   `gorm:"size:16;unique;" json:"phone" form:"phone"`                               // 电话
	Username         sql.NullString   `gorm:"size:32;unique;" json:"username" form:"username"`                         // 用户名
	Email            sql.NullString   `gorm:"size:128;unique;" json:"email" form:"email"`                              // 邮箱
	EmailVerified    bool             `gorm:"not null;default:false" json:"emailVerified" form:"emailVerified"`        // 邮箱是否验证
	Nickname         string           `gorm:"size:16;" json:"nickname" form:"nickname"`                                // 昵称
	Avatar           string           `gorm:"type:text" json:"avatar" form:"avatar"`                                   // 头像
	Gender           constants.Gender `gorm:"size:16;default:''" json:"gender" form:"gender"`                          // 性别
	Birthday         *time.Time       `json:"birthday" form:"birthday"`                                                // 生日
	BackgroundImage  string           `gorm:"type:text" json:"backgroundImage" form:"backgroundImage"`                 // 个人中心背景图片
	Password         string           `gorm:"size:512" json:"password" form:"password"`                                // 密码
	HomePage         string           `gorm:"size:1024" json:"homePage" form:"homePage"`                               // 个人主页
	Description      string           `gorm:"type:text" json:"description" form:"description"`                         // 个人描述
	Score            int              `gorm:"type:int(11);not null;index:idx_user_score" json:"score" form:"score"`    // 积分
	Exp              int              `gorm:"type:int(11);not null;default:0" json:"exp" form:"exp"`                   // 经验
	Level            int              `gorm:"type:int(11);not null;default:1" json:"level" form:"level"`               // 等级（从 1 开始）
	Status           int              `gorm:"type:int(11);index:idx_user_status;not null" json:"status" form:"status"` // 状态
	TopicCount       int              `gorm:"type:int(11);not null" json:"topicCount" form:"topicCount"`               // 帖子数量
	CommentCount     int              `gorm:"type:int(11);not null" json:"commentCount" form:"commentCount"`           // 跟帖数量
	FollowCount      int              `gorm:"type:int(11);not null" json:"followCount" form:"followCount"`             // 关注数量
	FansCount        int              `gorm:"type:int(11);not null" json:"fansCount" form:"fansCount"`                 // 粉丝数量
	Roles            string           `gorm:"type:text" json:"roles" form:"roles"`                                     // 角色
	ForbiddenEndTime int64            `gorm:"not null;default:0" json:"forbiddenEndTime" form:"forbiddenEndTime"`      // 禁言结束时间
	CreateTime       int64            `json:"createTime" form:"createTime"`                                            // 创建时间
	UpdateTime       int64            `json:"updateTime" form:"updateTime"`                                            // 更新时间
}

type UserToken struct {
	Model
	Token      string `gorm:"size:32;unique;not null" json:"token" form:"token"`
	UserId     int64  `gorm:"not null;index:idx_user_token_user_id;" json:"userId" form:"userId"`
	ExpiredAt  int64  `gorm:"not null" json:"expiredAt" form:"expiredAt"`
	Status     int    `gorm:"type:int(11);not null;index:idx_user_token_status" json:"status" form:"status"`
	CreateTime int64  `gorm:"not null" json:"createTime" form:"createTime"`
}

type ThirdUser struct {
	Model
	UserId     int64               `gorm:"not null;uniqueIndex:idx_third_user_user_id" json:"userId" form:"userId"`
	OpenId     string              `gorm:"size:64;not null;uniqueIndex:idx_open_id" json:"openId" form:"openId"`
	ThirdType  constants.ThirdType `gorm:"size:32;not null;uniqueIndex:idx_open_id;uniqueIndex:idx_third_user_user_id" json:"thirdType" form:"thirdType"`
	Nickname   string              `gorm:"size:32" json:"nickname" form:"nickname"`
	Avatar     string              `gorm:"size:1024" json:"avatar" form:"avatar"`
	ExtraData  string              `gorm:"type:longtext" json:"extraData" form:"extraData"`
	CreateTime int64               `json:"createTime" form:"createTime"`
	UpdateTime int64               `json:"updateTime" form:"updateTime"`
}

// 标签
type Tag struct {
	Model
	Name        string `gorm:"size:32;unique;not null" json:"name" form:"name"`
	Description string `gorm:"size:1024" json:"description" form:"description"`
	Status      int    `gorm:"type:int(11);index:idx_tag_status;not null" json:"status" form:"status"`
	CreateTime  int64  `json:"createTime" form:"createTime"`
	UpdateTime  int64  `json:"updateTime" form:"updateTime"`
}

// 文章
type Article struct {
	Model
	UserId       int64                 `gorm:"index:idx_article_user_id" json:"userId" form:"userId"`             // 所属用户编号
	Title        string                `gorm:"size:128;not null;" json:"title" form:"title"`                      // 标题
	Summary      string                `gorm:"type:text" json:"summary" form:"summary"`                           // 摘要
	Content      string                `gorm:"type:longtext;not null;" json:"content" form:"content"`             // 内容
	ContentType  constants.ContentType `gorm:"type:varchar(32);not null" json:"contentType" form:"contentType"`   // 内容类型：markdown、html
	Cover        string                `gorm:"type:text;" json:"cover" form:"cover"`                              // 封面图
	Status       int                   `gorm:"type:int(11);index:idx_article_status" json:"status" form:"status"` // 状态
	SourceUrl    string                `gorm:"type:text" json:"sourceUrl" form:"sourceUrl"`                       // 原文链接
	ViewCount    int64                 `gorm:"not null;" json:"viewCount" form:"viewCount"`                       // 查看数量
	CommentCount int64                 `gorm:"default:0" json:"commentCount" form:"commentCount"`                 // 评论数量
	LikeCount    int64                 `gorm:"default:0" json:"likeCount" form:"likeCount"`                       // 点赞数量
	CreateTime   int64                 `json:"createTime" form:"createTime"`                                      // 创建时间
	UpdateTime   int64                 `json:"updateTime" form:"updateTime"`                                      // 更新时间
}

// 文章标签
type ArticleTag struct {
	Model
	ArticleId  int64 `gorm:"not null;index:idx_article_id;" json:"articleId" form:"articleId"`  // 文章编号
	TagId      int64 `gorm:"not null;index:idx_article_tag_tag_id;" json:"tagId" form:"tagId"`  // 标签编号
	Status     int64 `gorm:"not null;index:idx_article_tag_status" json:"status" form:"status"` // 状态：正常、删除
	CreateTime int64 `json:"createTime" form:"createTime"`                                      // 创建时间
}

// 评论
type Comment struct {
	Model
	UserId       int64                 `gorm:"index:idx_comment_user_id;not null" json:"userId" form:"userId"`                     // 用户编号
	EntityType   string                `gorm:"size:64;index:idx_comment_entity_type;not null" json:"entityType" form:"entityType"` // 被评论实体类型
	EntityId     int64                 `gorm:"index:idx_comment_entity_id;not null" json:"entityId" form:"entityId"`               // 被评论实体编号
	Content      string                `gorm:"type:text;not null" json:"content" form:"content"`                                   // 内容
	ImageList    string                `gorm:"type:longtext" json:"imageList" form:"imageList"`                                    // 图片
	ContentType  constants.ContentType `gorm:"type:varchar(32);not null" json:"contentType" form:"contentType"`                    // 内容类型：markdown、html
	QuoteId      int64                 `gorm:"not null"  json:"quoteId" form:"quoteId"`                                            // 引用的评论编号
	LikeCount    int64                 `gorm:"not null;default:0" json:"likeCount" form:"likeCount"`                               // 点赞数量
	CommentCount int64                 `gorm:"not null;default:0" json:"commentCount" form:"commentCount"`                         // 评论数量
	UserAgent    string                `gorm:"size:1024" json:"userAgent" form:"userAgent"`                                        // UserAgent
	Ip           string                `gorm:"size:128" json:"ip" form:"ip"`                                                       // IP
	IpLocation   string                `gorm:"size:64" json:"ipLocation" form:"ipLocation"`                                        // IP属地
	Status       int                   `gorm:"type:int(11);index:idx_comment_status" json:"status" form:"status"`                  // 状态：0：待审核、1：审核通过、2：审核失败、3：已发布
	CreateTime   int64                 `json:"createTime" form:"createTime"`                                                       // 创建时间
}

// 收藏
type Favorite struct {
	Model
	UserId     int64  `gorm:"index:idx_favorite_user_id;not null" json:"userId" form:"userId"`                     // 用户编号
	EntityType string `gorm:"size:32;index:idx_favorite_entity_type;not null" json:"entityType" form:"entityType"` // 收藏实体类型
	EntityId   int64  `gorm:"index:idx_favorite_entity_id;not null" json:"entityId" form:"entityId"`               // 收藏实体编号
	CreateTime int64  `json:"createTime" form:"createTime"`                                                        // 创建时间
}

// TopicNode 话题节点
type TopicNode struct {
	Model
	Name        string `gorm:"size:32;unique" json:"name" form:"name"`                     // 名称
	Description string `gorm:"size:1024" json:"description" form:"description"`            // 描述
	Logo        string `gorm:"size:1024" json:"logo" form:"logo"`                          // 图标
	SortNo      int    `gorm:"type:int(11);index:idx_sort_no" json:"sortNo" form:"sortNo"` // 排序编号
	Status      int    `gorm:"type:int(11);not null" json:"status" form:"status"`          // 状态
	CreateTime  int64  `json:"createTime" form:"createTime"`                               // 创建时间
}

// 话题节点
type Topic struct {
	Model
	Type              constants.TopicType   `gorm:"type:int(11);not null:default:0" json:"type" form:"type"`                         // 类型
	NodeId            int64                 `gorm:"not null;index:idx_node_id;" json:"nodeId" form:"nodeId"`                         // 节点编号
	UserId            int64                 `gorm:"not null;index:idx_topic_user_id;" json:"userId" form:"userId"`                   // 用户
	Title             string                `gorm:"size:128" json:"title" form:"title"`                                              // 标题
	ContentType       constants.ContentType `gorm:"size:32;default:markdown" json:"contentType" form:"contentType"`                  // 内容类型（html/markdown）
	Content           string                `gorm:"type:longtext" json:"content" form:"content"`                                     // 内容
	ImageList         string                `gorm:"type:longtext" json:"imageList" form:"imageList"`                                 // 图片
	HideContent       string                `gorm:"type:longtext" json:"hideContent" form:"hideContent"`                             // 回复可见内容
	VoteId            int64                 `gorm:"not null;default:0" json:"voteId" form:"voteId"`                                  // 投票ID
	Recommend         bool                  `gorm:"not null;index:idx_recommend" json:"recommend" form:"recommend"`                  // 是否推荐
	RecommendTime     int64                 `gorm:"not null" json:"recommendTime" form:"recommendTime"`                              // 推荐时间
	Sticky            bool                  `gorm:"not null;index:idx_sticky_sticky_time" json:"sticky" form:"sticky"`               // 置顶
	StickyTime        int64                 `gorm:"not null;index:idx_sticky_sticky_time" json:"stickyTime" form:"stickyTime"`       // 置顶时间
	ViewCount         int64                 `gorm:"not null" json:"viewCount" form:"viewCount"`                                      // 查看数量
	CommentCount      int64                 `gorm:"not null" json:"commentCount" form:"commentCount"`                                // 跟帖数量
	LikeCount         int64                 `gorm:"not null" json:"likeCount" form:"likeCount"`                                      // 点赞数量
	Status            int                   `gorm:"type:int(11);index:idx_topic_status;" json:"status" form:"status"`                // 状态：0：正常、1：删除
	LastCommentTime   int64                 `gorm:"index:idx_topic_last_comment_time" json:"lastCommentTime" form:"lastCommentTime"` // 最后回复时间
	LastCommentUserId int64                 `json:"lastCommentUserId" form:"lastCommentUserId"`                                      // 最后回复用户
	UserAgent         string                `gorm:"size:1024" json:"userAgent" form:"userAgent"`                                     // UserAgent
	Ip                string                `gorm:"size:128" json:"ip" form:"ip"`                                                    // IP
	IpLocation        string                `gorm:"size:64" json:"ipLocation" form:"ipLocation"`                                     // IP属地
	CreateTime        int64                 `gorm:"index:idx_topic_create_time" json:"createTime" form:"createTime"`                 // 创建时间
	ExtraData         string                `gorm:"type:text" json:"extraData" form:"extraData"`                                     // 扩展数据
}

// Vote 投票
type Vote struct {
	Model
	Type        constants.VoteType `json:"type" form:"type" redis:"type"`                                          // 投票类型(1:单选 / 2:多选)
	Title       string             `gorm:"size:128" json:"title" form:"title" redis:"title"`                       // 标题
	ExpiredAt   int64              `gorm:"not null" json:"expiredAt" form:"expiredAt" redis:"expiredAt"`           // 截止日期
	TopicId     int64              `gorm:"not null" json:"topicId" form:"topicId" redis:"topicId"`                 // 帖子ID
	UserId      int64              `gorm:"not null" json:"userId" form:"userId" redis:"userId"`                    // 用户ID
	VoteNum     int                `gorm:"not null" json:"voteNum" form:"voteNum" redis:"voteNum"`                 // 可投票数量
	OptionCount int                `gorm:"not null" json:"optionCount" form:"optionCount" redis:"optionCount"`     // 选项数量
	VoteCount   int                `gorm:"not null;default:0" json:"voteCount" form:"voteCount" redis:"voteCount"` // 投票数量
	CreateTime  int64              `gorm:"not null" json:"createTime" form:"createTime" redis:"createTime"`        // 创建时间
}

// VoteOption 投票选项
type VoteOption struct {
	Model
	VoteId     int64  `gorm:"not null;index:idx_vote_id" json:"voteId" form:"voteId" redis:"voteId"`  // 投票ID
	Content    string `gorm:"size:256" json:"content" form:"content" redis:"content"`                 // 选项内容
	SortNo     int    `gorm:"not null" json:"sortNo" form:"sortNo" redis:"sortNo"`                    // 排序
	VoteCount  int    `gorm:"not null;default:0" json:"voteCount" form:"voteCount" redis:"voteCount"` // 票数
	CreateTime int64  `gorm:"not null" json:"createTime" form:"createTime" redis:"createTime"`        // 创建时间
}

// VoteRecord 投票记录
type VoteRecord struct {
	Model
	UserId     int64  `gorm:"uniqueIndex:idx_user_vote" json:"userId" form:"userId"` // 用户ID
	VoteId     int64  `gorm:"uniqueIndex:idx_user_vote" json:"voteId" form:"voteId"` // 投票ID
	OptionIds  string `gorm:"type:text" json:"optionIds" form:"optionIds"`           // 选项ID列表，逗号分隔
	CreateTime int64  `json:"createTime" form:"createTime"`                          // 投票时间
}

// 主题标签
type TopicTag struct {
	Model
	TopicId           int64 `gorm:"not null;index:idx_topic_tag_topic_id;" json:"topicId" form:"topicId"`                // 主题编号
	TagId             int64 `gorm:"not null;index:idx_topic_tag_tag_id;" json:"tagId" form:"tagId"`                      // 标签编号
	Status            int64 `gorm:"not null;index:idx_topic_tag_status" json:"status" form:"status"`                     // 状态：正常、删除
	LastCommentTime   int64 `gorm:"index:idx_topic_tag_last_comment_time" json:"lastCommentTime" form:"lastCommentTime"` // 最后回复时间
	LastCommentUserId int64 `json:"lastCommentUserId" form:"lastCommentUserId"`                                          // 最后回复用户
	CreateTime        int64 `json:"createTime" form:"createTime"`                                                        // 创建时间
}

// 用户点赞
type UserLike struct {
	Model
	UserId     int64  `gorm:"not null;uniqueIndex:idx_user_like_unique;" json:"userId" form:"userId"`                                            // 用户
	EntityId   int64  `gorm:"not null;uniqueIndex:idx_user_like_unique;index:idx_user_like_entity;" json:"topicId" form:"topicId"`               // 实体编号
	EntityType string `gorm:"not null;size:32;uniqueIndex:idx_user_like_unique;index:idx_user_like_entity;" json:"entityType" form:"entityType"` // 实体类型
	CreateTime int64  `json:"createTime" form:"createTime"`                                                                                      // 创建时间
}

// 消息
type Message struct {
	Model
	FromId       int64  `gorm:"not null" json:"fromId" form:"fromId"`                            // 消息发送人
	UserId       int64  `gorm:"not null;index:idx_message_user_id;" json:"userId" form:"userId"` // 用户编号(消息接收人)
	Title        string `gorm:"size:1024" json:"title" form:"title"`                             // 消息标题
	Content      string `gorm:"type:text;not null" json:"content" form:"content"`                // 消息内容
	QuoteContent string `gorm:"type:text" json:"quoteContent" form:"quoteContent"`               // 引用内容
	Type         int    `gorm:"type:int(11);not null" json:"type" form:"type"`                   // 消息类型
	ExtraData    string `gorm:"type:text" json:"extraData" form:"extraData"`                     // 扩展数据
	Status       int    `gorm:"type:int(11);not null" json:"status" form:"status"`               // 状态：0：未读、1：已读
	CreateTime   int64  `json:"createTime" form:"createTime"`                                    // 创建时间
}

// 系统配置
type SysConfig struct {
	Model
	Key         string `gorm:"not null;size:128;unique" json:"key" form:"key"`  // 配置key
	Value       string `gorm:"type:text" json:"value" form:"value"`             // 配置值
	Name        string `gorm:"not null;size:128" json:"name" form:"name"`       // 配置名称
	Description string `gorm:"size:1024" json:"description" form:"description"` // 配置描述
	CreateTime  int64  `gorm:"not null" json:"createTime" form:"createTime"`    // 创建时间
	UpdateTime  int64  `gorm:"not null" json:"updateTime" form:"updateTime"`    // 更新时间
}

// 友链
type Link struct {
	Model
	Url        string `gorm:"not null;type:text" json:"url" form:"url"`          // 链接
	Title      string `gorm:"not null;size:128" json:"title" form:"title"`       // 标题
	Summary    string `gorm:"size:1024" json:"summary" form:"summary"`           // 站点描述
	Status     int    `gorm:"type:int(11);not null" json:"status" form:"status"` // 状态
	CreateTime int64  `gorm:"not null" json:"createTime" form:"createTime"`      // 创建时间
}

// TaskConfig 任务配置
type TaskConfig struct {
	Model
	GroupName   constants.TaskGroup `gorm:"size:32;not null;default:'newbie';index:idx_task_config_group_name" json:"groupName" form:"groupName"` // 任务分组
	EventType   string              `gorm:"size:64;not null;index:idx_task_config_event_type" json:"eventType" form:"eventType"`                  // 事件类型
	Title       string              `gorm:"size:64;not null" json:"title" form:"title"`                                                           // 标题（单语言）
	Description string              `gorm:"size:512;not null" json:"description" form:"description"`                                              // 描述（单语言）

	Score   int   `gorm:"type:int(11);not null;default:0" json:"score" form:"score"`    // 完成一次获得积分
	Exp     int   `gorm:"type:int(11);not null;default:0" json:"exp" form:"exp"`        // 完成一次获得经验
	BadgeId int64 `gorm:"type:bigint;not null;default:0" json:"badgeId" form:"badgeId"` // 完成一次授予勋章（0 表示无）

	Period         int `gorm:"type:int(11);not null;default:0;index:idx_task_config_period" json:"period" form:"period"` // 0一次性/1每日...
	MaxFinishCount int `gorm:"type:int(11);not null;default:1" json:"maxFinishCount" form:"maxFinishCount"`              // 周期内最多完成次数
	EventCount     int `gorm:"type:int(11);not null;default:1" json:"eventCount" form:"eventCount"`                      // 多少次事件算完成一次

	BtnName   string `gorm:"size:32" json:"btnName" form:"btnName"`                                  // 按钮文案（单语言）
	ActionUrl string `gorm:"size:1024" json:"actionUrl" form:"actionUrl"`                            // 按钮跳转
	SortNo    int    `gorm:"type:int(11);index:idx_task_config_sort_no" json:"sortNo" form:"sortNo"` // 排序

	StartTime int64 `gorm:"type:bigint;not null;default:0;index:idx_task_config_time" json:"startTime" form:"startTime"` // 生效时间（0 表示立即）
	EndTime   int64 `gorm:"type:bigint;not null;default:0;index:idx_task_config_time" json:"endTime" form:"endTime"`     // 结束时间（0 表示不结束）

	Status     int   `gorm:"type:int(11); not null;default:0;index:idx_task_config_status" json:"status" form:"status"` // 状态
	CreateTime int64 `gorm:"type:bigint;not null" json:"createTime" form:"createTime"`                                  // 创建时间
	UpdateTime int64 `gorm:"type:bigint;not null" json:"updateTime" form:"updateTime"`                                  // 更新时间
}

// UserTaskEvent 用户任务事件累计（UserId + PeriodKey + TaskId 唯一）
type UserTaskEvent struct {
	Model
	UserId          int64 `gorm:"type:bigint;not null;uniqueIndex:uk_user_task_event_upt;index:idx_user_task_event_user_id" json:"userId" form:"userId"`
	PeriodKey       int   `gorm:"type:int(11);not null;uniqueIndex:uk_user_task_event_upt;index:idx_user_task_event_period_key" json:"periodKey" form:"periodKey"` // 一次性=0；每日=yyyyMMdd
	TaskId          int64 `gorm:"type:bigint;not null;uniqueIndex:uk_user_task_event_upt;index:idx_user_task_event_task_id" json:"taskId" form:"taskId"`
	EventCount      int   `gorm:"type:int(11);not null;default:0" json:"eventCount" form:"eventCount"`            // 事件次数（余量）
	TaskFinishCount int   `gorm:"type:int(11); not null;default:0" json:"taskFinishCount" form:"taskFinishCount"` // 已完成次数（周期内）
	CreateTime      int64 `gorm:"type:bigint;not null" json:"createTime" form:"createTime"`                       // 创建时间
	UpdateTime      int64 `gorm:"type:bigint;not null" json:"updateTime" form:"updateTime"`                       // 更新时间
}

// UserTaskLog 用户任务完成/发奖记录（每完成一次生成一条记录）
type UserTaskLog struct {
	Model
	UserId    int64 `gorm:"type:bigint;not null;uniqueIndex:uk_user_task_log_uptf;index:idx_user_task_log_user_id" json:"userId" form:"userId"`           // 用户编号
	PeriodKey int   `gorm:"type:int(11);not null;uniqueIndex:uk_user_task_log_uptf;index:idx_user_task_log_period_key" json:"periodKey" form:"periodKey"` // 一次性=0；每日=yyyyMMdd
	TaskId    int64 `gorm:"type:bigint;not null;uniqueIndex:uk_user_task_log_uptf;index:idx_user_task_log_task_id" json:"taskId" form:"taskId"`           // 任务ID
	FinishNo  int   `gorm:"type:int(11);not null;default:1;uniqueIndex:uk_user_task_log_uptf" json:"finishNo" form:"finishNo"`                            // 周期内第几次完成（1..MaxFinishCount）

	Score   int   `gorm:"type:int(11);not null;default:0" json:"score" form:"score"`    // 发放积分
	Exp     int   `gorm:"type:int(11);not null;default:0" json:"exp" form:"exp"`        // 发放经验
	BadgeId int64 `gorm:"type:bigint;not null;default:0" json:"badgeId" form:"badgeId"` // 发放勋章（0 表示无）

	CreateTime int64 `gorm:"type:bigint;not null" json:"createTime" form:"createTime"` // 创建时间
	UpdateTime int64 `gorm:"type:bigint;not null" json:"updateTime" form:"updateTime"` // 更新时间
}

// Badge 勋章配置
type Badge struct {
	Model
	Name        string `gorm:"size:64;not null;uniqueIndex:uk_badge_name" json:"name" form:"name"` // 稳定标识
	Title       string `gorm:"size:64;not null" json:"title" form:"title"`                         // 标题（单语言）
	Description string `gorm:"size:512" json:"description" form:"description"`                     // 描述（单语言）
	Icon        string `gorm:"size:1024" json:"icon" form:"icon"`                                  // 图标
	SortNo      int    `gorm:"type:int(11);index:idx_badge_sort_no" json:"sortNo" form:"sortNo"`   // 排序
	Status      int    `gorm:"type:int(11);not null;default:0;index:idx_badge_status" json:"status" form:"status"`
	CreateTime  int64  `gorm:"type:bigint;not null" json:"createTime" form:"createTime"`
	UpdateTime  int64  `gorm:"type:bigint;not null" json:"updateTime" form:"updateTime"`
}

// UserBadge 用户勋章（避免重复授予：UserId + BadgeId 唯一）
type UserBadge struct {
	Model
	UserId     int64  `gorm:"type:bigint;not null;uniqueIndex:uk_user_badge_ub;index:idx_user_badge_user_id" json:"userId" form:"userId"`
	BadgeId    int64  `gorm:"type:bigint;not null;uniqueIndex:uk_user_badge_ub;index:idx_user_badge_badge_id" json:"badgeId" form:"badgeId"`
	SourceType string `gorm:"size:32;not null;index:idx_user_badge_source" json:"sourceType" form:"sourceType"` // 例如：task
	SourceId   string `gorm:"size:64;not null;index:idx_user_badge_source" json:"sourceId" form:"sourceId"`     // 例如：UserTaskLog.Id
	IsWorn     bool   `gorm:"not null;default:false" json:"isWorn" form:"isWorn"`                               // 是否佩戴（可选能力）
	SortNo     int    `gorm:"type:int(11)" json:"sortNo" form:"sortNo"`                                         // 佩戴/展示排序（可选）
	CreateTime int64  `gorm:"type:bigint;not null" json:"createTime" form:"createTime"`
	UpdateTime int64  `gorm:"type:bigint;not null" json:"updateTime" form:"updateTime"`
}

// LevelConfig 等级配置（Level -> NeedExp）
type LevelConfig struct {
	Model
	Level      int    `gorm:"type:int(11);not null;uniqueIndex:uk_level_config_level" json:"level" form:"level"` // 等级（必须从 1 开始且连续）
	NeedExp    int    `gorm:"type:int(11);not null" json:"needExp" form:"needExp"`                               // 达到该等级所需累计经验（必须严格递增）
	Title      string `gorm:"size:64" json:"title" form:"title"`                                                 // 等级称号（可选）
	Status     int    `gorm:"type:int(11);not null;default:0;index:idx_level_config_status" json:"status" form:"status"`
	CreateTime int64  `gorm:"type:bigint;not null" json:"createTime" form:"createTime"`
	UpdateTime int64  `gorm:"type:bigint;not null" json:"updateTime" form:"updateTime"`
}

// 用户积分流水
type UserScoreLog struct {
	Model
	UserId      int64  `gorm:"index:idx_user_score_log_user_id" json:"userId" form:"userId"`           // 用户编号
	SourceType  string `gorm:"size:32;index:idx_user_score_score" json:"sourceType" form:"sourceType"` // 积分来源类型
	SourceId    string `gorm:"size:32;index:idx_user_score_score" json:"sourceId" form:"sourceId"`     // 积分来源编号
	Description string `json:"description" form:"description"`                                         // 描述
	Type        int    `gorm:"type:int(11)" json:"type" form:"type"`                                   // 类型(增加、减少)
	Score       int    `gorm:"type:int(11)" json:"score" form:"score"`                                 // 积分
	CreateTime  int64  `json:"createTime" form:"createTime"`                                           // 创建时间
}

// 用户经验流水
type UserExpLog struct {
	Model
	UserId      int64  `gorm:"index:idx_user_exp_log_user_id" json:"userId" form:"userId"`                // 用户编号
	SourceType  string `gorm:"size:32;index:idx_user_exp_log_source" json:"sourceType" form:"sourceType"` // 经验来源类型
	SourceId    string `gorm:"size:64;index:idx_user_exp_log_source" json:"sourceId" form:"sourceId"`     // 经验来源编号
	Description string `json:"description" form:"description"`                                            // 描述
	Type        int    `gorm:"type:int(11)" json:"type" form:"type"`                                      // 类型(增加、减少)
	Exp         int    `gorm:"type:int(11)" json:"exp" form:"exp"`                                        // 经验
	CreateTime  int64  `json:"createTime" form:"createTime"`                                              // 创建时间
}

// 操作日志
type OperateLog struct {
	Model
	UserId      int64  `gorm:"not null;index:idx_operate_log_user_id" json:"userId" form:"userId"`          // 用户编号
	OpType      string `gorm:"not null;index:idx_op_type;size:32" json:"opType" form:"opType"`              // 操作类型
	DataType    string `gorm:"not null;index:idx_operate_log_data;size:32" json:"dataType" form:"dataType"` // 数据类型
	DataId      int64  `gorm:"not null;index:idx_operate_log_data" json:"dataId" form:"dataId" `            // 数据编号
	Description string `gorm:"not null;size:1024" json:"description" form:"description"`                    // 描述
	Ip          string `gorm:"size:128" json:"ip" form:"ip"`                                                // ip地址
	UserAgent   string `gorm:"type:text" json:"userAgent" form:"userAgent"`                                 // UserAgent
	Referer     string `gorm:"type:text" json:"referer" form:"referer"`                                     // Referer
	CreateTime  int64  `json:"createTime" form:"createTime"`                                                // 创建时间
}

// 邮箱验证码
type EmailCode struct {
	Model
	UserId     int64  `gorm:"not null;index:idx_email_code_user_id" json:"userId" form:"userId"` // 用户编号
	Email      string `gorm:"not null;size:128" json:"email" form:"email"`                       // 邮箱
	Code       string `gorm:"not null;size:8" json:"code" form:"code"`                           // 验证码
	Token      string `gorm:"not null;size:32;unique" json:"token" form:"token"`                 // 验证码token
	Title      string `gorm:"size:1024" json:"title" form:"title"`                               // 标题
	Content    string `gorm:"type:text" json:"content" form:"content"`                           // 内容
	Used       bool   `gorm:"not null" json:"used" form:"used"`                                  // 是否使用
	CreateTime int64  `json:"createTime" form:"createTime"`                                      // 创建时间
}

// 短信验证码
type SmsCode struct {
	Model
	SmsId      string `gorm:"size:32;unique" json:"smsId" form:"smsId"`
	Phone      string `gorm:"size:32" json:"phone" form:"phone"`
	Code       string `gorm:"size:16" json:"code" form:"code"`
	ExpireAt   int64  `json:"expireAt" form:"expireAt"`
	Status     int    `json:"status" form:"status"`
	CreateTime int64  `json:"createTime" form:"createTime"`
}

// 签到
type CheckIn struct {
	Model
	UserId          int64 `gorm:"not null;uniqueIndex:idx_check_in_user_id" json:"userId" form:"userId"` // 用户编号
	LatestDayName   int   `gorm:"type:int(11);not null;index:idx_latest" json:"dayName" form:"dayName"`  // 最后一次签到
	ConsecutiveDays int   `gorm:"type:int(11);not null;" json:"consecutiveDays" form:"consecutiveDays"`  // 连续签到天数
	CreateTime      int64 `json:"createTime" form:"createTime"`                                          // 创建时间
	UpdateTime      int64 `gorm:"index:idx_latest" json:"updateTime" form:"updateTime"`                  // 更新时间
}

// UserFollow 粉丝关注
type UserFollow struct {
	Model
	UserId     int64 `gorm:"not null;uniqueIndex:idx_user_follow" json:"userId"`       // 用户编号
	OtherId    int64 `gorm:"not null;uniqueIndex:idx_user_follow" json:"otherId"`      // 对方的ID（被关注用户编号）
	Status     int   `gorm:"type:int(11);not null" json:"status"`                      // 关注状态
	CreateTime int64 `gorm:"type:bigint;not null" json:"createTime" form:"createTime"` // 创建时间
}

// UserFeed 用户信息流
type UserFeed struct {
	Model
	UserId     int64  `gorm:"not null;uniqueIndex:idx_data;index:idx_user_feed_user_id;index:idx_search" json:"userId"`                 // 用户编号
	DataId     int64  `gorm:"not null;uniqueIndex:idx_data;index:idx_data_id" json:"dataId" form:"dataId"`                              // 数据ID
	DataType   string `gorm:"not null;uniqueIndex:idx_data;index:idx_data_id;index:idx_search;size:32" json:"dataType" form:"dataType"` // 数据类型
	AuthorId   int64  `gorm:"not null;index:idx_user_feed_user_id" json:"authorId" form:"authorId"`                                     // 作者编号
	CreateTime int64  `gorm:"type:bigint;not null;index:idx_search" json:"createTime" form:"createTime"`                                // 数据的创建时间
}

// UserReport 用户举报
type UserReport struct {
	Model
	DataId      int64  `json:"dataId" form:"dataId"`           // 举报数据ID
	DataType    string `json:"dataType" form:"dataType"`       // 举报数据类型
	UserId      int64  `json:"userId" form:"userId"`           // 举报人ID
	Reason      string `json:"reason" form:"reason"`           // 举报原因
	AuditStatus int64  `json:"auditStatus" form:"auditStatus"` // 审核状态
	AuditTime   int64  `json:"auditTime" form:"auditTime"`     // 审核时间
	AuditUserId int64  `json:"auditUserId" form:"auditUserId"` // 审核人ID
	CreateTime  int64  `json:"createTime" form:"createTime"`   // 举报时间
}

// ForbiddenWord 违禁词
type ForbiddenWord struct {
	Model
	Type       string `gorm:"size:16" json:"type" form:"type"`       // 类型：word/regex
	Word       string `gorm:"size:128" json:"word" form:"word"`      // 违禁词
	Remark     string `gorm:"size:1024" json:"remark" form:"remark"` // 备注
	CreateTime int64  `json:"createTime" form:"createTime"`          // 举报时间
}
