package services

import (
	"bbs-go/model/constants"
	"github.com/mlogclub/simple"

	"bbs-go/model"
	"bbs-go/repositories"
)

var TweetService = newTweetService()

func newTweetService() *tweetService {
	return &tweetService{}
}

type tweetService struct {
}

func (s *tweetService) Get(id int64) *model.Tweet {
	return repositories.TweetRepository.Get(simple.DB(), id)
}

func (s *tweetService) Take(where ...interface{}) *model.Tweet {
	return repositories.TweetRepository.Take(simple.DB(), where...)
}

func (s *tweetService) Find(cnd *simple.SqlCnd) []model.Tweet {
	return repositories.TweetRepository.Find(simple.DB(), cnd)
}

func (s *tweetService) FindOne(cnd *simple.SqlCnd) *model.Tweet {
	return repositories.TweetRepository.FindOne(simple.DB(), cnd)
}

func (s *tweetService) FindPageByParams(params *simple.QueryParams) (list []model.Tweet, paging *simple.Paging) {
	return repositories.TweetRepository.FindPageByParams(simple.DB(), params)
}

func (s *tweetService) FindPageByCnd(cnd *simple.SqlCnd) (list []model.Tweet, paging *simple.Paging) {
	return repositories.TweetRepository.FindPageByCnd(simple.DB(), cnd)
}

func (s *tweetService) Count(cnd *simple.SqlCnd) int64 {
	return repositories.TweetRepository.Count(simple.DB(), cnd)
}

func (s *tweetService) GetTweets(cursor int64) (tweets []model.Tweet, nextCursor int64) {
	cnd := simple.NewSqlCnd().Eq("status", constants.StatusOk).Desc("id").Limit(20)
	if cursor > 0 {
		cnd.Lt("id", cursor)
	}
	tweets = repositories.TweetRepository.Find(simple.DB(), cnd)
	if len(tweets) > 0 {
		nextCursor = tweets[len(tweets)-1].Id
	} else {
		nextCursor = cursor
	}
	return
}

func (s *tweetService) GetNewest() (tweets []model.Tweet) {
	cnd := simple.NewSqlCnd().Eq("status", constants.StatusOk).Desc("id").Limit(6)
	tweets = repositories.TweetRepository.Find(simple.DB(), cnd)
	return
}

func (s *tweetService) Publish(userId int64, content, imageList string) (*model.Tweet, error) {
	tweet := &model.Tweet{
		UserId:     userId,
		Content:    content,
		ImageList:  imageList,
		Status:     constants.StatusOk,
		CreateTime: simple.NowTimestamp(),
	}
	if err := repositories.TweetRepository.Create(simple.DB(), tweet); err != nil {
		return nil, err
	}
	return tweet, nil
}

func (s *tweetService) OnComment(tweetId int64) {
	simple.DB().Exec("update t_tweet set comment_count = comment_count + 1 where id = ?", tweetId)
}

func (s *tweetService) Create(t *model.Tweet) error {
	return repositories.TweetRepository.Create(simple.DB(), t)
}

func (s *tweetService) Update(t *model.Tweet) error {
	return repositories.TweetRepository.Update(simple.DB(), t)
}

func (s *tweetService) Updates(id int64, columns map[string]interface{}) error {
	return repositories.TweetRepository.Updates(simple.DB(), id, columns)
}

func (s *tweetService) UpdateColumn(id int64, name string, value interface{}) error {
	return repositories.TweetRepository.UpdateColumn(simple.DB(), id, name, value)
}

func (s *tweetService) Delete(id int64) {
	repositories.TweetRepository.Delete(simple.DB(), id)
}
