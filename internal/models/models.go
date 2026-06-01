package models

import (
	"bbs-go/internal/models/constants"
	"database/sql"
	"time"
)

var Models = []interface{}{
	&Migration{},
	&UserRole{}, &Role{}, &Permission{}, &RolePermission{}, &DictType{}, &Dict{},

	&User{}, &UserToken{}, &ThirdUser{}, &Tag{}, &Article{}, &ArticleTag{}, &Comment{}, &Favorite{}, &Topic{}, &Category{},
	&TopicTag{}, &UserLike{}, &Message{}, &SysConfig{}, &Link{},
	&TaskConfig{}, &UserTaskEvent{}, &UserTaskLog{},
	&Badge{}, &UserBadge{},
	&LevelConfig{},
	&Vote{}, &VoteOption{}, &VoteRecord{},
	&UserScoreLog{}, &UserExpLog{},
	&OperateLog{}, &EmailLog{}, &EmailCode{}, &SmsCode{}, &CheckIn{}, &UserFollow{}, &UserFeed{}, &UserReport{},
	&ForbiddenWord{},
	&Attachment{}, &AttachmentDownloadLog{},

	&TopicStake{}, &UserHeatLog{}, &UserHeatStats{}, &HeatPublicPool{},
	&SystemMintLog{}, &TopicInteractionSnapshot{}, &HeatCirculationSnapshot{},
	&DailyFlameOffset{}, &SettlementTaskLog{}, &ColdStorageLog{},
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

type UserRole struct {
	Model
	UserId     int64 `gorm:"uniqueIndex:idx_user_role" json:"userId" form:"userId"`
	RoleId     int64 `gorm:"uniqueIndex:idx_user_role" json:"roleId" form:"roleId"`
	CreateTime int64 `gorm:"not null;default:0" json:"createTime" form:"createTime"` // 创建时间
}

type Permission struct {
	Model
	Type        string `gorm:"size:32;not null;index:idx_permission_type" json:"type" form:"type"`                 // 权限类型
	Code        string `gorm:"size:128;not null;uniqueIndex:uk_permission_code" json:"code" form:"code"`           // 权限编码
	Name        string `gorm:"size:64;not null" json:"name" form:"name"`                                           // 权限名称
	GroupName   string `gorm:"size:64;not null;index:idx_permission_group_name" json:"groupName" form:"groupName"` // 权限分组
	Description string `gorm:"size:256" json:"description" form:"description"`                                     // 描述
	SortNo      int    `gorm:"not null;default:0" json:"sortNo" form:"sortNo"`                                     // 排序
	Status      int    `gorm:"not null;default:0;index:idx_permission_status" json:"status" form:"status"`         // 状态
	CreateTime  int64  `gorm:"not null;default:0" json:"createTime" form:"createTime"`                             // 创建时间
	UpdateTime  int64  `gorm:"not null;default:0" json:"updateTime" form:"updateTime"`                             // 更新时间
}

type RolePermission struct {
	Model
	RoleId       int64 `gorm:"not null;uniqueIndex:uk_role_permission_role_permission;index:idx_role_permission_role_id" json:"roleId" form:"roleId"`
	PermissionId int64 `gorm:"not null;uniqueIndex:uk_role_permission_role_permission;index:idx_role_permission_permission_id" json:"permissionId" form:"permissionId"`
	CreateTime   int64 `gorm:"not null;default:0" json:"createTime" form:"createTime"` // 创建时间
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
	HeatPoints       int              `gorm:"type:int(11);not null;default:0" json:"heatPoints" form:"heatPoints"`     // 热度点
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

// Category 话题节点（支持一级 parent_id=0 / 二级 parent_id>0）
type Category struct {
	Model
	ParentId    int64                  `gorm:"not null;default:0;index:idx_category_parent_id" json:"parentId" form:"parentId"` // 父节点ID，0=一级
	Name        string                 `gorm:"size:32;unique" json:"name" form:"name"`                                          // 名称（一级全局唯一；二级同父下唯一，应用层校验）
	Type        constants.CategoryType `gorm:"size:16;not null;default:normal;index:idx_category_type" json:"type" form:"type"` // 节点类型：normal/qa
	Description string                 `gorm:"size:1024" json:"description" form:"description"`                                 // 描述
	Logo        string                 `gorm:"size:1024" json:"logo" form:"logo"`                                               // 图标
	SortNo      int                    `gorm:"type:int(11);index:idx_category_sort_no" json:"sortNo" form:"sortNo"`             // 排序编号
	Status      int                    `gorm:"type:int(11);not null" json:"status" form:"status"`                               // 状态
	CreateTime  int64                  `json:"createTime" form:"createTime"`                                                    // 创建时间
}

// 话题节点
type Topic struct {
	Model
	Type              constants.TopicType   `gorm:"type:int(11);not null:default:0;index:idx_topic_type_category_id,priority:1;index:idx_topic_type_qa_status,priority:1" json:"type" form:"type"` // 类型
	CategoryId        int64                 `gorm:"not null;index:idx_category_id;index:idx_topic_type_category_id,priority:2" json:"categoryId" form:"categoryId"`                                // 节点编号
	QaStatus          constants.QaStatus    `gorm:"size:16;not null;default:unsolved;index:idx_topic_type_qa_status,priority:2" json:"qaStatus" form:"qaStatus"`                                   // 问答状态
	AcceptedCommentId int64                 `gorm:"not null;default:0;index:idx_topic_accepted_comment_id" json:"acceptedCommentId" form:"acceptedCommentId"`                                      // 采纳评论ID
	SolvedAt          int64                 `gorm:"not null;default:0" json:"solvedAt" form:"solvedAt"`                                                                                            // 解决时间
	BountyScore       int                   `gorm:"type:int(11);not null;default:0" json:"bountyScore" form:"bountyScore"`                                                                         // 悬赏积分（仅问答帖有效，0 表示无悬赏）
	UserId            int64                 `gorm:"not null;index:idx_topic_user_id;" json:"userId" form:"userId"`                                                                                 // 用户
	Title             string                `gorm:"size:128" json:"title" form:"title"`                                                                                                            // 标题
	ContentType       constants.ContentType `gorm:"size:32;default:markdown" json:"contentType" form:"contentType"`                                                                                // 内容类型（html/markdown）
	Content           string                `gorm:"type:longtext" json:"content" form:"content"`                                                                                                   // 内容
	ImageList         string                `gorm:"type:longtext" json:"imageList" form:"imageList"`                                                                                               // 图片
	HideContent       string                `gorm:"type:longtext" json:"hideContent" form:"hideContent"`                                                                                           // 回复可见内容
	VoteId            int64                 `gorm:"not null;default:0" json:"voteId" form:"voteId"`                                                                                                // 投票ID
	Recommend         bool                  `gorm:"not null;index:idx_recommend" json:"recommend" form:"recommend"`                                                                                // 是否推荐
	RecommendTime     int64                 `gorm:"not null" json:"recommendTime" form:"recommendTime"`                                                                                            // 推荐时间
	Sticky            bool                  `gorm:"not null;index:idx_sticky_sticky_time" json:"sticky" form:"sticky"`                                                                             // 置顶
	StickyTime        int64                 `gorm:"not null;index:idx_sticky_sticky_time" json:"stickyTime" form:"stickyTime"`                                                                     // 置顶时间
	ViewCount         int64                 `gorm:"not null" json:"viewCount" form:"viewCount"`                                                                                                    // 查看数量
	CommentCount      int64                 `gorm:"not null" json:"commentCount" form:"commentCount"`                                                                                              // 跟帖数量
	LikeCount         int64                 `gorm:"not null" json:"likeCount" form:"likeCount"`                                                                                                    // 点赞数量
	Status            int                   `gorm:"type:int(11);index:idx_topic_status;" json:"status" form:"status"`                                                                              // 状态：0：正常、1：删除
	LastCommentTime   int64                 `gorm:"index:idx_topic_last_comment_time" json:"lastCommentTime" form:"lastCommentTime"`                                                               // 最后回复时间
	LastCommentUserId int64                 `json:"lastCommentUserId" form:"lastCommentUserId"`                                                                                                    // 最后回复用户
	UserAgent         string                `gorm:"size:1024" json:"userAgent" form:"userAgent"`                                                                                                   // UserAgent
	Ip                string                `gorm:"size:128" json:"ip" form:"ip"`                                                                                                                  // IP
	IpLocation        string                `gorm:"size:64" json:"ipLocation" form:"ipLocation"`                                                                                                   // IP 属地
	CreateTime        int64                 `gorm:"index:idx_topic_create_time" json:"createTime" form:"createTime"`                                                                               // 创建时间
	ExtraData         string                `gorm:"type:text" json:"extraData" form:"extraData"`                                                                                                   // 扩展数据
	EverViral         bool                  `gorm:"not null;default:false" json:"everViral" form:"everViral"`                                                                                      // 曾经达到过热门（不可逆）
	FlameLockedLevel  int                   `gorm:"default:0" json:"flameLockedLevel" form:"flameLockedLevel"`                                                                                     // 0=算法控制，>0=管理员锁定火焰等级
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
	Type         int    `gorm:"type:int(11);not null" json:"type" form:"type"`                   // 消息类型：评论/点赞/收藏/推荐/删除/文章评论/等级提升/获得勋章
	ExtraData    string `gorm:"type:text" json:"extraData" form:"extraData"`                     // 扩展数据
	Status       int    `gorm:"type:int(11);not null" json:"status" form:"status"`               // 状态：0：未读、1：已读
	CreateTime   int64  `json:"createTime" form:"createTime"`                                    // 创建时间
}

// 系统配置
type SysConfig struct {
	Model
	Key          string `gorm:"not null;size:128;unique" json:"key" form:"key"`              // 配置key
	Value        string `gorm:"type:text" json:"value" form:"value"`                         // 当前生效的配置值
	PendingValue string `gorm:"type:text" json:"pendingValue" form:"pendingValue"`           // 待生效值（下次结算时覆盖 Value，然后清空）
	Name         string `gorm:"not null;size:128" json:"name" form:"name"`                   // 配置名称
	Description  string `gorm:"size:1024" json:"description" form:"description"`             // 配置描述
	CreateTime   int64  `gorm:"not null" json:"createTime" form:"createTime"`                // 创建时间
	UpdateTime   int64  `gorm:"not null" json:"updateTime" form:"updateTime"`                // 更新时间
}

// 友链
type Link struct {
	Model
	Url        string `gorm:"not null;type:text" json:"url" form:"url"`                        // 链接
	Title      string `gorm:"not null;size:128" json:"title" form:"title"`                     // 标题
	Summary    string `gorm:"size:1024" json:"summary" form:"summary"`                         // 站点描述
	SortNo     int    `gorm:"type:int(11);index:idx_link_sort_no" json:"sortNo" form:"sortNo"` // 排序编号
	Status     int    `gorm:"type:int(11);not null" json:"status" form:"status"`               // 状态
	CreateTime int64  `gorm:"not null" json:"createTime" form:"createTime"`                    // 创建时间
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
	BizType    string `gorm:"not null;default:'';size:32;index:idx_email_code_biz_type" json:"bizType" form:"bizType"`
	Email      string `gorm:"not null;size:128" json:"email" form:"email"`       // 邮箱
	Code       string `gorm:"not null;size:8" json:"code" form:"code"`           // 验证码
	Token      string `gorm:"not null;size:32;unique" json:"token" form:"token"` // 验证码token
	Title      string `gorm:"size:1024" json:"title" form:"title"`               // 标题
	Content    string `gorm:"type:text" json:"content" form:"content"`           // 内容
	Used       bool   `gorm:"not null" json:"used" form:"used"`                  // 是否使用
	CreateTime int64  `json:"createTime" form:"createTime"`                      // 创建时间
}

// 邮件发送记录
type EmailLog struct {
	Model
	ToEmail    string `gorm:"not null;size:128;index:idx_email_log_to_email" json:"toEmail" form:"toEmail"`
	Subject    string `gorm:"size:1024" json:"subject" form:"subject"`
	Content    string `gorm:"type:longtext" json:"content" form:"content"`
	BizType    string `gorm:"not null;default:'';size:32;index:idx_email_log_biz_type" json:"bizType" form:"bizType"`
	Status     int    `gorm:"not null;index:idx_email_log_status" json:"status" form:"status"`
	ErrorMsg   string `gorm:"type:text" json:"errorMsg" form:"errorMsg"`
	CreateTime int64  `gorm:"index:idx_email_log_create_time" json:"createTime" form:"createTime"`
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

// Attachment 帖子附件
type Attachment struct {
	Id            string `gorm:"primaryKey;size:64" json:"id" form:"id"`
	TopicId       int64  `gorm:"not null;default:0;index:idx_attachment_topic_id" json:"topicId" form:"topicId"`          // 所属帖子 ID
	UserId        int64  `gorm:"not null;index:idx_attachment_user_id" json:"userId" form:"userId"`                       // 上传者（发帖人）ID
	FileName      string `gorm:"size:256" json:"fileName" form:"fileName"`                                                // 原始文件名
	FileUrl       string `gorm:"size:1024" json:"fileUrl" form:"fileUrl"`                                                 // 访问地址（相对路径或完整 URL，上传返回）
	FileSize      int64  `gorm:"not null;default:0" json:"fileSize" form:"fileSize"`                                      // 文件大小（字节）
	FileType      string `gorm:"size:64" json:"fileType" form:"fileType"`                                                 // MIME 或扩展名
	DownloadScore int    `gorm:"type:int(11);not null;default:0" json:"downloadScore" form:"downloadScore"`               // 下载所需积分，0 表示免费
	DownloadCount int    `gorm:"type:int(11);not null;default:0" json:"downloadCount" form:"downloadCount"`               // 下载次数
	Status        int    `gorm:"type:int(11);not null;index:idx_attachment_status" json:"status" form:"status"`           // 状态：正常/删除
	CreateTime    int64  `gorm:"not null;default:0;index:idx_attachment_create_time" json:"createTime" form:"createTime"` // 创建时间（毫秒）
	UpdateTime    int64  `gorm:"not null;default:0" json:"updateTime" form:"updateTime"`                                  // 更新时间（毫秒）
}

// AttachmentDownloadLog 附件下载记录（用户已购买，同一附件再次下载不扣积分）
type AttachmentDownloadLog struct {
	Model
	UserId       int64  `gorm:"not null;uniqueIndex:uk_attachment_download_log_ua" json:"userId" form:"userId"`                     // 下载用户 ID
	AttachmentId string `gorm:"not null;size:64;uniqueIndex:uk_attachment_download_log_ua" json:"attachmentId" form:"attachmentId"` // 附件 ID（UUID）
	CreateTime   int64  `gorm:"not null;default:0" json:"createTime" form:"createTime"`                                             // 首次下载（支付）时间（毫秒）
}

// TopicStake 帖子质押记录（分区表）
// 分区策略：按 update_time/100000000 范围分区（每半年一个分区）
// 注意：update_time 必须在主键中以满足分区要求
type TopicStake struct {
	Model
	TopicId        int64  `gorm:"not null;index:idx_stake_topic" json:"topicId" form:"topicId"`
	UserId         int64  `gorm:"not null;index:idx_stake_user" json:"userId" form:"userId"`
	HeatPoints     int    `gorm:"not null" json:"heatPoints" form:"heatPoints"`         // 当前押注的热度点（含滚入利息）
	OriginalPoints int    `gorm:"not null" json:"originalPoints" form:"originalPoints"` // 最初押注的热度点
	StakeDay       string `gorm:"size:8;not null;index" json:"stakeDay" form:"stakeDay"` // 首次质押日 yyyyMMdd
	Status         int    `gorm:"not null;index" json:"status" form:"status"`           // 质押中/已赎回/已结算失败
	LastSettleDay  string `gorm:"size:8" json:"lastSettleDay" form:"lastSettleDay"`     // 上次结算日
	CreateTime     int64  `json:"createTime" form:"createTime"`
	UpdateTime     int64  `gorm:"not null" json:"updateTime" form:"updateTime"` // 分区字段（必须在主键中）
}

// UserHeatLog 用户热度点流水
type UserHeatLog struct {
	Model
	UserId     int64  `gorm:"not null;index" json:"userId" form:"userId"`
	ChangeType string `gorm:"size:32;not null" json:"changeType" form:"changeType"` // DailyCheckIn / StakeOut / SettleProfit / SettleLoss / Redeem / Decay
	Amount     int    `gorm:"not null" json:"amount" form:"amount"`                 // 变动量（正或负）
	Balance    int    `gorm:"not null" json:"balance" form:"balance"`               // 变动后余额
	RefId      string `gorm:"size:64" json:"refId" form:"refId"`                    // 关联 ID
	Remark     string `gorm:"size:256" json:"remark" form:"remark"`
	CreateTime int64  `json:"createTime" form:"createTime"`
}

// UserHeatStats 用户热度点统计（利用率计算基础）
type UserHeatStats struct {
	UserId             int64 `gorm:"primaryKey" json:"userId" form:"userId"`
	TotalPoints        int   `gorm:"not null" json:"totalPoints" form:"totalPoints"`             // 持有总量（余额 + 活跃质押）
	StakedInWindow     int   `gorm:"not null" json:"stakedInWindow" form:"stakedInWindow"`       // 近 7 天累计质押量
	LastStakeTime      int64 `gorm:"default:0" json:"lastStakeTime" form:"lastStakeTime"`
	DecayedAccumulated int   `gorm:"default:0" json:"decayedAccumulated" form:"decayedAccumulated"`
	CooldownPoints     int   `gorm:"not null;default:0" json:"cooldownPoints" form:"cooldownPoints"` // 冷却中的热度点数
	CooldownUntil      int64 `gorm:"default:0" json:"cooldownUntil" form:"cooldownUntil"`           // 冷却期到期时间（毫秒时间戳）
	UpdateTime         int64 `json:"updateTime" form:"updateTime"`
}

// HeatPublicPool 公共奖池记录
type HeatPublicPool struct {
	Model
	Source       string `gorm:"size:64;not null" json:"source" form:"source"` // DecayInflow / StakeLoss
	Amount       int    `gorm:"not null" json:"amount" form:"amount"`
	RefId        string `gorm:"size:64" json:"refId" form:"refId"`
	BalanceAfter int    `gorm:"not null" json:"balanceAfter" form:"balanceAfter"` // 变动后奖池余额
	Remark       string `gorm:"size:256" json:"remark" form:"remark"`
	CreateTime   int64  `json:"createTime" form:"createTime"`
}

// SystemMintLog 系统铸币日志（仅创世和补贴期使用）
type SystemMintLog struct {
	Model
	MintType    string `gorm:"size:32;not null" json:"mintType" form:"mintType"` // GenesisAirdrop / SubsidyMint
	Amount      int    `gorm:"not null" json:"amount" form:"amount"`
	RecipientId int64  `gorm:"default:0" json:"recipientId" form:"recipientId"` // 0 表示进入公共奖池
	Remark      string `gorm:"size:256" json:"remark" form:"remark"`
	CreateTime  int64  `json:"createTime" form:"createTime"`
}

// TopicInteractionSnapshot 每日互动快照
type TopicInteractionSnapshot struct {
	Model
	TopicId          int64  `gorm:"not null;uniqueIndex:uk_snapshot_topic_date" json:"topicId" form:"topicId"`
	SnapshotDate     string `gorm:"size:8;not null;uniqueIndex:uk_snapshot_topic_date" json:"snapshotDate" form:"snapshotDate"`
	UniqueLikers     int    `gorm:"not null;default:0" json:"uniqueLikers" form:"uniqueLikers"`
	UniqueCommenters int    `gorm:"not null;default:0" json:"uniqueCommenters" form:"uniqueCommenters"`
	ValidComments    int    `gorm:"not null;default:0" json:"validComments" form:"validComments"`
	TotalLikes       int    `gorm:"not null;default:0" json:"totalLikes" form:"totalLikes"`
	TotalComments    int    `gorm:"not null;default:0" json:"totalComments" form:"totalComments"`
	StakeTotal       int    `gorm:"not null;default:0" json:"stakeTotal" form:"stakeTotal"` // 快照时的质押总量
	CreateTime       int64  `json:"createTime" form:"createTime"`
}

// HeatCirculationSnapshot 活跃流通快照
type HeatCirculationSnapshot struct {
	SnapshotDate      string `gorm:"primaryKey;size:8" json:"snapshotDate" form:"snapshotDate"`
	TotalSupply       int    `gorm:"not null" json:"totalSupply" form:"totalSupply"`
	ActiveCirculation int    `gorm:"not null" json:"activeCirculation" form:"activeCirculation"`
	StakedTotal       int    `gorm:"not null" json:"stakedTotal" form:"stakedTotal"`
	ActiveUserCount   int    `gorm:"not null" json:"activeUserCount" form:"activeUserCount"`
	DailyDecayTruncated int  `gorm:"not null;default:0" json:"dailyDecayTruncated" form:"dailyDecayTruncated"` // 当日截断的衰减量
	CreateTime        int64  `json:"createTime" form:"createTime"`
}

// DailyFlameOffset 每日火焰等级偏移量
type DailyFlameOffset struct {
	Date           string  `gorm:"primaryKey;size:8" json:"date" form:"date"`
	Phase12Offset  float64 `gorm:"not null" json:"phase12Offset" form:"phase12Offset"`
	Phase23Offset  float64 `gorm:"not null" json:"phase23Offset" form:"phase23Offset"`
	Flame2Offset   float64 `gorm:"not null" json:"flame2Offset" form:"flame2Offset"`
	Flame3Offset   float64 `gorm:"not null" json:"flame3Offset" form:"flame3Offset"`
	Flame4Offset   float64 `gorm:"not null" json:"flame4Offset" form:"flame4Offset"`
	Flame5Offset   float64 `gorm:"not null" json:"flame5Offset" form:"flame5Offset"`
	CreateTime     int64   `json:"createTime" form:"createTime"`
}

// SettlementTaskLog 结算任务日志
type SettlementTaskLog struct {
	Model
	TaskDate       string `gorm:"size:8;not null;uniqueIndex:uk_task_date_type" json:"taskDate" form:"taskDate"`
	TaskType       string `gorm:"size:32;not null;uniqueIndex:uk_task_date_type" json:"taskType" form:"taskType"`
	Status         int    `gorm:"not null;default:0" json:"status" form:"status"` // 0=执行中，1=已完成，2=失败
	StartedAt      int64  `gorm:"not null" json:"startedAt" form:"startedAt"`
	FinishedAt     int64  `gorm:"default:0" json:"finishedAt" form:"finishedAt"`
	BatchCount     int    `gorm:"default:0" json:"batchCount" form:"batchCount"`
	TotalProcessed int    `gorm:"default:0" json:"totalProcessed" form:"totalProcessed"`
	ErrorMsg       string `gorm:"type:text" json:"errorMsg" form:"errorMsg"`
}

// ColdStorageLog 冷存储归档执行日志
type ColdStorageLog struct {
	Model
	TableName    string `gorm:"size:64;not null" json:"tableName"`
	ArchiveTable string `gorm:"size:64;not null" json:"archiveTable"`
	TriggeredBy  string `gorm:"size:16;not null" json:"triggeredBy"`  // time / row_count / both
	RowsBefore   int64  `gorm:"not null" json:"rowsBefore"`
	RowsMigrated int64  `gorm:"not null" json:"rowsMigrated"`
	RowsAfter    int64  `gorm:"not null" json:"rowsAfter"`
	DurationMs   int64  `gorm:"not null" json:"durationMs"`
	Status       int    `gorm:"not null;default:0" json:"status"`     // 0=成功，1=部分失败，2=失败
	ErrorMsg     string `gorm:"type:text" json:"errorMsg"`
	CreateTime   int64  `gorm:"not null" json:"createTime"`
}
