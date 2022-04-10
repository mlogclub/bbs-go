package services

import (
	"bbs-go/model/constants"
	"bbs-go/pkg/github"
	"bbs-go/pkg/osc"
	"bbs-go/pkg/qq"
	"database/sql"
	"strconv"
	"strings"

	"github.com/mlogclub/simple/common/dates"
	"github.com/mlogclub/simple/common/jsons"
	"github.com/mlogclub/simple/sqls"
	"github.com/mlogclub/simple/web/params"

	"bbs-go/model"
	"bbs-go/repositories"
)

var ThirdAccountService = newThirdAccountService()

func newThirdAccountService() *thirdAccountService {
	return &thirdAccountService{}
}

type thirdAccountService struct {
}

func (s *thirdAccountService) Get(id int64) *model.ThirdAccount {
	return repositories.ThirdAccountRepository.Get(sqls.DB(), id)
}

func (s *thirdAccountService) Take(where ...interface{}) *model.ThirdAccount {
	return repositories.ThirdAccountRepository.Take(sqls.DB(), where...)
}

func (s *thirdAccountService) Find(cnd *sqls.Cnd) []model.ThirdAccount {
	return repositories.ThirdAccountRepository.Find(sqls.DB(), cnd)
}

func (s *thirdAccountService) FindOne(cnd *sqls.Cnd) *model.ThirdAccount {
	return repositories.ThirdAccountRepository.FindOne(sqls.DB(), cnd)
}

func (s *thirdAccountService) FindPageByParams(params *params.QueryParams) (list []model.ThirdAccount, paging *sqls.Paging) {
	return repositories.ThirdAccountRepository.FindPageByParams(sqls.DB(), params)
}

func (s *thirdAccountService) FindPageByCnd(cnd *sqls.Cnd) (list []model.ThirdAccount, paging *sqls.Paging) {
	return repositories.ThirdAccountRepository.FindPageByCnd(sqls.DB(), cnd)
}

func (s *thirdAccountService) Create(t *model.ThirdAccount) error {
	return repositories.ThirdAccountRepository.Create(sqls.DB(), t)
}

func (s *thirdAccountService) Update(t *model.ThirdAccount) error {
	return repositories.ThirdAccountRepository.Update(sqls.DB(), t)
}

func (s *thirdAccountService) Updates(id int64, columns map[string]interface{}) error {
	return repositories.ThirdAccountRepository.Updates(sqls.DB(), id, columns)
}

func (s *thirdAccountService) UpdateColumn(id int64, name string, value interface{}) error {
	return repositories.ThirdAccountRepository.UpdateColumn(sqls.DB(), id, name, value)
}

func (s *thirdAccountService) Delete(id int64) {
	repositories.ThirdAccountRepository.Delete(sqls.DB(), id)
}

func (s *thirdAccountService) GetThirdAccount(thirdType string, thirdId string) *model.ThirdAccount {
	return repositories.ThirdAccountRepository.Take(sqls.DB(), "third_type = ? and third_id = ?", thirdType, thirdId)
}

func (s *thirdAccountService) GetOrCreateByGithub(code, state string) (*model.ThirdAccount, error) {
	userInfo, err := github.GetUserInfoByCode(code, state)
	if err != nil {
		return nil, err
	}

	account := s.GetThirdAccount(constants.ThirdAccountTypeGithub, strconv.FormatInt(userInfo.Id, 10))
	if account != nil {
		return account, nil
	}

	nickname := userInfo.Login
	if len(userInfo.Name) > 0 {
		nickname = strings.TrimSpace(userInfo.Name)
	}

	userInfoJson, _ := jsons.ToStr(userInfo)
	account = &model.ThirdAccount{
		UserId:     sql.NullInt64{},
		Avatar:     userInfo.AvatarUrl,
		Nickname:   nickname,
		ThirdType:  constants.ThirdAccountTypeGithub,
		ThirdId:    strconv.FormatInt(userInfo.Id, 10),
		ExtraData:  userInfoJson,
		CreateTime: dates.NowTimestamp(),
		UpdateTime: dates.NowTimestamp(),
	}
	err = s.Create(account)
	if err != nil {
		return nil, err
	}
	return account, nil
}

func (s *thirdAccountService) GetOrCreateByOSC(code, state string) (*model.ThirdAccount, error) {
	userInfo, err := osc.GetUserInfoByCode(code, state)
	if err != nil {
		return nil, err
	}

	account := s.GetThirdAccount(constants.ThirdAccountTypeOSC, strconv.FormatInt(userInfo.Id, 10))
	if account != nil {
		return account, nil
	}

	nickname := userInfo.Name
	if len(userInfo.Name) > 0 {
		nickname = strings.TrimSpace(userInfo.Name)
	}

	userInfoJson, _ := jsons.ToStr(userInfo)
	account = &model.ThirdAccount{
		UserId:     sql.NullInt64{},
		Avatar:     userInfo.Avatar,
		Nickname:   nickname,
		ThirdType:  constants.ThirdAccountTypeOSC,
		ThirdId:    strconv.FormatInt(userInfo.Id, 10),
		ExtraData:  userInfoJson,
		CreateTime: dates.NowTimestamp(),
		UpdateTime: dates.NowTimestamp(),
	}
	err = s.Create(account)
	if err != nil {
		return nil, err
	}
	return account, nil
}

func (s *thirdAccountService) GetOrCreateByQQ(code, state string) (*model.ThirdAccount, error) {
	userInfo, err := qq.GetUserInfoByCode(code, state)
	if err != nil {
		return nil, err
	}

	account := s.GetThirdAccount(constants.ThirdAccountTypeQQ, userInfo.Unionid)
	if account != nil {
		return account, nil
	}

	userInfoJson, _ := jsons.ToStr(userInfo)
	account = &model.ThirdAccount{
		UserId:     sql.NullInt64{},
		Avatar:     userInfo.FigureurlQQ1,
		Nickname:   strings.TrimSpace(userInfo.Nickname),
		ThirdType:  constants.ThirdAccountTypeQQ,
		ThirdId:    userInfo.Unionid,
		ExtraData:  userInfoJson,
		CreateTime: dates.NowTimestamp(),
		UpdateTime: dates.NowTimestamp(),
	}
	err = s.Create(account)
	if err != nil {
		return nil, err
	}
	return account, nil
}
