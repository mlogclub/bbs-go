package es

import (
	"bbs-go/model"
	"context"
	"github.com/mlogclub/simple"
	"github.com/mlogclub/simple/json"
	"github.com/olivere/elastic/v7"
	"github.com/sirupsen/logrus"
	"strconv"
)

type TopicDocument struct {
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

func (t *TopicDocument) ToStr() string {
	str, err := json.ToStr(t)
	if err != nil {
		logrus.Error(err)
	}
	return str
}

func NewTopicDoc(topic *model.Topic) *TopicDocument {
	if topic == nil {
		return nil
	}
	return &TopicDocument{
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

func UpdateTopicIndex(topic *model.Topic) {
	if initClient() == nil {
		logrus.Error(noConfigErr)
		return
	}
	doc := NewTopicDoc(topic)
	if doc == nil {
		logrus.Error("Topic doc is null. ")
		return
	}
	logrus.Infof("Es add index topic, id = %d", topic.Id)
	if response, err := es.Index().
		Index(index).
		BodyJson(doc).
		Id(strconv.FormatInt(doc.Id, 10)).
		Do(context.Background()); err == nil {
		logrus.Info(response.Result)
	} else {
		logrus.Error(err)
	}
}

func SearchTopic(keyword string, page, limit int) (docs []TopicDocument, paging *simple.Paging, err error) {
	if initClient() == nil {
		err = noConfigErr
		return
	}

	paging = &simple.Paging{Page: page, Limit: limit}

	query := elastic.NewMultiMatchQuery(keyword, "title", "content")
	searchResult, err := es.Search().
		Index(index).
		Query(query).
		From(paging.Offset()).Size(paging.Limit).
		Do(context.Background())
	if err != nil {
		return
	}
	logrus.Infof("Query took %d milliseconds\n", searchResult.TookInMillis)

	if totalHits := searchResult.TotalHits(); totalHits > 0 {
		paging.Total = totalHits
		for _, hit := range searchResult.Hits.Hits {
			var doc TopicDocument
			if json.Parse(string(hit.Source), &doc) != nil {
				docs = append(docs, doc)
			}
		}
	}
	return
}
