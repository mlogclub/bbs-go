package es

import (
	"errors"
	"server/pkg/config"
	"sync"

	"github.com/mlogclub/simple/common/strs"
	"github.com/olivere/elastic/v7"
	"github.com/sirupsen/logrus"
)

var (
	client        *elastic.Client
	indexTopic    string
	indexArticles string
	once          sync.Once
	errNoConfig   = errors.New("es config not found. ")
)

func initClient() *elastic.Client {
	once.Do(func() {
		var err error
		if !strs.IsAnyBlank(config.Instance.Es.Url, config.Instance.Es.IndexTopic) {
			indexTopic = config.Instance.Es.IndexTopic
			indexArticles = config.Instance.Es.IndexArticles
			client, err = elastic.NewClient(
				elastic.SetURL(config.Instance.Es.Url),
				elastic.SetHealthcheck(false),
				elastic.SetSniff(false),
			)
		} else {
			err = errNoConfig
		}
		if err != nil {
			logrus.Error(err)
		}
	})
	return client
}
