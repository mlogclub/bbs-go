package render

import (
	"bbs-go/model"

	"github.com/mlogclub/simple/common/jsons"
	"github.com/mlogclub/simple/common/strs"
	"github.com/sirupsen/logrus"
)

func buildImageList(imageListStr string) (imageList []model.ImageInfo) {
	if strs.IsNotBlank(imageListStr) {
		var images []model.ImageDTO
		if err := jsons.Parse(imageListStr, &images); err == nil {
			if len(images) > 0 {
				for _, image := range images {
					imageList = append(imageList, model.ImageInfo{
						Url:     HandleOssImageStyleDetail(image.Url),
						Preview: HandleOssImageStylePreview(image.Url),
					})
				}
			}
		} else {
			logrus.Error(err)
		}
	}
	return
}
