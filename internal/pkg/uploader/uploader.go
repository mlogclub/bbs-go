package uploader

import (
	"bbs-go/internal/models/dto"
	"bbs-go/internal/pkg/config"
	"mime"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/mlogclub/simple/common/dates"
	"github.com/mlogclub/simple/common/digests"
	"github.com/mlogclub/simple/common/strs"
)

type Uploader interface {
	PutImage(cfg dto.UploadConfig, data []byte, contentType string) (string, error)
	PutObject(cfg dto.UploadConfig, key string, data []byte, contentType string) (string, error)
	CopyImage(cfg dto.UploadConfig, originUrl string) (string, error)
}

// generateKey 生成图片Key
func generateImageKey(data []byte, contentType string) string {
	if config.IsProd() {
		return generateImageKeyByPrefix("images", data, contentType)
	}
	return generateImageKeyByPrefix("test/images", data, contentType)
}

func generateImageKeyByPrefix(prefix string, data []byte, contentType string) string {
	md5 := digests.MD5Bytes(data)
	ext := getImageExt(contentType)
	datePath := dates.Format(time.Now(), "2006/01/02/")
	if strs.IsBlank(prefix) {
		return datePath + md5 + ext
	}
	cleanPrefix := strings.Trim(strings.TrimSpace(prefix), "/")
	return cleanPrefix + "/" + datePath + md5 + ext
}

func getImageExt(contentType string) string {
	if strs.IsBlank(contentType) {
		return ""
	}

	// 先解析 Content-Type，去掉参数（如 charset=utf-8）
	mediaType, _, _ := mime.ParseMediaType(contentType)
	if mediaType == "" {
		return ""
	}

	// 处理一些非标准的 MIME 类型
	switch mediaType {
	case "image/jfif":
		return ".jpg"
	case "image/pjpeg":
		return ".jpg"
	default:
		exts, _ := mime.ExtensionsByType(mediaType)
		if len(exts) > 0 {
			// 对于 image/jpeg，优先使用 .jpg 而不是 .jpe
			if mediaType == "image/jpeg" {
				for _, e := range exts {
					if e == ".jpg" {
						return ".jpg"
					}
				}
			}
			return exts[0]
		}
	}
	return ""
}

func download(url string) ([]byte, string, error) {
	rsp, err := resty.New().R().Get(url)
	if err != nil {
		return nil, "", err
	}
	return rsp.Body(), rsp.Header().Get("Content-Type"), nil
}
