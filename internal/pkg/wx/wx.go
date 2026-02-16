package wx

import (
	"github.com/silenceper/wechat/v2/cache"
	"github.com/silenceper/wechat/v2/officialaccount"
	offConfig "github.com/silenceper/wechat/v2/officialaccount/config"
)

// 微信文档：https://developers.weixin.qq.com/doc/oplatform/Website_App/WeChat_Login/Wechat_Login.html
// SDK文档：https://silenceper.com/wechat/openplatform/

// var (
// 	wc *wechat.Wechat
// )
// func init() {
// 	wc = wechat.NewWechat()
// 	wc.SetCache(cache.NewMemory())
// }

var (
	memCache = cache.NewMemory()
)

func NewOfficialAccount(appId, appSecret string) *officialaccount.OfficialAccount {
	return officialaccount.NewOfficialAccount(&offConfig.Config{
		AppID:     appId,
		AppSecret: appSecret,
		Cache:     memCache,
	})
}
