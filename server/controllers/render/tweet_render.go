package render

import (
	"bbs-go/model"
	"github.com/mlogclub/simple"
	"github.com/mlogclub/simple/json"
	"github.com/sirupsen/logrus"
)

func BuildTweet(tweet *model.Tweet) *model.TweetResponse {
	if tweet == nil {
		return nil
	}

	rsp := &model.TweetResponse{
		TweetId:      tweet.Id,
		User:         BuildUserDefaultIfNull(tweet.UserId),
		Content:      tweet.Content,
		CommentCount: tweet.CommentCount,
		LikeCount:    tweet.LikeCount,
		Status:       tweet.Status,
		CreateTime:   tweet.CreateTime,
	}
	if simple.IsNotBlank(tweet.ImageList) {
		var images []string
		if err := json.Parse(tweet.ImageList, &images); err == nil {
			if len(images) > 0 {
				var imageList []model.ImageInfo
				for _, image := range images {
					imageList = append(imageList, model.ImageInfo{
						Url:     HandleOssImageStyleDetail(image),
						Preview: HandleOssImageStylePreview(image),
					})
				}
				rsp.ImageList = imageList
			}
		} else {
			logrus.Error(err)
		}
	}
	return rsp
}

func BuildTweets(tweets []model.Tweet) []model.TweetResponse {
	var ret []model.TweetResponse
	for _, tweet := range tweets {
		ret = append(ret, *BuildTweet(&tweet))
	}
	return ret
}
