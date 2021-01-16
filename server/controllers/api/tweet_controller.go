package api

import (
	"bbs-go/model"
	"bbs-go/model/constants"
	"github.com/mlogclub/simple/json"
	"github.com/sirupsen/logrus"
	"strconv"
	"strings"

	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple"

	"bbs-go/controllers/render"
	"bbs-go/services"
)

type TweetController struct {
	Ctx iris.Context
}

func (c *TweetController) GetFuck() *simple.JsonResult {
	go func() {
		tweets := services.TweetService.Find(simple.NewSqlCnd().Asc("id"))
		for _, tweet := range tweets {
			form := model.CreateTopicForm{
				Type:      constants.TopicTypeTweet,
				NodeId:    1,
				Content:   tweet.Content,
				ImageList: nil,
			}

			var images []string
			if err := json.Parse(tweet.ImageList, &images); err == nil {
				if len(images) > 0 {
					var imageList []model.ImageDTO
					for _, image := range images {
						imageList = append(imageList, model.ImageDTO{Url: image})
					}
					form.ImageList = imageList
				}
			}
			topic, err := services.TopicService.Publish(tweet.UserId, form)
			if err != nil {
				logrus.Error(err)
			} else {
				_ = services.TopicService.Updates(topic.Id, map[string]interface{}{
					"status":            tweet.Status,
					"create_time":       tweet.CreateTime,
					"last_comment_time": tweet.CreateTime,
					"comment_count":     tweet.CommentCount,
				})
			}
			simple.DB().Exec("update t_comment set entity_type = 'topic', entity_id = ? where entity_type = 'tweet' and entity_id = ?", topic.Id, tweet.Id)
			simple.DB().Exec("update t_user_like set entity_type = 'topic', entity_id = ? where entity_type = 'tweet' and entity_id = ?", topic.Id, tweet.Id)
			logrus.Info("tweet -> topic, tweetId=" + strconv.FormatInt(tweet.Id, 10) + ", topicId=" + strconv.FormatInt(topic.Id, 10))
		}
	}()
	return simple.JsonSuccess()
}

func (c *TweetController) PostCreate() *simple.JsonResult {
	user := services.UserTokenService.GetCurrent(c.Ctx)
	if err := services.UserService.CheckPostStatus(user); err != nil {
		return simple.JsonError(err)
	}
	content := strings.TrimSpace(simple.FormValue(c.Ctx, "content"))
	imageList := simple.FormValue(c.Ctx, "imageList")
	tweets, err := services.TweetService.Publish(user.Id, content, imageList)
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	return simple.JsonData(render.BuildTweet(tweets))
}

func (c *TweetController) GetList() *simple.JsonResult {
	cursor := simple.FormValueInt64Default(c.Ctx, "cursor", 0)
	tweets, cursor := services.TweetService.GetTweets(0, cursor)
	return simple.JsonCursorData(render.BuildTweets(tweets), strconv.FormatInt(cursor, 10))
}

func (c *TweetController) GetBy(tweetId int64) *simple.JsonResult {
	tweet := services.TweetService.Get(tweetId)
	return simple.JsonData(render.BuildTweet(tweet))
}

func (c *TweetController) PostLikeBy(tweetId int64) *simple.JsonResult {
	user := services.UserTokenService.GetCurrent(c.Ctx)
	if user == nil {
		return simple.JsonError(simple.ErrorNotLogin)
	}
	err := services.UserLikeService.TweetLike(user.Id, tweetId)
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	return simple.JsonSuccess()
}

func (c *TweetController) GetNewest() *simple.JsonResult {
	tweets := services.TweetService.GetNewest()
	return simple.JsonData(render.BuildTweets(tweets))
}

func (c *TweetController) GetUserTweets() *simple.JsonResult {
	userId, err := simple.FormValueInt64(c.Ctx, "userId")
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	cursor := simple.FormValueInt64Default(c.Ctx, "cursor", 0)
	tweets, cursor := services.TweetService.GetTweets(userId, cursor)
	return simple.JsonCursorData(render.BuildTweets(tweets), strconv.FormatInt(cursor, 10))
}
