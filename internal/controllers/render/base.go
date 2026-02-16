package render

import (
	"bbs-go/internal/models/req"
	"bbs-go/internal/models/resp"
	"log/slog"

	"github.com/mlogclub/simple/common/jsons"
	"github.com/mlogclub/simple/common/strs"
)

func BuildImageList(imageListStr string) (imageList []resp.ImageInfo) {
	if strs.IsNotBlank(imageListStr) {
		var images []req.ImageDTO
		if err := jsons.Parse(imageListStr, &images); err == nil {
			if len(images) > 0 {
				for _, image := range images {
					imageList = append(imageList, resp.ImageInfo{
						Url:     HandleOssImageStyleDetail(image.Url),
						Preview: HandleOssImageStylePreview(image.Url),
					})
				}
			}
		} else {
			slog.Error(err.Error(), slog.Any("err", err))
		}
	}
	return
}

func BuildImage(imageStr string) *resp.ImageInfo {
	if strs.IsBlank(imageStr) {
		return nil
	}
	var img *req.ImageDTO
	if err := jsons.Parse(imageStr, &img); err != nil {
		slog.Error(err.Error(), slog.Any("err", err))
		return nil
	} else {
		return &resp.ImageInfo{
			Url:     HandleOssImageStyleDetail(img.Url),
			Preview: HandleOssImageStylePreview(img.Url),
		}
	}
}
