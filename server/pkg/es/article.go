package es

import (
	"bbs-go/cache"
	"bbs-go/model"
	"bbs-go/model/constants"
	html2 "bbs-go/pkg/html"
	"bbs-go/pkg/markdown"
	"bbs-go/repositories"
	"context"
	"html"
	"strconv"
	"time"

	"github.com/mlogclub/simple/common/dates"
	"github.com/mlogclub/simple/common/jsons"
	"github.com/mlogclub/simple/common/strs"
	"github.com/mlogclub/simple/sqls"

	"github.com/olivere/elastic/v7"
	"github.com/sirupsen/logrus"
)

type ArticleDocument struct {
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

func (t *ArticleDocument) ToStr() string {
	str, err := jsons.ToStr(t)
	if err != nil {
		logrus.Error(err)
	}
	return str
}

func NewArticleDoc(article *model.Article) *ArticleDocument {
	if article == nil {
		return nil
	}
	doc := &ArticleDocument{
		Id:         article.Id,
		UserId:     article.UserId,
		Title:      article.Title,
		Status:     article.Status,
		CreateTime: article.CreateTime,
	}

	// 处理内容
	content := markdown.ToHTML(article.Content)
	content = html2.GetHtmlText(content)
	content = html.EscapeString(content)

	doc.Content = content

	// 处理用户
	user := cache.UserCache.Get(article.UserId)
	if user != nil {
		doc.Nickname = user.Nickname
	}

	// 处理标签
	tags := getArticleTags(article.Id)
	var tagsArr []string
	for _, tag := range tags {
		tagsArr = append(tagsArr, tag.Name)
	}
	doc.Tags = tagsArr

	return doc
}

func getArticleTags(topicId int64) []model.Tag {
	articleTags := repositories.ArticleTagRepository.Find(sqls.DB(), sqls.NewCnd().Where("article_id = ?", topicId))

	var tagIds []int64
	for _, articleTag := range articleTags {
		tagIds = append(tagIds, articleTag.TagId)
	}
	return cache.TagCache.GetList(tagIds)
}

func UpdateArticleIndexAsync(topic *model.Article) {
	if err := indexPool.Submit(func() {
		UpdateArticleIndex(topic)
	}); err != nil {
		logrus.Error(err)
	}
}

func UpdateArticleIndex(article *model.Article) {
	if article == nil {
		return
	}
	if initClient() == nil {
		logrus.Error(errNoConfig)
		return
	}
	doc := NewArticleDoc(article)
	if doc == nil {
		logrus.Error("Topic doc is null. ")
		return
	}
	logrus.Infof("Es add index topic, id = %d", article.Id)
	if response, err := client.Index().
		Index(indexArticles).
		BodyJson(doc).
		Id(strconv.FormatInt(doc.Id, 10)).
		Do(context.Background()); err == nil {
		logrus.Info(response.Result)
	} else {
		logrus.Error(err)
	}
}

func SearchArticle(keyword string, timeRange, page, limit int) (docs []ArticleDocument, paging *sqls.Paging, err error) {
	if initClient() == nil {
		err = errNoConfig
		return
	}

	paging = &sqls.Paging{Page: page, Limit: limit}

	query := elastic.NewBoolQuery().
		Must(elastic.NewTermQuery("status", constants.StatusOk))

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
		Index(indexArticles).
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
			var doc ArticleDocument
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
				logrus.Error(err)
			}
		}
	}
	return
}
