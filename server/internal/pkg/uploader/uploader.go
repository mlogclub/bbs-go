package uploader

import (
	"mime"
	"time"

	"bbs-go/internal/models/dto"
	"bbs-go/internal/pkg/config"

	"github.com/go-resty/resty/v2"

	"bbs-go/internal/pkg/simple/common/dates"
	"bbs-go/internal/pkg/simple/common/digests"
	"bbs-go/internal/pkg/simple/common/strs"
)

type Uploader interface {
	PutImage(cfg dto.UploadConfig, data []byte, contentType string) (string, error)
	PutObject(cfg dto.UploadConfig, key string, data []byte, contentType string) (string, error)
	CopyImage(cfg dto.UploadConfig, originUrl string) (string, error)
}

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
	if config.IsProd() {
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
