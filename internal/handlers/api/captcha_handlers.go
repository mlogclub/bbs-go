package api

import (
	newCaptcha "bbs-go/internal/pkg/captcha"
	"bytes"
	"encoding/base64"

	"github.com/gin-gonic/gin"

	"bbs-go/internal/pkg/ginx"

	"github.com/dchest/captcha"
)

func CaptchaRequest(ctx *gin.Context) {

	captchaId := captcha.NewLen(4)
	var buf bytes.Buffer
	if err := captcha.WriteImage(&buf, captchaId, captcha.StdWidth, captcha.StdHeight); err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}
	ginx.WriteJSON(ctx, map[string]any{
		"captchaId":     captchaId,
		"captchaBase64": base64.StdEncoding.EncodeToString(buf.Bytes()),
	})

}

func CaptchaVerify(ctx *gin.Context) {
	captchaId := ctx.Query("captchaId")
	captchaCode := ctx.Query("captchaCode")
	success := captcha.VerifyString(captchaId, captchaCode)
	ginx.WriteJSON(ctx, map[string]any{"success": success})

}

func CaptchaRequestAngle(ctx *gin.Context) {

	data, err := newCaptcha.Generate()
	if err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}
	ginx.WriteJSON(ctx, data)

}
