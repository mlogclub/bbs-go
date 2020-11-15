package es

import (
	"bbs-go/model"
	"github.com/mlogclub/simple/json"
	"github.com/sirupsen/logrus"
)

type TopicDoc struct {
	Id              int64  `json:"id"`
	NodeId          int64  `json:"nodeId"`
	UserId          int64  `json:"userId"`
	Title           string `json:"title"`
	Content         string `json:"content"`
	Recommend       bool   `json:"recommend"`
	LastCommentTime int64  `json:"lastCommentTime"`
	Status          int    `json:"status"`
	CommentCount    int64  `json:"commentCount"`
	LikeCount       int64  `json:"likeCount"`
	CreateTime      int64  `json:"createTime"`
}

func (t *TopicDoc) ToStr() string {
	str, err := json.ToStr(t)
	if err != nil {
		logrus.Error(err)
	}
	return str
}

func NewTopicDoc(topic *model.Topic) *TopicDoc {
	if topic == nil {
		return nil
	}
	return &TopicDoc{
		Id:              topic.Id,
		NodeId:          topic.NodeId,
		UserId:          topic.UserId,
		Title:           topic.Title,
		Content:         topic.Content,
		Recommend:       topic.Recommend,
		LastCommentTime: topic.LastCommentTime,
		Status:          topic.Status,
		CommentCount:    topic.CommentCount,
		LikeCount:       topic.LikeCount,
		CreateTime:      topic.CreateTime,
	}
}
