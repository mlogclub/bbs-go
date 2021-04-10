package es

import (
	"bbs-go/package/config"
	"errors"
	"github.com/mlogclub/simple"
	"github.com/olivere/elastic/v7"
	"github.com/sirupsen/logrus"
	"sync"
)

var (
	es          *elastic.Client
	index       string
	once        sync.Once
	noConfigErr = errors.New("Es config not found. ")
)

func initClient() *elastic.Client {
	once.Do(func() {
		var err error
		if !simple.IsAnyBlank(config.Instance.Es.Url, config.Instance.Es.Index) {
			index = config.Instance.Es.Index
			es, err = elastic.NewClient(
				elastic.SetURL(config.Instance.Es.Url),
				elastic.SetHealthcheck(false),
				elastic.SetSniff(false),
			)
		} else {
			err = noConfigErr
		}
		if err != nil {
			logrus.Error(err)
		}
	})
	return es
}
