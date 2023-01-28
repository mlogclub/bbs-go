package alipay

import (
	"github.com/sirupsen/logrus"
	"github.com/smartwalle/alipay/v3"
	"server/pkg/config"
)

var (
	client *alipay.Client
	err    error
)

func Init() {
	client, err = alipay.New(config.Instance.Alipay.AppId, config.Instance.Alipay.AppPrivateKey, true)
	if err != nil {
		logrus.Errorf("初始化支付宝支付失败 : [%s]", err.Error())
		return
	}
	err := client.LoadAliPayPublicKey(config.Instance.Alipay.AppPublicKey)
	if err != nil {
		logrus.Errorf("初始化支付宝公钥失败 : [%s]", err.Error())
		return
	}
}

func GetClient() *alipay.Client {
	return client
}
