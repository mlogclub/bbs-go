package api

import (
	"bbs-go/package/urls"
	"github.com/dchest/captcha"
	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple"
	"github.com/sirupsen/logrus"
)

type CaptchaController struct {
	Ctx iris.Context
}

func (c *CaptchaController) GetRequest() *simple.JsonResult {
	captchaId := c.Ctx.FormValue("captchaId")
	if simple.IsNotBlank(captchaId) { // reload
		if !captcha.Reload(captchaId) {
			// reload 失败，重新加载验证码
			captchaId = captcha.NewLen(4)
		}
	} else {
		captchaId = captcha.NewLen(4)
	}
	captchaUrl := urls.AbsUrl("/api/captcha/show?captchaId=" + captchaId + "&r=" + simple.UUID())
	return simple.NewEmptyRspBuilder().
		Put("captchaId", captchaId).
		Put("captchaUrl", captchaUrl).
		JsonResult()
}

func (c *CaptchaController) GetShow() {
	captchaId := c.Ctx.URLParam("captchaId")

	if captchaId == "" {
		c.Ctx.StatusCode(404)
		return
	}

	if !captcha.Reload(captchaId) {
		c.Ctx.StatusCode(404)
		return
	}

	c.Ctx.Header("Content-Type", "image/png")
	if err := captcha.WriteImage(c.Ctx.ResponseWriter(), captchaId, captcha.StdWidth, captcha.StdHeight); err != nil {
		logrus.Error(err)
	}
}

func (c *CaptchaController) GetVerify() *simple.JsonResult {
	captchaId := c.Ctx.URLParam("captchaId")
	captchaCode := c.Ctx.URLParam("captchaCode")
	success := captcha.VerifyString(captchaId, captchaCode)
	return simple.NewEmptyRspBuilder().Put("success", success).JsonResult()
}
