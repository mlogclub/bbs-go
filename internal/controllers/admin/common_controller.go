package admin

import (
	"bbs-go/internal/models"
	"bbs-go/internal/models/constants"
	"bbs-go/internal/pkg/config"
	"bbs-go/internal/repositories"
	"time"

	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple/sqls"
	"github.com/mlogclub/simple/web"
)

type CommonController struct {
	Ctx iris.Context
}

type dashboardRecentItem struct {
	Id         int64  `json:"id"`
	Title      string `json:"title,omitempty"`
	Content    string `json:"content,omitempty"`
	Nickname   string `json:"nickname,omitempty"`
	CreateTime int64  `json:"createTime"`
}

func (c *CommonController) GetOverview() *web.JsonResult {
	now := time.Now()
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).UnixMilli()
	db := sqls.DB()

	metrics := map[string]int64{
		"totalUsers":    repositories.UserRepository.Count(db, sqls.NewCnd().Eq("status", constants.StatusOk)),
		"totalTopics":   repositories.TopicRepository.Count(db, sqls.NewCnd().Eq("status", constants.StatusOk)),
		"totalArticles": repositories.ArticleRepository.Count(db, sqls.NewCnd().Eq("status", constants.StatusOk)),
		"todayUsers":    repositories.UserRepository.Count(db, sqls.NewCnd().Eq("status", constants.StatusOk).Gte("create_time", todayStart)),
		"todayTopics":   repositories.TopicRepository.Count(db, sqls.NewCnd().Eq("status", constants.StatusOk).Gte("create_time", todayStart)),
	}

	pending := map[string]int64{
		"pendingTopics":   repositories.TopicRepository.Count(db, sqls.NewCnd().Eq("status", constants.StatusReview)),
		"pendingArticles": repositories.ArticleRepository.Count(db, sqls.NewCnd().Eq("status", constants.StatusReview)),
		"pendingReports":  repositories.UserReportRepository.Count(db, sqls.NewCnd().Eq("audit_status", 0)),
		"failedEmails":    repositories.EmailLogRepository.Count(db, sqls.NewCnd().Eq("status", constants.EmailLogStatusFailed)),
	}

	recentTopics := repositories.TopicRepository.Find(db, sqls.NewCnd().Eq("status", constants.StatusOk).Desc("id").Limit(5))
	recentUsers := repositories.UserRepository.Find(db, sqls.NewCnd().Eq("status", constants.StatusOk).Desc("id").Limit(5))

	return web.NewEmptyRspBuilder().
		Put("metrics", metrics).
		Put("pending", pending).
		Put("recent", map[string]interface{}{
			"topics": buildRecentTopicItems(recentTopics),
			"users":  buildRecentUserItems(recentUsers),
		}).
		JsonResult()
}

func buildRecentTopicItems(topics []models.Topic) []dashboardRecentItem {
	items := make([]dashboardRecentItem, 0, len(topics))
	for _, topic := range topics {
		items = append(items, dashboardRecentItem{
			Id:         topic.Id,
			Title:      topic.Title,
			CreateTime: topic.CreateTime,
		})
	}
	return items
}

func buildRecentUserItems(users []models.User) []dashboardRecentItem {
	items := make([]dashboardRecentItem, 0, len(users))
	for _, user := range users {
		items = append(items, dashboardRecentItem{
			Id:         user.Id,
			Nickname:   user.Nickname,
			CreateTime: user.CreateTime,
		})
	}
	return items
}

type TaskEventTypeItem struct {
	Value string `json:"value"`
	Title string `json:"title"`
}

// GetTask_event_types 获取任务事件类型枚举（用于后台下拉选择）
func (c *CommonController) GetTask_event_types() *web.JsonResult {
	lang := config.Instance.Language
	if !lang.IsValid() {
		lang = config.DefaultLanguage
	}

	items := []TaskEventTypeItem{
		{Value: constants.TaskEventTypeUserLogin},
		{Value: constants.TaskEventTypeCheckIn},
		{Value: constants.TaskEventTypeTopicCreate},
		{Value: constants.TaskEventTypeQaQuestion},
		{Value: constants.TaskEventTypeQaAnswerAccept},
		{Value: constants.TaskEventTypeCommentCreate},
		{Value: constants.TaskEventTypeFollowCreate},
		{Value: constants.TaskEventTypeFavoriteCreate},
		{Value: constants.TaskEventTypeLikeCreate},
		{Value: constants.TaskEventTypeLevel10},
	}

	if lang == config.LanguageEnUS {
		titleMap := map[string]string{
			constants.TaskEventTypeUserLogin:      "Daily login",
			constants.TaskEventTypeCheckIn:        "Check-in",
			constants.TaskEventTypeTopicCreate:    "Create topic",
			constants.TaskEventTypeQaQuestion:     "Publish question",
			constants.TaskEventTypeQaAnswerAccept: "Answer accepted",
			constants.TaskEventTypeCommentCreate:  "Create comment",
			constants.TaskEventTypeFollowCreate:   "Follow user",
			constants.TaskEventTypeFavoriteCreate: "Favorite",
			constants.TaskEventTypeLikeCreate:     "Like",
			constants.TaskEventTypeLevel10:        "Reach level 10",
		}
		for i := range items {
			items[i].Title = titleMap[items[i].Value]
		}
	} else {
		titleMap := map[string]string{
			constants.TaskEventTypeUserLogin:      "每日登录",
			constants.TaskEventTypeCheckIn:        "签到",
			constants.TaskEventTypeTopicCreate:    "发帖",
			constants.TaskEventTypeQaQuestion:     "发布问题",
			constants.TaskEventTypeQaAnswerAccept: "回答被采纳",
			constants.TaskEventTypeCommentCreate:  "评论",
			constants.TaskEventTypeFollowCreate:   "关注用户",
			constants.TaskEventTypeFavoriteCreate: "收藏",
			constants.TaskEventTypeLikeCreate:     "点赞",
			constants.TaskEventTypeLevel10:        "达到等级 10",
		}
		for i := range items {
			items[i].Title = titleMap[items[i].Value]
		}
	}

	return web.JsonData(items)
}
