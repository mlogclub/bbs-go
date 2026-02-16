package admin

import (
	"bbs-go/internal/models/constants"
	"bbs-go/internal/pkg/config"
	"os"
	"runtime"

	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple/web"
)

type CommonController struct {
	Ctx iris.Context
}

func (c *CommonController) GetSystem_info() *web.JsonResult {
	hostname, _ := os.Hostname()
	return web.NewEmptyRspBuilder().
		Put("os", runtime.GOOS).
		Put("arch", runtime.GOARCH).
		Put("numCpu", runtime.NumCPU()).
		Put("goVersion", runtime.Version()).
		Put("hostname", hostname).
		JsonResult()
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
