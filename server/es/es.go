package es

import (
	"bbs-go/model"
	"context"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/sirupsen/logrus"
	"strconv"
	"strings"
	"sync"
)

var (
	es   *elasticsearch.Client
	once sync.Once
)

const (
	TopicIndexName = "bbsgo_topic"
)

func GetClient() *elasticsearch.Client {
	once.Do(func() {
		var err error
		es, err = elasticsearch.NewClient(elasticsearch.Config{
			Addresses: []string{"http://127.0.0.1:9200"},
		})
		if err != nil {
			logrus.Error(err)
		}
	})
	return es
}

func AddIndex(topic *model.Topic) {
	GetClient()
	doc := NewTopicDoc(topic)
	if doc == nil {
		return
	}
	logrus.Infof("Es add index topic, id = %d", topic.Id)
	request := esapi.IndexRequest{
		Index:      TopicIndexName,
		DocumentID: strconv.FormatInt(topic.Id, 10),
		Body:       strings.NewReader(doc.ToStr()),
	}
	response, err := request.Do(context.Background(), es)
	if err != nil {
		logrus.Error(err)
		return
	}
	defer func() {
		_ = response.Body.Close()
	}()
	if response.IsError() {
		logrus.Printf("[%s] Error indexing document ID=%d", response.Status(), topic.Id)
	} else {
		logrus.Info(response.String())
	}
}
