package captcha

import (
	"time"

	"github.com/goburrow/cache"
	"github.com/wenlng/go-captcha/v2/rotate"
)

var captchaCache cache.Cache

func init() {
	captchaCache = cache.New(
		cache.WithMaximumSize(1000),
		cache.WithExpireAfterAccess(10*time.Minute),
	)
}

func Get(captchaId string) *rotate.Block {
	if v, ok := captchaCache.GetIfPresent(captchaId); ok {
		return v.(*rotate.Block)
	}
	return nil
}

func Set(captchaId string, captcha *rotate.Block) {
	captchaCache.Put(captchaId, captcha)
}
