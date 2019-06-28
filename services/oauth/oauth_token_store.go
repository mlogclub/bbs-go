package oauth

import (
	"errors"
	"github.com/json-iterator/go"
	"github.com/mlogclub/mlog/model"
	"github.com/mlogclub/mlog/repositories"
	"github.com/mlogclub/simple"
	"gopkg.in/oauth2.v3"
	"gopkg.in/oauth2.v3/models"
)

// TokenStore实现
type OauthTokenStore struct {
	OauthTokenRepository *repositories.OauthTokenRepository
}

// 存储oauth token
func NewOauthTokenStore() *OauthTokenStore {
	return &OauthTokenStore{OauthTokenRepository: repositories.NewOauthTokenRepository()}
}

func (s *OauthTokenStore) Create(info oauth2.TokenInfo) error {
	buf, _ := jsoniter.Marshal(info)
	item := &model.OauthToken{
		Data: string(buf),
	}

	if code := info.GetCode(); code != "" {
		item.Code = code
		item.ExpiredAt = info.GetCodeCreateAt().Add(info.GetCodeExpiresIn()).Unix()
	} else {
		item.AccessToken = info.GetAccess()
		item.ExpiredAt = info.GetAccessCreateAt().Add(info.GetAccessExpiresIn()).Unix()
		if refresh := info.GetRefresh(); refresh != "" {
			item.RefreshToken = info.GetRefresh()
			item.ExpiredAt = info.GetRefreshCreateAt().Add(info.GetRefreshExpiresIn()).Unix()
		}
	}
	return s.OauthTokenRepository.Create(simple.GetDB(), item)
}

func (s *OauthTokenStore) RemoveByCode(code string) error {
	s.OauthTokenRepository.RemoveByCode(simple.GetDB(), code)
	return nil
}

func (s *OauthTokenStore) RemoveByAccess(access string) error {
	s.OauthTokenRepository.RemoveByAccessToken(simple.GetDB(), access)
	return nil
}

func (s *OauthTokenStore) RemoveByRefresh(refresh string) error {
	s.OauthTokenRepository.RemoveByRefreshToken(simple.GetDB(), refresh)
	return nil
}

func (s *OauthTokenStore) GetByCode(code string) (oauth2.TokenInfo, error) {
	if len(code) == 0 {
		return nil, nil
	}
	oauthToken := s.OauthTokenRepository.GetByCode(simple.GetDB(), code)
	if oauthToken == nil {
		return nil, errors.New("invalidate code")
	}
	return s.toTokenInfo(oauthToken.Data)
}

func (s *OauthTokenStore) GetByAccess(access string) (oauth2.TokenInfo, error) {
	if len(access) == 0 {
		return nil, nil
	}
	oauthToken := s.OauthTokenRepository.GetByAccessToken(simple.GetDB(), access)
	if oauthToken == nil {
		return nil, errors.New("invalidate access token")
	}
	return s.toTokenInfo(oauthToken.Data)
}

func (s *OauthTokenStore) GetByRefresh(refresh string) (oauth2.TokenInfo, error) {
	if len(refresh) == 0 {
		return nil, nil
	}
	oauthToken := s.OauthTokenRepository.GetByRefreshToken(simple.GetDB(), refresh)
	if oauthToken == nil {
		return nil, errors.New("invalidate refresh token")
	}
	return s.toTokenInfo(oauthToken.Data)
}

func (s *OauthTokenStore) toTokenInfo(data string) (oauth2.TokenInfo, error) {
	var tm models.Token
	err := jsoniter.Unmarshal([]byte(data), &tm)
	if err != nil {
		return nil, err
	}
	return &tm, nil
}
