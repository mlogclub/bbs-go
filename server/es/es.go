package es

import (
	"bbs-go/model"
	"context"
	"errors"
	"fmt"
	"github.com/mlogclub/simple"
	"github.com/mlogclub/simple/json"
	"github.com/olivere/elastic/v7"
	"github.com/sirupsen/logrus"
	"sync"
)

var (
	es   *elastic.Client
	once sync.Once
)

const (
	TopicIndexName = "bbsgo_topic"
)

func initClient() *elastic.Client {
	once.Do(func() {
		var err error
		es, err = elastic.NewClient(
			elastic.SetURL("http://127.0.0.1:9200"),
			elastic.SetHealthcheck(false),
			elastic.SetSniff(false),
		)
		if err != nil {
			logrus.Error(err)
		}
	})
	return es
}

func AddIndex(topic *model.Topic) error {
	initClient()
	doc := NewTopicDoc(topic)
	if doc == nil {
		return errors.New("Topic doc is null. ")
	}
	logrus.Infof("Es add index topic, id = %d", topic.Id)
	if response, err := es.Index().
		Index(TopicIndexName).
		BodyJson(doc).
		Id(doc.Id).
		Do(context.Background()); err == nil {
		logrus.Error(response.Result)
		return nil
	} else {
		return err
	}
}

func Search(keyword string, page, limit int) *simple.PageResult {
	initClient()

	query := elastic.NewMultiMatchQuery(keyword, "title", "content")
	paging := &simple.Paging{Page: page, Limit: limit}
	searchResult, err := es.Search().
		Index(TopicIndexName).
		Query(query).
		From(paging.Offset()).Size(paging.Limit).
		Do(context.Background())
	if err != nil {
		logrus.Error(err)
		return nil
	}
	fmt.Printf("Query took %d milliseconds\n", searchResult.TookInMillis)

	var (
		docs      []Document
		totalHits int64
	)
	if totalHits = searchResult.TotalHits(); totalHits > 0 {
		paging.Total = totalHits
		for _, hit := range searchResult.Hits.Hits {
			var doc Document
			err := json.Parse(string(hit.Source), &doc)
			if err != nil {
				docs = append(docs, doc)
			}
		}
	}

	return &simple.PageResult{
		Page:    paging,
		Results: docs,
	}
}
