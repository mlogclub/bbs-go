package services

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/mlogclub/simple"
	"github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"

	"bbs-go/common"
	"bbs-go/common/avatar"
	"bbs-go/common/uploader"
	"bbs-go/services/cache"

	"bbs-go/model"
	"bbs-go/repositories"
)

type ScanUserCallback func(users []model.User)

var UserService = newUserService()

func newUserService() *userService {
	return &userService{}
}

type userService struct {
}

func (s *userService) Get(id int64) *model.User {
	return repositories.UserRepository.Get(simple.DB(), id)
}

func (s *userService) Take(where ...interface{}) *model.User {
	return repositories.UserRepository.Take(simple.DB(), where...)
}

func (s *userService) Find(cnd *simple.SqlCnd) []model.User {
	return repositories.UserRepository.Find(simple.DB(), cnd)
}

func (s *userService) FindOne(cnd *simple.SqlCnd) *model.User {
	return repositories.UserRepository.FindOne(simple.DB(), cnd)
}

func (s *userService) FindPageByParams(params *simple.QueryParams) (list []model.User, paging *simple.Paging) {
	return repositories.UserRepository.FindPageByParams(simple.DB(), params)
}

func (s *userService) FindPageByCnd(cnd *simple.SqlCnd) (list []model.User, paging *simple.Paging) {
	return repositories.UserRepository.FindPageByCnd(simple.DB(), cnd)
}

func (s *userService) Create(t *model.User) error {
	err := repositories.UserRepository.Create(simple.DB(), t)
	if err == nil {
		cache.UserCache.Invalidate(t.Id)
	}
	return nil
}

func (s *userService) Update(t *model.User) error {
	err := repositories.UserRepository.Update(simple.DB(), t)
	cache.UserCache.Invalidate(t.Id)
	return err
}

func (s *userService) Updates(id int64, columns map[string]interface{}) error {
	err := repositories.UserRepository.Updates(simple.DB(), id, columns)
	cache.UserCache.Invalidate(id)
	return err
}

func (s *userService) UpdateColumn(id int64, name string, value interface{}) error {
	err := repositories.UserRepository.UpdateColumn(simple.DB(), id, name, value)
	cache.UserCache.Invalidate(id)
	return err
}

func (s *userService) Delete(id int64) {
	repositories.UserRepository.Delete(simple.DB(), id)
	cache.UserCache.Invalidate(id)
}

// Scan 扫描
func (s *userService) Scan(cb ScanUserCallback) {
	var cursor int64
	for {
		list := repositories.UserRepository.Find(simple.DB(), simple.NewSqlCnd().Where("id > ?", cursor).Asc("id").Limit(100))
		if list == nil || len(list) == 0 {
			break
		}
		cursor = list[len(list)-1].Id
		cb(list)
	}
}

// GetByEmail 根据邮箱查找
func (s *userService) GetByEmail(email string) *model.User {
	return repositories.UserRepository.GetByEmail(simple.DB(), email)
}

// GetByUsername 根据用户名查找
func (s *userService) GetByUsername(username string) *model.User {
	return repositories.UserRepository.GetByUsername(simple.DB(), username)
}

// SignUp 注册
func (s *userService) SignUp(username, email, nickname, password, rePassword string) (*model.User, error) {
	username = strings.TrimSpace(username)
	email = strings.TrimSpace(email)
	nickname = strings.TrimSpace(nickname)

	// 验证昵称
	if len(nickname) == 0 {
		return nil, errors.New("昵称不能为空")
	}

	// 验证密码
	err := common.IsValidatePassword(password, rePassword)
	if err != nil {
		return nil, err
	}

	// 验证邮箱
	if len(email) > 0 {
		if err := common.IsValidateEmail(email); err != nil {
			return nil, err
		}
		if s.GetByEmail(email) != nil {
			return nil, errors.New("邮箱：" + email + " 已被占用")
		}
	} else {
		return nil, errors.New("请输入邮箱")
	}

	// 验证用户名
	if len(username) > 0 {
		if err := common.IsValidateUsername(username); err != nil {
			return nil, err
		}
		if s.isUsernameExists(username) {
			return nil, errors.New("用户名：" + username + " 已被占用")
		}
	}

	user := &model.User{
		Username:   simple.SqlNullString(username),
		Email:      simple.SqlNullString(email),
		Nickname:   nickname,
		Password:   simple.EncodePassword(password),
		Status:     model.StatusOk,
		CreateTime: simple.NowTimestamp(),
		UpdateTime: simple.NowTimestamp(),
	}

	err = simple.Tx(simple.DB(), func(tx *gorm.DB) error {
		if err := repositories.UserRepository.Create(tx, user); err != nil {
			return err
		}

		avatarUrl, err := s.HandleAvatar(user.Id, "")
		if err != nil {
			return err
		}

		if err := repositories.UserRepository.UpdateColumn(tx, user.Id, "avatar", avatarUrl); err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return user, nil
}

// SignIn 登录
func (s *userService) SignIn(username, password string) (*model.User, error) {
	if len(username) == 0 {
		return nil, errors.New("用户名/邮箱不能为空")
	}
	if len(password) == 0 {
		return nil, errors.New("密码不能为空")
	}
	var user *model.User = nil
	if err := common.IsValidateEmail(username); err == nil { // 如果用户输入的是邮箱
		user = s.GetByEmail(username)
	} else {
		user = s.GetByUsername(username)
	}
	if user == nil || user.Status != model.StatusOk {
		return nil, errors.New("用户不存在或被禁用")
	}
	if !simple.ValidatePassword(user.Password, password) {
		return nil, errors.New("密码错误")
	}
	return user, nil
}

// SignInByThirdAccount 第三方账号登录
func (s *userService) SignInByThirdAccount(thirdAccount *model.ThirdAccount) (*model.User, *simple.CodeError) {
	user := s.Get(thirdAccount.UserId.Int64)
	if user != nil {
		if user.Status != model.StatusOk {
			return nil, simple.NewErrorMsg("用户已被禁用")
		}
		return user, nil
	}

	var homePage string
	var description string
	if thirdAccount.ThirdType == model.ThirdAccountTypeGithub {
		if blog := gjson.Get(thirdAccount.ExtraData, "blog"); blog.Exists() && len(blog.String()) > 0 {
			homePage = blog.String()
		} else if htmlUrl := gjson.Get(thirdAccount.ExtraData, "html_url"); htmlUrl.Exists() && len(htmlUrl.String()) > 0 {
			homePage = htmlUrl.String()
		}

		description = gjson.Get(thirdAccount.ExtraData, "bio").String()
	}

	user = &model.User{
		Username:    sql.NullString{},
		Nickname:    thirdAccount.Nickname,
		Status:      model.StatusOk,
		HomePage:    homePage,
		Description: description,
		CreateTime:  simple.NowTimestamp(),
		UpdateTime:  simple.NowTimestamp(),
	}
	err := simple.Tx(simple.DB(), func(tx *gorm.DB) error {
		if err := repositories.UserRepository.Create(tx, user); err != nil {
			return err
		}

		if err := repositories.ThirdAccountRepository.UpdateColumn(tx, thirdAccount.Id, "user_id", user.Id); err != nil {
			return err
		}

		avatarUrl, err := s.HandleAvatar(user.Id, thirdAccount.Avatar)
		if err != nil {
			return err
		}

		if err := repositories.UserRepository.UpdateColumn(tx, user.Id, "avatar", avatarUrl); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, simple.FromError(err)
	}
	cache.UserCache.Invalidate(user.Id)
	return user, nil
}

// HandleAvatar 处理头像，优先级如下：1. 如果第三方登录带有来头像；2. 生成随机默认头像
// thirdAvatar: 第三方登录带过来的头像
func (s *userService) HandleAvatar(userId int64, thirdAvatar string) (string, error) {
	if len(thirdAvatar) > 0 {
		return uploader.CopyImage(thirdAvatar)
	}

	avatarBytes, err := avatar.Generate(userId)
	if err != nil {
		return "", err
	}
	return uploader.PutImage(avatarBytes)
}

// isEmailExists 邮箱是否存在
func (s *userService) isEmailExists(email string) bool {
	if len(email) == 0 { // 如果邮箱为空，那么就认为是不存在
		return false
	}
	return s.GetByEmail(email) != nil
}

// isUsernameExists 用户名是否存在
func (s *userService) isUsernameExists(username string) bool {
	return s.GetByUsername(username) != nil
}

// SetAvatar 更新头像
func (s *userService) UpdateAvatar(userId int64, avatar string) error {
	return s.UpdateColumn(userId, "avatar", avatar)
}

// SetUsername 设置用户名
func (s *userService) SetUsername(userId int64, username string) error {
	username = strings.TrimSpace(username)
	if err := common.IsValidateUsername(username); err != nil {
		return err
	}

	user := s.Get(userId)
	if len(user.Username.String) > 0 {
		return errors.New("你已设置了用户名，无法重复设置。")
	}
	if s.isUsernameExists(username) {
		return errors.New("用户名：" + username + " 已被占用")
	}
	return s.UpdateColumn(userId, "username", username)
}

// SetEmail 设置密码
func (s *userService) SetEmail(userId int64, email string) error {
	email = strings.TrimSpace(email)
	if err := common.IsValidateEmail(email); err != nil {
		return err
	}
	if s.isEmailExists(email) {
		return errors.New("邮箱：" + email + " 已被占用")
	}
	return s.UpdateColumn(userId, "email", email)
}

// SetPassword 设置密码
func (s *userService) SetPassword(userId int64, password, rePassword string) error {
	if err := common.IsValidatePassword(password, rePassword); err != nil {
		return err
	}
	user := s.Get(userId)
	if len(user.Password) > 0 {
		return errors.New("你已设置了密码，如需修改请前往修改页面。")
	}
	password = simple.EncodePassword(password)
	return s.UpdateColumn(userId, "password", password)
}

// UpdatePassword 修改密码
func (s *userService) UpdatePassword(userId int64, oldPassword, password, rePassword string) error {
	if err := common.IsValidatePassword(password, rePassword); err != nil {
		return err
	}
	user := s.Get(userId)

	if len(user.Password) == 0 {
		return errors.New("你没设置密码，请先设置密码")
	}

	if !simple.ValidatePassword(user.Password, oldPassword) {
		return errors.New("旧密码验证失败")
	}

	return s.UpdateColumn(userId, "password", simple.EncodePassword(password))
}

// IncrTopicCount topic_count + 1
func (s *userService) IncrTopicCount(userId int64) int {
	t := repositories.UserRepository.Get(simple.DB(), userId)
	if t == nil {
		return 0
	}
	topicCount := t.TopicCount + 1
	if err := repositories.UserRepository.UpdateColumn(simple.DB(), userId, "topic_count", topicCount); err != nil {
		logrus.Error(err)
	} else {
		cache.UserCache.Invalidate(userId)
	}
	return topicCount
}

// IncrCommentCount comment_count + 1
func (s *userService) IncrCommentCount(userId int64) int {
	t := repositories.UserRepository.Get(simple.DB(), userId)
	if t == nil {
		return 0
	}
	commentCount := t.CommentCount + 1
	if err := repositories.UserRepository.UpdateColumn(simple.DB(), userId, "comment_count", commentCount); err != nil {
		logrus.Error(err)
	} else {
		cache.UserCache.Invalidate(userId)
	}
	return commentCount
}

// SyncUserCount 同步用户计数
func (s *userService) SyncUserCount() {
	s.Scan(func(users []model.User) {
		for _, user := range users {
			topicCount := repositories.TopicRepository.Count(simple.DB(), simple.NewSqlCnd().Eq("user_id", user.Id).Eq("status", model.StatusOk))
			commentCount := repositories.CommentRepository.Count(simple.DB(), simple.NewSqlCnd().Eq("user_id", user.Id).Eq("status", model.StatusOk))
			_ = repositories.UserRepository.UpdateColumn(simple.DB(), user.Id, "topic_count", topicCount)
			_ = repositories.UserRepository.UpdateColumn(simple.DB(), user.Id, "comment_count", commentCount)
			cache.UserCache.Invalidate(user.Id)
		}
	})
}
