package render

import (
	"bbs-go/internal/models"
	"log/slog"

	"github.com/mlogclub/simple/common/jsons"
	"github.com/mlogclub/simple/common/strs"
)

func BuildImageList(imageListStr string) (imageList []models.ImageInfo) {
	if strs.IsNotBlank(imageListStr) {
		var images []models.ImageDTO
		if err := jsons.Parse(imageListStr, &images); err == nil {
			if len(images) > 0 {
				for _, image := range images {
					imageList = append(imageList, models.ImageInfo{
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

func BuildImage(imageStr string) *models.ImageInfo {
	if strs.IsBlank(imageStr) {
		return nil
	}
	var img *models.ImageDTO
	if err := jsons.Parse(imageStr, &img); err != nil {
		slog.Error(err.Error(), slog.Any("err", err))
		return nil
	} else {
		return &models.ImageInfo{
			Url:     HandleOssImageStyleDetail(img.Url),
			Preview: HandleOssImageStylePreview(img.Url),
		}
	}
}
