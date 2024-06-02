package uploader

import (
	"bbs-go/internal/pkg/config"
	"mime"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/mlogclub/simple/common/dates"
	"github.com/mlogclub/simple/common/digests"
	"github.com/mlogclub/simple/common/strs"
)

// generateKey 生成图片Key
func generateImageKey(data []byte, contentType string) string {
	md5 := digests.MD5Bytes(data)
	ext := ""
	if strs.IsNotBlank(contentType) {
		exts, err := mime.ExtensionsByType(contentType)
		if err == nil || len(exts) > 0 {
			ext = exts[0]
		}
	}
	if config.Instance.IsProd() {
		return "images/" + dates.Format(time.Now(), "2006/01/02/") + md5 + ext
	} else {
		return "test/images/" + dates.Format(time.Now(), "2006/01/02/") + md5 + ext
	}
}

func download(url string) ([]byte, string, error) {
	rsp, err := resty.New().R().Get(url)
	if err != nil {
		return nil, "", err
	}
	return rsp.Body(), rsp.Header().Get("Content-Type"), nil
}
