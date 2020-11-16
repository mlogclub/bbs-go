package es

import (
	"bbs-go/model"
	"context"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/mlogclub/simple/json"
	"github.com/sirupsen/logrus"
	"log"
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

func initClient() *elasticsearch.Client {
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
	initClient()
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

func Search(keyword string) {
	initClient()
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				"title": "test",
			},
		},
	}
	body, err := json.ToStr(query)
	if err != nil {
		logrus.Errorf("Error encoding query: %s", err)
		return
	}

	response, err := es.Search(
		es.Search.WithContext(context.Background()),
		es.Search.WithIndex(TopicIndexName),
		es.Search.WithBody(strings.NewReader(body)),
	)
	if err != nil {
		logrus.Error(err)
		return
	}
	defer func() {
		_ = response.Body.Close()
	}()

	if response.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			log.Fatalf("Error parsing the response body: %s", err)
		} else {
			// Print the response status and error information.
			log.Fatalf("[%s] %s: %s",
				res.Status(),
				e["error"].(map[string]interface{})["type"],
				e["error"].(map[string]interface{})["reason"],
			)
		}
	}

	// TODO 搜索
}
