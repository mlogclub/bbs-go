package services

import (
	"bbs-go/model"
	"bbs-go/repositories"
	"github.com/mlogclub/simple"
)

var EmailCodeService = newEmailCodeService()

func newEmailCodeService() *emailCodeService {
	return &emailCodeService{}
}

type emailCodeService struct {
}

func (s *emailCodeService) Get(id int64) *model.EmailCode {
	return repositories.EmailCodeRepository.Get(simple.DB(), id)
}

func (s *emailCodeService) Take(where ...interface{}) *model.EmailCode {
	return repositories.EmailCodeRepository.Take(simple.DB(), where...)
}

func (s *emailCodeService) Find(cnd *simple.SqlCnd) []model.EmailCode {
	return repositories.EmailCodeRepository.Find(simple.DB(), cnd)
}

func (s *emailCodeService) FindOne(cnd *simple.SqlCnd) *model.EmailCode {
	return repositories.EmailCodeRepository.FindOne(simple.DB(), cnd)
}

func (s *emailCodeService) FindPageByParams(params *simple.QueryParams) (list []model.EmailCode, paging *simple.Paging) {
	return repositories.EmailCodeRepository.FindPageByParams(simple.DB(), params)
}

func (s *emailCodeService) FindPageByCnd(cnd *simple.SqlCnd) (list []model.EmailCode, paging *simple.Paging) {
	return repositories.EmailCodeRepository.FindPageByCnd(simple.DB(), cnd)
}

func (s *emailCodeService) Count(cnd *simple.SqlCnd) int64 {
	return repositories.EmailCodeRepository.Count(simple.DB(), cnd)
}

func (s *emailCodeService) Create(t *model.EmailCode) error {
	return repositories.EmailCodeRepository.Create(simple.DB(), t)
}

func (s *emailCodeService) Update(t *model.EmailCode) error {
	return repositories.EmailCodeRepository.Update(simple.DB(), t)
}

func (s *emailCodeService) Updates(id int64, columns map[string]interface{}) error {
	return repositories.EmailCodeRepository.Updates(simple.DB(), id, columns)
}

func (s *emailCodeService) UpdateColumn(id int64, name string, value interface{}) error {
	return repositories.EmailCodeRepository.UpdateColumn(simple.DB(), id, name, value)
}

func (s *emailCodeService) Delete(id int64) {
	repositories.EmailCodeRepository.Delete(simple.DB(), id)
}

// 发送邮箱验证码
func (s *emailCodeService) SendVerifyEmail(userId int64, emailStr, title, template string) error {
	// email.SendTemplateEmail()
	// email.SendTemplateEmail(email, , "邮箱验证", "邮件内容:"+code, "", "https://mlog.club")
	// TODO 发送验证码
	// TODO 插入EmailCode记录
	// email.SendTemplateEmail(emailStr)
	return nil
}
