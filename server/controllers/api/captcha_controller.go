package api

import (
	"bbs-go/pkg/bbsurls"

	"github.com/dchest/captcha"
	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple/common/strs"
	"github.com/mlogclub/simple/mvc"
	"github.com/sirupsen/logrus"
)

type CaptchaController struct {
	Ctx iris.Context
}

func (c *CaptchaController) GetRequest() *mvc.JsonResult {
	captchaId := c.Ctx.FormValue("captchaId")
	if strs.IsNotBlank(captchaId) { // reload
		if !captcha.Reload(captchaId) {
			// reload 失败，重新加载验证码
			captchaId = captcha.NewLen(4)
		}
	} else {
		captchaId = captcha.NewLen(4)
	}
	captchaUrl := bbsurls.AbsUrl("/api/captcha/show?captchaId=" + captchaId + "&r=" + strs.UUID())
	return mvc.NewEmptyRspBuilder().
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

func (c *CaptchaController) GetVerify() *mvc.JsonResult {
	captchaId := c.Ctx.URLParam("captchaId")
	captchaCode := c.Ctx.URLParam("captchaCode")
	success := captcha.VerifyString(captchaId, captchaCode)
	return mvc.NewEmptyRspBuilder().Put("success", success).JsonResult()
}
