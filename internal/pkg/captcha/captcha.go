package captcha

import (
	"log"

	"github.com/mlogclub/simple/common/strs"
	"github.com/spf13/cast"
	"github.com/wenlng/go-captcha-assets/resources/images"
	"github.com/wenlng/go-captcha/v2/base/option"
	"github.com/wenlng/go-captcha/v2/rotate"
)

var rotateCapt rotate.Captcha

func init() {
	builder := rotate.NewBuilder(
		rotate.WithRangeAnglePos([]option.RangeVal{
			{Min: 20, Max: 330},
		}),
	)

	// background images
	images, err := images.GetImages()
	if err != nil {
		log.Fatalln(err)
	}

	// set resources
	builder.SetResources(
		rotate.WithImages(images),
	)

	rotateCapt = builder.Make()
}

func Generate() (*CaptchaData, error) {
	data, err := rotateCapt.Generate()
	if err != nil {
		return nil, err
	}

	imageBase64, err := data.GetMasterImage().ToBase64()
	if err != nil {
		return nil, err
	}
	thumbBase64, err := data.GetThumbImage().ToBase64()
	if err != nil {
		return nil, err
	}

	id := strs.UUID()

	Set(id, data.GetData())

	return &CaptchaData{
		Id:          id,
		ImageBase64: imageBase64,
		ThumbBase64: thumbBase64,
		ThumbSize:   data.GetData().Width,
	}, nil
}

func Verify(captchaId string, captchaCode string) bool {
	data := Get(captchaId)
	if data == nil {
		return false
	}
	angle := cast.ToFloat64(captchaCode)
	return rotate.CheckAngle(int64(angle), int64(data.Angle), 2)
}
