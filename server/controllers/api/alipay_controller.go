package api

import (
	alipay2 "bbs-go/pkg/alipay"
	"bbs-go/services"
	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple/web"
	"github.com/mlogclub/simple/web/params"
	"github.com/smartwalle/alipay/v3"
	"net/url"
)

type AlipayController struct {
	Ctx iris.Context
}

func (c *AlipayController) GetUrl() *web.JsonResult {
	user := services.UserTokenService.GetCurrent(c.Ctx)
	if err := services.UserService.CheckPostStatus(user); err != nil {
		return web.JsonError(err)
	}

	price := params.FormValue(c.Ctx, "price")
	urls := services.AlipayService.GetUrl(user.Id, price)
	return web.JsonData(urls)
}

func (c *AlipayController) GetVerify() *web.JsonResult {
	x, _ := url.Parse(c.Ctx.Request().URL.String())
	sign := services.AlipayService.PayVerify(x.Query())
	return web.JsonData(sign)
}

func (c *AlipayController) PostNotification() {
	var noti, _ = alipay2.GetClient().GetTradeNotification(c.Ctx.Request())
	if noti != nil {
		services.AlipayService.NotifyVerifyOrder(noti.OutTradeNo, noti.TradeStatus)
	}
	alipay.AckNotification(c.Ctx.ResponseWriter()) // 确认收到通知消息
}
