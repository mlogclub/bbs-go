package api

import (
	newCaptcha "bbs-go/internal/pkg/captcha"
	"bytes"
	"encoding/base64"

	"github.com/dchest/captcha"
	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple/web"
)

type CaptchaController struct {
	Ctx iris.Context
}

func (c *CaptchaController) GetRequest() *web.JsonResult {
	captchaId := captcha.NewLen(4)
	var buf bytes.Buffer
	if err := captcha.WriteImage(&buf, captchaId, captcha.StdWidth, captcha.StdHeight); err != nil {
		return web.JsonError(err)
	}
	return web.NewEmptyRspBuilder().
		Put("captchaId", captchaId).
		Put("captchaBase64", base64.StdEncoding.EncodeToString(buf.Bytes())).
		JsonResult()
}

func (c *CaptchaController) GetVerify() *web.JsonResult {
	captchaId := c.Ctx.URLParam("captchaId")
	captchaCode := c.Ctx.URLParam("captchaCode")
	success := captcha.VerifyString(captchaId, captchaCode)
	return web.NewEmptyRspBuilder().Put("success", success).JsonResult()
}

func (c *CaptchaController) GetRequest_angle() *web.JsonResult {
	data, err := newCaptcha.Generate()
	if err != nil {
		return web.JsonError(err)
	}
	return web.JsonData(data)
}
