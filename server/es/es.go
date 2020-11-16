package es

import (
	"bbs-go/model"
	"context"
	"errors"
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
