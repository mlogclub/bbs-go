package sms

import (
	"bbs-go/internal/models/dto"
	"log/slog"

	dysmsapi "github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
	"github.com/mlogclub/simple/common/jsons"
	"github.com/mlogclub/simple/common/strs"
)

func SendSms(cfg dto.AliyunSmsConfig, phone string, templateParam map[string]string) error {
	if strs.IsBlank(phone) {
		return nil
	}

	client, err := dysmsapi.NewClientWithAccessKey("cn-hangzhou", cfg.AccessKeyId, cfg.AccessKeySecret)
	if err != nil {
		return err
	}

	request := dysmsapi.CreateSendSmsRequest()
	request.SignName = cfg.SignName
	request.PhoneNumbers = phone
	request.TemplateCode = cfg.TemplateCode
	if len(templateParam) > 0 {
		request.TemplateParam = jsons.ToJsonStr(templateParam)
	}

	response, err := client.SendSms(request)
	if err != nil {
		slog.Error(err.Error())
		return err
	}
	slog.Info("短信发送成功", slog.Any("requestId", response.RequestId))
	return nil
}

func SendSmsCode(cfg dto.AliyunSmsConfig, phone, smsCode string) error {
	return SendSms(cfg, phone, map[string]string{
		"code": smsCode,
	})
}
