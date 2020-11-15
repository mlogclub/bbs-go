package es

import (
	"bbs-go/model"
	"github.com/mlogclub/simple/json"
	"github.com/sirupsen/logrus"
	"strconv"
)

type Document struct {
	Id              string `json:"id"`
	OriginId        int64  `json:"originId"`
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

func (t *Document) ToStr() string {
	str, err := json.ToStr(t)
	if err != nil {
		logrus.Error(err)
	}
	return str
}

func NewTopicDoc(topic *model.Topic) *Document {
	if topic == nil {
		return nil
	}
	return &Document{
		Id:              "topic-" + strconv.FormatInt(topic.Id, 10),
		OriginId:        topic.Id,
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

func NewArticleDoc(article *model.Article) *Document {
	if article == nil {
		return nil
	}
	return &Document{
		Id:         "article-" + strconv.FormatInt(article.Id, 10),
		OriginId:   article.Id,
		UserId:     article.UserId,
		Title:      article.Title,
		Content:    article.Content,
		Status:     article.Status,
		CreateTime: article.CreateTime,
	}
}
