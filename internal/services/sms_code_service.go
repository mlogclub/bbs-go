package services

import (
	"bbs-go/internal/models"
	"bbs-go/internal/pkg/sms"
	"bbs-go/internal/repositories"
	"errors"
	"math/rand"

	"github.com/mlogclub/simple/common/dates"
	"github.com/mlogclub/simple/common/strs"
	"github.com/mlogclub/simple/sqls"
	"github.com/mlogclub/simple/web/params"
	"github.com/spf13/cast"
)

var SmsCodeService = newSmsCodeService()

func newSmsCodeService() *smsCodeService {
	return &smsCodeService{}
}

type smsCodeService struct {
}

func (s *smsCodeService) Get(id int64) *models.SmsCode {
	return repositories.SmsCodeRepository.Get(sqls.DB(), id)
}

func (s *smsCodeService) Take(where ...interface{}) *models.SmsCode {
	return repositories.SmsCodeRepository.Take(sqls.DB(), where...)
}

func (s *smsCodeService) Find(cnd *sqls.Cnd) []models.SmsCode {
	return repositories.SmsCodeRepository.Find(sqls.DB(), cnd)
}

func (s *smsCodeService) FindOne(cnd *sqls.Cnd) *models.SmsCode {
	return repositories.SmsCodeRepository.FindOne(sqls.DB(), cnd)
}

func (s *smsCodeService) FindPageByParams(params *params.QueryParams) (list []models.SmsCode, paging *sqls.Paging) {
	return repositories.SmsCodeRepository.FindPageByParams(sqls.DB(), params)
}

func (s *smsCodeService) FindPageByCnd(cnd *sqls.Cnd) (list []models.SmsCode, paging *sqls.Paging) {
	return repositories.SmsCodeRepository.FindPageByCnd(sqls.DB(), cnd)
}

func (s *smsCodeService) Count(cnd *sqls.Cnd) int64 {
	return repositories.SmsCodeRepository.Count(sqls.DB(), cnd)
}

func (s *smsCodeService) Create(t *models.SmsCode) error {
	return repositories.SmsCodeRepository.Create(sqls.DB(), t)
}

func (s *smsCodeService) Update(t *models.SmsCode) error {
	return repositories.SmsCodeRepository.Update(sqls.DB(), t)
}

func (s *smsCodeService) Updates(id int64, columns map[string]interface{}) error {
	return repositories.SmsCodeRepository.Updates(sqls.DB(), id, columns)
}

func (s *smsCodeService) UpdateColumn(id int64, name string, value interface{}) error {
	return repositories.SmsCodeRepository.UpdateColumn(sqls.DB(), id, name, value)
}

func (s *smsCodeService) Delete(id int64) {
	repositories.SmsCodeRepository.Delete(sqls.DB(), id)
}

type SmsCode struct {
	Code  string
	Phone string
}

func (s *smsCodeService) SendSms(phone string) (string, error) {
	var (
		smsId = strs.UUID()
		code  = cast.ToString(rand.Intn(900000) + 100000)
		cfg   = SysConfigService.GetLoginConfig().SmsLogin.Aliyun
	)

	if err := sms.SendSmsCode(cfg, phone, code); err != nil {
		return "", err
	}

	s.Create(&models.SmsCode{
		SmsId:      smsId,
		Phone:      phone,
		Code:       code,
		ExpireAt:   dates.NowTimestamp() + 60000,
		Status:     0,
		CreateTime: dates.NowTimestamp(),
	})
	return smsId, nil
}

func (s *smsCodeService) Verify(smsId, code string) (string, error) {
	smsCode := s.FindOne(sqls.NewCnd().Where("sms_id = ? and status = 0", smsId))
	if smsCode == nil {
		return "", errors.New("验证码错误")
	}
	if smsCode.Code != code {
		return "", errors.New("验证码错误")
	}

	s.UpdateColumn(smsCode.Id, "status", 1)

	return smsCode.Phone, nil
}
