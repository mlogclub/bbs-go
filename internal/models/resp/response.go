package resp

import (
	"bbs-go/internal/models/constants"
	"time"

	"github.com/mlogclub/simple/web"
)

// UserInfo 用户简单信息
type UserInfo struct {
	Id           string           `json:"id"`
	Type         int              `json:"type"`
	Nickname     string           `json:"nickname"`
	Avatar       string           `json:"avatar"`
	SmallAvatar  string           `json:"smallAvatar"`
	Gender       constants.Gender `json:"gender"`
	Birthday     *time.Time       `json:"birthday"`
	TopicCount   int              `json:"topicCount"`   // 话题数量
	CommentCount int              `json:"commentCount"` // 跟帖数量
	FansCount    int              `json:"fansCount"`    // 粉丝数量
	FollowCount  int              `json:"followCount"`  // 关注数量
	Score        int              `json:"score"`        // 积分
	Exp          int              `json:"exp"`          // 经验值
	Level        int              `json:"level"`        // 等级
	LevelTitle   string           `json:"levelTitle"`   // 等级称号
	Description  string           `json:"description"`
	CreateTime   int64            `json:"createTime"`

	Forbidden bool `json:"forbidden"` // 是否禁言
	Followed  bool `json:"followed"`  // 是否关注

	// ExpProgress 经验值进度（当前等级内进度条数据），由 BuildUserInfo 根据 LevelConfig 计算填充；未登录或异常时为 nil
	ExpProgress *ExpProgressResponse `json:"expProgress,omitempty"`
}

// ExpProgressResponse 用户经验值进度（用于当前等级内的进度条展示）
// 计算依据：LevelConfig 中 NeedExp 表示达到该等级所需的累计经验，严格递增。
// 当前等级区间为 [当前级 NeedExp, 下一级 NeedExp)，进度 = 在此区间内已获得的经验占比。
type ExpProgressResponse struct {
	// CurrentExp 用户当前累计经验值（与 UserInfo.Exp 一致，便于组件只读进度）
	// 计算方式：直接取 user.Exp。
	CurrentExp int `json:"currentExp"`

	// Level 当前等级（与 UserInfo.Level 一致）
	// 计算方式：直接取 user.Level。
	Level int `json:"level"`

	// LevelTitle 当前等级称号（与 UserInfo.LevelTitle 一致）
	// 计算方式：由 LevelConfig(level).Title 得到。
	LevelTitle string `json:"levelTitle"`

	// ExpInCurrentLevel 当前等级内已获得的经验数（用于文案展示，如「120 / 350」中的 120）
	// 计算方式：当前累计经验 - 当前等级起始所需累计经验 = user.Exp - LevelConfig(level).NeedExp。
	// 若 user.Exp < 当前级 NeedExp，取 0；若已超过下一级 NeedExp，取 expNeedForNextLevel（封顶）。
	ExpInCurrentLevel int `json:"expInCurrentLevel"`

	// ExpNeedForNextLevel 从当前等级升到下一级，在本等级段内需要的经验数（即区间长度，用于文案中的「/ 350」）
	// 计算方式：下一级所需累计经验 - 当前级所需累计经验 = LevelConfig(level+1).NeedExp - LevelConfig(level).NeedExp。
	// 若已是最高等级（无下一级配置），则为 0，前端可配合 isMaxLevel 显示「已满级」或 100%。
	ExpNeedForNextLevel int `json:"expNeedForNextLevel"`

	// ExpProgressPercent 当前等级内经验进度百分比，取值 0～100，供进度条直接使用
	// 计算方式：round(ExpInCurrentLevel / ExpNeedForNextLevel * 100)。
	// 当 ExpNeedForNextLevel 为 0（满级）时取 100；若分母为 0 且未满级则取 0。
	ExpProgressPercent int `json:"expProgressPercent"`

	// IsMaxLevel 是否已为最高等级（无下一级可升）
	// 计算方式：不存在 LevelConfig(level+1) 或为配置中的最高级时为 true。
	// 为 true 时前端可显示 100% 或「已满级」。
	IsMaxLevel bool `json:"isMaxLevel"`
}

// UserDetail 用户详细信息
type UserDetail struct {
	UserInfo
	Username             string `json:"username"`
	BackgroundImage      string `json:"backgroundImage"`
	SmallBackgroundImage string `json:"smallBackgroundImage"`
	HomePage             string `json:"homePage"`
	Status               int    `json:"status"`
}

// UserProfile 用户个人信息
type UserProfile struct {
	UserDetail
	Roles         []string `json:"roles"`
	PasswordSet   bool     `json:"passwordSet"` // 密码已设置
	Email         string   `json:"email"`
	EmailVerified bool     `json:"emailVerified"`
}

type TagResponse struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

type ArticleSimpleResponse struct {
	Id           int64          `json:"id"`
	User         *UserInfo      `json:"user"`
	Tags         *[]TagResponse `json:"tags"`
	Title        string         `json:"title"`
	Summary      string         `json:"summary"`
	Cover        *ImageInfo     `json:"cover"`
	SourceUrl    string         `json:"sourceUrl"`
	ViewCount    int64          `json:"viewCount"`
	CommentCount int64          `json:"commentCount"`
	LikeCount    int64          `json:"likeCount"`
	CreateTime   int64          `json:"createTime"`
	Status       int            `json:"status"`
	Favorited    bool           `json:"favorited"`
}

type ArticleResponse struct {
	ArticleSimpleResponse
	Content string `json:"content"`
}

type NodeResponse struct {
	Id          int64  `json:"id"`
	Name        string `json:"name"`
	Logo        string `json:"logo"`
	Description string `json:"description"`
}

type SearchTopicResponse struct {
	Id         int64          `json:"id"`
	User       *UserInfo      `json:"user"`
	Node       *NodeResponse  `json:"node"`
	Tags       *[]TagResponse `json:"tags"`
	Title      string         `json:"title"`
	Summary    string         `json:"summary"`
	CreateTime int64          `json:"createTime"`
}

// 帖子列表返回实体
type TopicResponse struct {
	Id              string              `json:"id"`
	Type            constants.TopicType `json:"type"`
	User            *UserInfo           `json:"user"`
	Node            *NodeResponse       `json:"node"`
	Tags            *[]TagResponse      `json:"tags"`
	Title           string              `json:"title"`
	Summary         string              `json:"summary"`
	Content         string              `json:"content"`
	ImageList       []ImageInfo         `json:"imageList"`
	LastCommentTime int64               `json:"lastCommentTime"`
	ViewCount       int64               `json:"viewCount"`
	CommentCount    int64               `json:"commentCount"`
	LikeCount       int64               `json:"likeCount"`
	Liked           bool                `json:"liked"`
	CreateTime      int64               `json:"createTime"`
	Recommend       bool                `json:"recommend"`
	RecommendTime   int64               `json:"recommendTime"`
	Sticky          bool                `json:"sticky"`
	StickyTime      int64               `json:"stickyTime"`
	Status          int                 `json:"status"`
	Favorited       bool                `json:"favorited"`
	IpLocation      string              `json:"ipLocation"`
	Vote            *VoteResponse       `json:"vote"`
}

type VoteResponse struct {
	Id          int64                `json:"id"`
	Type        constants.VoteType   `json:"type"`
	Title       string               `json:"title"`
	ExpiredAt   int64                `json:"expiredAt"`
	VoteNum     int                  `json:"voteNum"`
	OptionCount int                  `json:"optionCount"`
	VoteCount   int                  `json:"voteCount"`
	Expired     bool                 `json:"expired"`
	Voted       bool                 `json:"voted"`
	OptionIds   []int64              `json:"optionIds"`
	Options     []VoteOptionResponse `json:"options"`
}

type VoteOptionResponse struct {
	Id        int64   `json:"id"`
	Content   string  `json:"content"`
	SortNo    int     `json:"sortNo"`
	VoteCount int     `json:"voteCount"`
	Percent   float64 `json:"percent"`
	Voted     bool    `json:"voted"`
}

// CommentResponse 评论返回数据
type CommentResponse struct {
	Id           int64                 `json:"id"`
	User         *UserInfo             `json:"user"`
	EntityType   string                `json:"entityType"`
	EntityId     int64                 `json:"entityId"`
	ContentType  constants.ContentType `json:"contentType"`
	Content      string                `json:"content"`
	ImageList    []ImageInfo           `json:"imageList"`
	LikeCount    int64                 `json:"likeCount"`
	CommentCount int64                 `json:"commentCount"`
	Liked        bool                  `json:"liked"`
	QuoteId      int64                 `json:"quoteId"`
	Quote        *CommentResponse      `json:"quote"`
	Replies      *web.CursorResult     `json:"replies"`
	IpLocation   string                `json:"ipLocation"`
	Status       int                   `json:"status"`
	CreateTime   int64                 `json:"createTime"`
}

// 收藏返回数据
type FavoriteResponse struct {
	Id         int64     `json:"id"`
	EntityType string    `json:"entityType"`
	EntityId   int64     `json:"entityId"`
	Deleted    bool      `json:"deleted"`
	Title      string    `json:"title"`
	Content    string    `json:"content"`
	User       *UserInfo `json:"user"`
	Url        string    `json:"url"`
	CreateTime int64     `json:"createTime"`
}

// 消息
type MessageResponse struct {
	Id           int64     `json:"id"`
	From         *UserInfo `json:"from"`    // 消息发送人
	UserId       int64     `json:"userId"`  // 消息接收人编号
	Title        string    `json:"title"`   // 标题
	Content      string    `json:"content"` // 消息内容
	QuoteContent string    `json:"quoteContent"`
	Type         int       `json:"type"`
	DetailUrl    string    `json:"detailUrl"` // 消息详情url
	ExtraData    string    `json:"extraData"`
	Status       int       `json:"status"`
	CreateTime   int64     `json:"createTime"`
}

// 图片
type ImageInfo struct {
	Url     string `json:"url"`
	Preview string `json:"preview"`
}

type TreeNode struct {
	Id       int64      `json:"id"`
	Key      int64      `json:"key"`
	Title    string     `json:"title"`
	Children []TreeNode `json:"children"`
}

type MenuResponse struct {
	Id         int64  `json:"id"`
	ParentId   *int64 `json:"parentId"`
	Type       string `json:"type"`
	Name       string `json:"name"`
	Title      string `json:"title"`
	Icon       string `json:"icon"`
	Path       string `json:"path"`
	Component  string `json:"component"`
	SortNo     int    `json:"sortNo"`
	Status     int    `json:"status"`
	CreateTime int64  `json:"createTime"`
	UpdateTime int64  `json:"updateTime"`
}

type MenuTreeResponse struct {
	MenuResponse
	Level    int                `json:"level"`
	Children []MenuTreeResponse `json:"children"`
}

type DictResponse struct {
	Id         int64  `json:"id"`
	TypeId     int64  `json:"typeId"`
	ParentId   *int64 `json:"parentId"`   // 上级分类
	Name       string `json:"name"`       // 名称
	Label      string `json:"label"`      // 标题
	Value      string `json:"value"`      // 值
	SortNo     int    `json:"sortNo"`     // 排序
	Status     int    `json:"status"`     // 状态
	CreateTime int64  `json:"createTime"` // 创建时间
	UpdateTime int64  `json:"updateTime"` // 更新时间
}

type DictListResponse struct {
	DictResponse
	Children []DictListResponse `json:"children"`
}

// TaskGroupInfo 任务分组信息（含多语言名称）
type TaskGroupInfo struct {
	Key  constants.TaskGroup `json:"key"`
	Name string              `json:"name"`
}

type TaskResponse struct {
	Id             int64                 `json:"id"`
	GroupName      constants.TaskGroup   `json:"groupName"`
	Title          string                `json:"title"`
	Description    string                `json:"description"`
	EventType      string                `json:"eventType"`
	Period         constants.TaskPeriod  `json:"period"`
	EventCount     int                   `json:"eventCount"`
	MaxFinishCount int                   `json:"maxFinishCount"`
	Score          int                   `json:"score"`
	Exp            int                   `json:"exp"`
	BadgeId        int64                 `json:"badgeId"`
	BtnName        string                `json:"btnName"`
	ActionUrl      string                `json:"actionUrl"`
	SortNo         int                   `json:"sortNo"`
	StartTime      int64                 `json:"startTime"`
	EndTime        int64                 `json:"endTime"`
	Status         int                   `json:"status"`
	UserProgress   *TaskProgressResponse `json:"userProgress,omitempty"`
}

// TaskProgressResponse 用户在某任务上的当前进度
type TaskProgressResponse struct {
	PeriodKey      int `json:"periodKey"`      // 当前周期 key（一次性为 0）
	EventProgress  int `json:"eventProgress"`  // 本周期已累计的事件次数
	EventTarget    int `json:"eventTarget"`    // 完成一次任务需要的事件次数
	FinishedCount  int `json:"finishedCount"`  // 本周期已完成次数
	MaxFinishCount int `json:"maxFinishCount"` // 本周期最多可完成次数
}

type BadgeResponse struct {
	Id          int64  `json:"id"`
	Name        string `json:"name"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
	SortNo      int    `json:"sortNo"`
	Status      int    `json:"status"`
	Owned       bool   `json:"owned"`      // 当前登录用户是否已获得
	Worn        bool   `json:"worn"`       // 是否已佩戴
	ObtainTime  int64  `json:"obtainTime"` // 获得时间（未获得为0）
}
