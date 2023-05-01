package services

import (
	"bbs-go/model"
	alipay2 "bbs-go/pkg/alipay"
	"bbs-go/pkg/config"
	"bbs-go/pkg/order"
	"bbs-go/repositories"
	"errors"
	"fmt"
	"github.com/mlogclub/simple/common/dates"
	"github.com/mlogclub/simple/sqls"
	"github.com/sirupsen/logrus"
	"github.com/smartwalle/alipay/v3"
	"gorm.io/gorm"
	"math"
	"net/url"
	"strconv"
	"strings"
)

var AlipayService = newAlipayService()

func newAlipayService() *alipayServices {
	return &alipayServices{}
}

type alipayServices struct {
}

func (s *alipayServices) GetUrl(userId int64, price string) string {

	orderId := fmt.Sprintf("%s_%d", order.GetOrder(), userId)
	var p = alipay.TradePagePay{}
	p.ReturnURL = config.Instance.Alipay.BaseUrl + "/pay/callback"         //订单付款后跳转的网址页面
	p.NotifyURL = config.Instance.Alipay.BaseUrl + "/api/pay/notification" //通知
	p.Subject = "积分购买"                                                     //付款标题
	p.OutTradeNo = orderId                                                 //商家订单号
	p.TotalAmount = price                                                  //价格
	p.ProductCode = "FAST_INSTANT_TRADE_PAY"
	urls, err := alipay2.GetClient().TradePagePay(p)
	if err != nil {
		return ""
	}
	var o model.Order
	o.UserId = userId
	o.Subject = "积分购买"
	o.OutTradeNo = orderId
	o.TotalAmount = price
	o.CreateTime = dates.NowTimestamp()

	err = sqls.DB().Transaction(func(tx *gorm.DB) error {
		err := repositories.OrderRepository.Create(tx, &o)
		return err
	})
	if err == nil {
		return urls.String()
	}
	return ""
}

func (s *alipayServices) PayVerify(data url.Values) bool {
	sign, _ := alipay2.GetClient().VerifySign(data)
	if sign {
		orderId := data["out_trade_no"][0]
		_order := repositories.OrderRepository.Get(sqls.DB(), orderId)
		if !_order.IsSuccess {
			verifyOrder(orderId, _order.TotalAmount)
		}
	}
	return sign
}

func (s *alipayServices) NotifyVerifyOrder(orderId string, tradeStatus alipay.TradeStatus) {
	if tradeStatus == alipay.TradeStatusSuccess {
		_order := repositories.OrderRepository.Get(sqls.DB(), orderId)
		if !_order.IsSuccess {
			verifyOrder(orderId, _order.TotalAmount)
		}
	}
}

func verifyOrder(orderId, totalAmount string) {
	orders := strings.Split(orderId, "_")
	if len(orders) < 2 {
		logrus.Errorf("解析用户ID失败 : [%s]", orderId)
	}
	userId, err := strconv.ParseInt(orders[1], 10, 64)
	if err != nil {
		logrus.Errorf("用户ID转换失败 : [%s]", orderId)
	}
	price, err := getPrice(totalAmount)
	if err != nil {
		logrus.Errorf("价格转换失败 : [%s]", orderId)
	}
	err = sqls.DB().Transaction(func(tx *gorm.DB) error {
		err = UserService.PayScore(userId, price*500)
		if err != nil {
			return errors.New("充值失败")
		}
		err = repositories.OrderRepository.UpdateColumn(sqls.DB(), orderId, "is_success", true)
		if err != nil {
			return errors.New("更新订单失败")
		}
		return err
	})
	if err != nil {
		logrus.Errorf("%s : [%s]", err.Error(), orderId)
	}
}

func getPrice(price string) (int64, error) {
	_price, err := strconv.ParseFloat(price, 64)
	if err != nil {
		logrus.Errorf("价格转换失败 : [%s]", price)
		return 0, err
	}
	return int64(math.Floor(_price)), nil
}
