package es

import (
	"bbs-go/internal/cache"
	"bbs-go/internal/models"
	"bbs-go/internal/models/constants"
	html2 "bbs-go/internal/pkg/html"
	"bbs-go/internal/pkg/markdown"
	"bbs-go/internal/repositories"
	"context"
	"html"
	"log/slog"
	"strconv"
	"time"

	"github.com/mlogclub/simple/common/dates"
	"github.com/mlogclub/simple/common/jsons"
	"github.com/mlogclub/simple/common/strs"
	"github.com/mlogclub/simple/sqls"

	"github.com/olivere/elastic/v7"
	"github.com/panjf2000/ants/v2"

	"github.com/sirupsen/logrus"
)

var indexPool, _ = ants.NewPool(8)

type TopicDocument struct {
	Id         int64    `json:"id"`
	NodeId     int64    `json:"nodeId"`
	UserId     int64    `json:"userId"`
	Nickname   string   `json:"nickname"`
	Title      string   `json:"title"`
	Content    string   `json:"content"`
	Tags       []string `json:"tags"`
	Recommend  bool     `json:"recommend"`
	Status     int      `json:"status"`
	CreateTime int64    `json:"createTime"`
}

func (t *TopicDocument) ToStr() string {
	str, err := jsons.ToStr(t)
	if err != nil {
		slog.Error(err.Error(), slog.Any("err", err))
	}
	return str
}

func NewTopicDoc(topic *models.Topic) *TopicDocument {
	if topic == nil {
		return nil
	}
	doc := &TopicDocument{
		Id:         topic.Id,
		NodeId:     topic.NodeId,
		UserId:     topic.UserId,
		Title:      topic.Title,
		Status:     topic.Status,
		Recommend:  topic.Recommend,
		CreateTime: topic.CreateTime,
	}

	// 处理内容
	content := markdown.ToHTML(topic.Content)
	content = html2.GetHtmlText(content)
	content = html.EscapeString(content)

	doc.Content = content

	// 处理用户
	user := cache.UserCache.Get(topic.UserId)
	if user != nil {
		doc.Nickname = user.Nickname
	}

	// 处理标签
	tags := getTopicTags(topic.Id)
	var tagsArr []string
	for _, tag := range tags {
		tagsArr = append(tagsArr, tag.Name)
	}
	doc.Tags = tagsArr

	return doc
}

func getTopicTags(topicId int64) []models.Tag {
	topicTags := repositories.TopicTagRepository.Find(sqls.DB(), sqls.NewCnd().Where("topic_id = ?", topicId))

	var tagIds []int64
	for _, topicTag := range topicTags {
		tagIds = append(tagIds, topicTag.TagId)
	}
	return cache.TagCache.GetList(tagIds)
}

func UpdateTopicIndexAsync(topic *models.Topic) {
	if err := indexPool.Submit(func() {
		UpdateTopicIndex(topic)
	}); err != nil {
		slog.Error(err.Error(), slog.Any("err", err))
	}
}

func UpdateTopicIndex(topic *models.Topic) {
	if topic == nil {
		return
	}
	if initClient() == nil {
		slog.Error("初始化ES客户端异常", slog.Any("err", errNoConfig))
		return
	}
	doc := NewTopicDoc(topic)
	if doc == nil {
		slog.Error("Topic doc is null. ")
		return
	}
	logrus.Infof("Es add index topic, id = %d", topic.Id)
	if response, err := client.Index().
		Index(index).
		BodyJson(doc).
		Id(strconv.FormatInt(doc.Id, 10)).
		Do(context.Background()); err == nil {
		slog.Info(response.Result)
	} else {
		slog.Error(err.Error(), slog.Any("err", err))
	}
}

func SearchTopic(keyword string, nodeId int64, timeRange, page, limit int) (docs []TopicDocument, paging *sqls.Paging, err error) {
	if initClient() == nil {
		err = errNoConfig
		return
	}

	paging = &sqls.Paging{Page: page, Limit: limit}

	query := elastic.NewBoolQuery().
		Must(elastic.NewTermQuery("status", constants.StatusOk))
	if nodeId != 0 {
		if nodeId == -1 { // 推荐
			query.Must(elastic.NewTermQuery("recommend", true))
		} else {
			query.Must(elastic.NewTermQuery("nodeId", nodeId))
		}
	}
	if timeRange == 1 { // 一天内
		beginTime := dates.Timestamp(time.Now().Add(-24 * time.Hour))
		query.Must(elastic.NewRangeQuery("createTime").Gte(beginTime))
	} else if timeRange == 2 { // 一周内
		beginTime := dates.Timestamp(time.Now().Add(-7 * 24 * time.Hour))
		query.Must(elastic.NewRangeQuery("createTime").Gte(beginTime))
	} else if timeRange == 3 { // 一月内
		beginTime := dates.Timestamp(time.Now().AddDate(0, -1, 0))
		query.Must(elastic.NewRangeQuery("createTime").Gte(beginTime))
	} else if timeRange == 4 { // 一年内
		beginTime := dates.Timestamp(time.Now().AddDate(-1, 0, 0))
		query.Must(elastic.NewRangeQuery("createTime").Gte(beginTime))
	}
	query.Must(elastic.NewMultiMatchQuery(keyword, "title", "content", "tags"))

	highlight := elastic.NewHighlight().
		PreTags("<span class='search-highlight'>").PostTags("</span>").
		Fields(elastic.NewHighlighterField("title"), elastic.NewHighlighterField("content"), elastic.NewHighlighterField("nickname"),
			elastic.NewHighlighterField("tags"))

	searchResult, err := client.Search().
		Index(index).
		Query(query).
		From(paging.Offset()).Size(paging.Limit).
		Highlight(highlight).
		Do(context.Background())
	if err != nil {
		return
	}
	// logrus.Infof("Query took %d milliseconds\n", searchResult.TookInMillis)

	if totalHits := searchResult.TotalHits(); totalHits > 0 {
		paging.Total = totalHits
		for _, hit := range searchResult.Hits.Hits {
			var doc TopicDocument
			if err := jsons.Parse(string(hit.Source), &doc); err == nil {
				if len(hit.Highlight["title"]) > 0 && strs.IsNotBlank(hit.Highlight["title"][0]) {
					doc.Title = hit.Highlight["title"][0]
				}
				if len(hit.Highlight["content"]) > 0 && strs.IsNotBlank(hit.Highlight["content"][0]) {
					doc.Content = hit.Highlight["content"][0]
				} else {
					doc.Content = html2.GetSummary(doc.Content, 128)
				}
				if len(hit.Highlight["nickname"]) > 0 && strs.IsNotBlank(hit.Highlight["nickname"][0]) {
					doc.Nickname = hit.Highlight["nickname"][0]
				} else if len(hit.Highlight["tags"]) > 0 {
					doc.Tags = hit.Highlight["tags"]
				}
				docs = append(docs, doc)
			} else {
				slog.Error(err.Error(), slog.Any("err", err))
			}
		}
	}
	return
}
