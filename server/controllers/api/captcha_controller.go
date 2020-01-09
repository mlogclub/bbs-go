package api

import (
	"github.com/dchest/captcha"
	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple"
	"github.com/sirupsen/logrus"

	"bbs-go/common/urls"
)

type CaptchaController struct {
	Ctx iris.Context
}

func (this *CaptchaController) GetRequest() *simple.JsonResult {
	captchaId := captcha.NewLen(4)
	captchaUrl := urls.AbsUrl("/api/captcha/show?captchaId=" + captchaId)
	return simple.NewEmptyRspBuilder().
		Put("captchaId", captchaId).
		Put("captchaUrl", captchaUrl).
		JsonResult()
}

func (this *CaptchaController) GetShow() {
	captchaId := this.Ctx.URLParam("captchaId")

	if captchaId == "" {
		this.Ctx.StatusCode(404)
		return
	}

	if !captcha.Reload(captchaId) {
		this.Ctx.StatusCode(404)
		return
	}

	this.Ctx.Header("Content-Type", "image/png")
	if err := captcha.WriteImage(this.Ctx.ResponseWriter(), captchaId, captcha.StdWidth, captcha.StdHeight); err != nil {
		logrus.Error(err)
	}
}

func (this *CaptchaController) GetVerify() *simple.JsonResult {
	captchaId := this.Ctx.URLParam("captchaId")
	captchaCode := this.Ctx.URLParam("captchaCode")
	success := captcha.VerifyString(captchaId, captchaCode)
	return simple.NewEmptyRspBuilder().Put("success", success).JsonResult()
}
