package es

import (
	"bbs-go/model"
	"bbs-go/model/constants"
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
	ViewCount       int64  `json:"viewCount"`
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
		ViewCount:       topic.ViewCount,
		CommentCount:    topic.CommentCount,
		LikeCount:       topic.LikeCount,
		CreateTime:      topic.CreateTime,
	}
}

func UpdateTopicIndex(topic *model.Topic) {
	if true { // 默认全部关闭
		return
	}
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

	query := elastic.NewBoolQuery().
		Must(elastic.NewTermQuery("status", constants.StatusOk)).
		Must(elastic.NewMultiMatchQuery(keyword, "title", "content"))

	highlight := elastic.NewHighlight().
		PreTags("<span class='search-highlight'>").PostTags("</span>").
		Fields(elastic.NewHighlighterField("title"), elastic.NewHighlighterField("content"))

	searchResult, err := es.Search().
		Index(index).
		Query(query).
		From(paging.Offset()).Size(paging.Limit).
		Highlight(highlight).
		Do(context.Background())
	if err != nil {
		return
	}
	logrus.Infof("Query took %d milliseconds\n", searchResult.TookInMillis)

	if totalHits := searchResult.TotalHits(); totalHits > 0 {
		paging.Total = totalHits
		for _, hit := range searchResult.Hits.Hits {
			var doc TopicDocument
			if err := json.Parse(string(hit.Source), &doc); err == nil {
				if len(hit.Highlight["title"]) > 0 && simple.IsNotBlank(hit.Highlight["title"][0]) {
					doc.Title = hit.Highlight["title"][0]
				}
				if len(hit.Highlight["content"]) > 0 && simple.IsNotBlank(hit.Highlight["content"][0]) {
					doc.Content = hit.Highlight["content"][0]
				} else {
					// 如果内容没有高亮的，那么就不显示内容
					doc.Content = ""
				}
				docs = append(docs, doc)
			} else {
				logrus.Error(err)
			}
		}
	}
	return
}
