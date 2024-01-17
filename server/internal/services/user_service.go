package services

import (
	"bbs-go/internal/models/constants"
	"bbs-go/internal/pkg/bbsurls"
	"bbs-go/internal/pkg/email"
	"bbs-go/internal/pkg/errs"
	"bbs-go/internal/pkg/validate"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/mlogclub/simple/common/dates"
	"github.com/mlogclub/simple/common/passwd"
	"github.com/mlogclub/simple/common/strs"
	"github.com/mlogclub/simple/sqls"
	"github.com/mlogclub/simple/web"
	"github.com/mlogclub/simple/web/params"
	"gorm.io/gorm"

	"bbs-go/internal/cache"

	"bbs-go/internal/models"
	"bbs-go/internal/repositories"
)

// 邮箱验证邮件有效期（小时）
const emailVerifyExpireHour = 24

var UserService = newUserService()

func newUserService() *userService {
	return &userService{}
}

type userService struct {
}

func (s *userService) Get(id int64) *models.User {
	return repositories.UserRepository.Get(sqls.DB(), id)
}

func (s *userService) Take(where ...interface{}) *models.User {
	return repositories.UserRepository.Take(sqls.DB(), where...)
}

func (s *userService) Find(cnd *sqls.Cnd) []models.User {
	return repositories.UserRepository.Find(sqls.DB(), cnd)
}

func (s *userService) FindOne(cnd *sqls.Cnd) *models.User {
	return repositories.UserRepository.FindOne(sqls.DB(), cnd)
}

func (s *userService) FindPageByParams(params *params.QueryParams) (list []models.User, paging *sqls.Paging) {
	return repositories.UserRepository.FindPageByParams(sqls.DB(), params)
}

func (s *userService) FindPageByCnd(cnd *sqls.Cnd) (list []models.User, paging *sqls.Paging) {
	return repositories.UserRepository.FindPageByCnd(sqls.DB(), cnd)
}

func (s *userService) Create(t *models.User) error {
	err := repositories.UserRepository.Create(sqls.DB(), t)
	if err == nil {
		cache.UserCache.Invalidate(t.Id)
	}
	return nil
}

func (s *userService) Update(t *models.User) error {
	err := repositories.UserRepository.Update(sqls.DB(), t)
	cache.UserCache.Invalidate(t.Id)
	return err
}

func (s *userService) Updates(id int64, columns map[string]interface{}) error {
	err := repositories.UserRepository.Updates(sqls.DB(), id, columns)
	cache.UserCache.Invalidate(id)
	return err
}

func (s *userService) UpdateColumn(id int64, name string, value interface{}) error {
	err := repositories.UserRepository.UpdateColumn(sqls.DB(), id, name, value)
	cache.UserCache.Invalidate(id)
	return err
}

func (s *userService) Delete(id int64) {
	repositories.UserRepository.Delete(sqls.DB(), id)
	cache.UserCache.Invalidate(id)
}

// Scan 扫描
func (s *userService) Scan(callback func(users []models.User)) {
	var cursor int64
	for {
		list := repositories.UserRepository.Find(sqls.DB(), sqls.NewCnd().Where("id > ?", cursor).Asc("id").Limit(100))
		if len(list) == 0 {
			break
		}
		cursor = list[len(list)-1].Id
		callback(list)
	}
}

// Forbidden 禁言
func (s *userService) Forbidden(operatorId, userId int64, days int, reason string, r *http.Request) error {
	var forbiddenEndTime int64
	if days == -1 { // 永久禁言
		forbiddenEndTime = -1
	} else if days > 0 {
		forbiddenEndTime = dates.Timestamp(time.Now().Add(time.Hour * 24 * time.Duration(days)))
	} else {
		return errors.New("禁言时间错误")
	}
	if repositories.UserRepository.UpdateColumn(sqls.DB(), userId, "forbidden_end_time", forbiddenEndTime) == nil {
		cache.UserCache.Invalidate(userId)
		description := ""
		if strs.IsNotBlank(reason) {
			description = "禁言原因：" + reason
		}
		OperateLogService.AddOperateLog(operatorId, constants.OpTypeForbidden, constants.EntityUser, userId,
			description, r)

		// 永久禁言
		if days == -1 {
			user := cache.UserCache.Get(userId)
			_ = s.DecrScore(userId, user.Score, constants.EntityUser, strconv.FormatInt(operatorId, 10), "永久禁言")
			go func() {
				// 删除话题
				TopicService.ScanByUser(userId, func(topics []models.Topic) {
					for _, topic := range topics {
						if topic.Status != constants.StatusDeleted {
							_ = TopicService.Delete(topic.Id, operatorId, nil)
						}
					}
				})

				// 删除文章
				ArticleService.ScanByUser(userId, func(articles []models.Article) {
					for _, article := range articles {
						if article.Status != constants.StatusDeleted {
							_ = ArticleService.Delete(article.Id)
						}
					}
				})

				// 删除评论
				CommentService.ScanByUser(userId, func(comments []models.Comment) {
					for _, comment := range comments {
						if comment.Status != constants.StatusDeleted {
							_ = CommentService.Delete(comment.Id)
						}
					}
				})

			}()
		}
	}
	return nil
}

// RemoveForbidden 移除禁言
func (s *userService) RemoveForbidden(operatorId, userId int64, r *http.Request) {
	user := s.Get(userId)
	if user == nil || !user.IsForbidden() {
		return
	}
	if repositories.UserRepository.UpdateColumn(sqls.DB(), userId, "forbidden_end_time", 0) == nil {
		cache.UserCache.Invalidate(user.Id)
		OperateLogService.AddOperateLog(operatorId, constants.OpTypeRemoveForbidden, constants.EntityUser, userId, "", r)
	}
}

// GetByEmail 根据邮箱查找
func (s *userService) GetByEmail(email string) *models.User {
	return repositories.UserRepository.GetByEmail(sqls.DB(), email)
}

// GetByUsername 根据用户名查找
func (s *userService) GetByUsername(username string) *models.User {
	return repositories.UserRepository.GetByUsername(sqls.DB(), username)
}

// SignUp 注册
func (s *userService) SignUp(username, email, nickname, password, rePassword string) (*models.User, error) {
	username = strings.TrimSpace(username)
	email = strings.TrimSpace(email)
	nickname = strings.TrimSpace(nickname)

	// 验证昵称
	if len(nickname) == 0 {
		return nil, errors.New("昵称不能为空")
	}

	// 验证密码
	err := validate.IsValidPassword(password, rePassword)
	if err != nil {
		return nil, err
	}

	// 验证邮箱
	if len(email) > 0 {
		if err := validate.IsEmail(email); err != nil {
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
		if err := validate.IsUsername(username); err != nil {
			return nil, err
		}
		if s.isUsernameExists(username) {
			return nil, errors.New("用户名：" + username + " 已被占用")
		}
	}

	user := &models.User{
		Username:   sqls.SqlNullString(username),
		Email:      sqls.SqlNullString(email),
		Nickname:   nickname,
		Password:   passwd.EncodePassword(password),
		Status:     constants.StatusOk,
		CreateTime: dates.NowTimestamp(),
		UpdateTime: dates.NowTimestamp(),
	}

	err = repositories.UserRepository.Create(sqls.DB(), user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// SignIn 登录
func (s *userService) SignIn(username, password string) (*models.User, error) {
	if strs.IsBlank(username) {
		return nil, errors.New("用户名/邮箱不能为空")
	}
	if strs.IsBlank(password) {
		return nil, errors.New("密码不能为空")
	}
	if err := validate.IsPassword(password); err != nil {
		return nil, err
	}
	var user *models.User = nil
	if err := validate.IsEmail(username); err == nil { // 如果用户输入的是邮箱
		user = s.GetByEmail(username)
	} else {
		user = s.GetByUsername(username)
	}
	if user == nil || user.Status != constants.StatusOk {
		return nil, errors.New("用户名或密码错误")
	}
	if !passwd.ValidatePassword(user.Password, password) {
		return nil, errors.New("用户名或密码错误")
	}
	return user, nil
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

// UpdateAvatar 更新头像
func (s *userService) UpdateAvatar(userId int64, avatar string) error {
	return s.UpdateColumn(userId, "avatar", avatar)
}

// UpdateNickname 更新昵称
func (s *userService) UpdateNickname(userId int64, nickname string) error {
	return s.UpdateColumn(userId, "nickname", nickname)
}

// UpdateDescription 更新简介
func (s *userService) UpdateDescription(userId int64, description string) error {
	return s.UpdateColumn(userId, "description", description)
}

// UpdateGender 修改性别
func (s *userService) UpdateGender(userId int64, gender string) error {
	if strs.IsBlank(gender) {
		return s.UpdateColumn(userId, "gender", "")
	} else {
		if gender != string(constants.GenderMale) && gender != string(constants.GenderFemale) {
			return errors.New("invalidate gender value")
		}
		return s.UpdateColumn(userId, "gender", gender)
	}
}

// UpdateBirthday 修改生日
func (s *userService) UpdateBirthday(userId int64, birthdayStr string) error {
	if strs.IsBlank(birthdayStr) {
		return s.UpdateColumn(userId, "birthday", "")
	} else {
		birthday, err := dates.Parse(birthdayStr, dates.FmtDate)
		if err != nil {
			return err
		}
		return s.UpdateColumn(userId, "birthday", birthday)
	}
}

// UpdateBackgroundImage 修改背景图
func (s *userService) UpdateBackgroundImage(userId int64, backgroundImage string) error {
	return s.UpdateColumn(userId, "background_image", backgroundImage)
}

// SetUsername 设置用户名
func (s *userService) SetUsername(userId int64, username string) error {
	username = strings.TrimSpace(username)
	if err := validate.IsUsername(username); err != nil {
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
	if err := validate.IsEmail(email); err != nil {
		return err
	}
	user := s.Get(userId)
	if user == nil {
		return errors.New("用户不存在")
	}
	if user.Email.String == email {
		// 用户邮箱没做变更
		return nil
	}
	if s.isEmailExists(email) {
		return errors.New("邮箱：" + email + " 已被占用")
	}
	return s.Updates(userId, map[string]interface{}{
		"email":          email,
		"email_verified": false,
	})
}

// SetPassword 设置密码
func (s *userService) SetPassword(userId int64, password, rePassword string) error {
	if err := validate.IsValidPassword(password, rePassword); err != nil {
		return err
	}
	user := s.Get(userId)
	if len(user.Password) > 0 {
		return errors.New("你已设置了密码，如需修改请前往修改页面。")
	}
	password = passwd.EncodePassword(password)
	return s.UpdateColumn(userId, "password", password)
}

// UpdatePassword 修改密码
func (s *userService) UpdatePassword(userId int64, oldPassword, password, rePassword string) error {
	if err := validate.IsValidPassword(password, rePassword); err != nil {
		return err
	}
	user := s.Get(userId)

	if len(user.Password) == 0 {
		return errors.New("你没设置密码，请先设置密码")
	}

	if !passwd.ValidatePassword(user.Password, oldPassword) {
		return errors.New("旧密码验证失败")
	}

	return s.UpdateColumn(userId, "password", passwd.EncodePassword(password))
}

// IncrTopicCount topic_count + 1
func (s *userService) IncrTopicCount(tx *gorm.DB, userId int64) error {
	if err := repositories.UserRepository.UpdateColumn(sqls.DB(), userId, "topic_count", gorm.Expr("topic_count + 1")); err != nil {
		slog.Error(err.Error(), slog.Any("err", err))
		return err
	}
	cache.UserCache.Invalidate(userId)
	return nil
}

// IncrCommentCount comment_count + 1
func (s *userService) IncrCommentCount(userId int64) int {
	t := repositories.UserRepository.Get(sqls.DB(), userId)
	if t == nil {
		return 0
	}
	commentCount := t.CommentCount + 1
	if err := repositories.UserRepository.UpdateColumn(sqls.DB(), userId, "comment_count", commentCount); err != nil {
		slog.Error(err.Error(), slog.Any("err", err))
	} else {
		cache.UserCache.Invalidate(userId)
	}
	return commentCount
}

// SyncUserCount 同步用户计数
func (s *userService) SyncUserCount() {
	s.Scan(func(users []models.User) {
		for _, user := range users {
			topicCount := repositories.TopicRepository.Count(sqls.DB(), sqls.NewCnd().Eq("user_id", user.Id).Eq("status", constants.StatusOk))
			commentCount := repositories.CommentRepository.Count(sqls.DB(), sqls.NewCnd().Eq("user_id", user.Id).Eq("status", constants.StatusOk))
			_ = repositories.UserRepository.UpdateColumn(sqls.DB(), user.Id, "topic_count", topicCount)
			_ = repositories.UserRepository.UpdateColumn(sqls.DB(), user.Id, "comment_count", commentCount)
			cache.UserCache.Invalidate(user.Id)
		}
	})
}

// SendEmailVerifyEmail 发送邮箱验证邮件
func (s *userService) SendEmailVerifyEmail(userId int64) error {
	user := s.Get(userId)
	if user == nil {
		return errors.New("用户不存在")
	}
	if user.EmailVerified {
		return errors.New("用户邮箱已验证")
	}
	if err := validate.IsEmail(user.Email.String); err != nil {
		return err
	}
	// 如果设置了邮箱白名单
	if emailWhitelist := SysConfigService.GetEmailWhitelist(); len(emailWhitelist) > 0 {
		isInWhitelist := false
		for _, whitelist := range emailWhitelist {
			if strings.Contains(strings.ToLower(user.Email.String), strings.ToLower(whitelist)) {
				isInWhitelist = true
				break
			}
		}
		if !isInWhitelist {
			// 直接返回，也不抛出异常了，就是不发邮件
			slog.Error("不支持使用该邮箱进行验证.", slog.String("email", user.Email.String))
			return errors.New("不支持该类型邮箱")
		}
	}
	var (
		token     = strs.UUID()
		url       = bbsurls.AbsUrl("/user/email/verify?token=" + token)
		link      = &models.ActionLink{Title: "点击这里验证邮箱>>", Url: url}
		siteTitle = cache.SysConfigCache.GetValue(constants.SysConfigSiteTitle)
		subject   = "邮箱验证 - " + siteTitle
		title     = "邮箱验证 - " + siteTitle
		content   = "该邮件用于验证你在 " + siteTitle + " 中设置邮箱的正确性，请在" + strconv.Itoa(emailVerifyExpireHour) + "小时内完成验证。验证链接：" + url
	)
	return sqls.DB().Transaction(func(tx *gorm.DB) error {
		if err := repositories.EmailCodeRepository.Create(tx, &models.EmailCode{
			Model:      models.Model{},
			UserId:     userId,
			Email:      user.Email.String,
			Code:       "",
			Token:      token,
			Title:      title,
			Content:    content,
			Used:       false,
			CreateTime: dates.NowTimestamp(),
		}); err != nil {
			return nil
		}
		if err := email.SendTemplateEmail(nil, user.Email.String, subject, title, content, "", link); err != nil {
			return err
		}
		return nil
	})
}

// VerifyEmail 验证邮箱
func (s *userService) VerifyEmail(token string) (string, error) {
	emailCode := EmailCodeService.FindOne(sqls.NewCnd().Eq("token", token))
	if emailCode == nil || emailCode.Used {
		return "", errors.New("非法请求")
	}

	user := s.Get(emailCode.UserId)
	if user == nil || emailCode.Email != user.Email.String {
		return "", errors.New("验证码过期")
	}
	if dates.FromTimestamp(emailCode.CreateTime).Add(time.Hour * time.Duration(emailVerifyExpireHour)).Before(time.Now()) {
		return "", errors.New("验证邮件已过期")
	}
	err := sqls.DB().Transaction(func(tx *gorm.DB) error {
		if err := repositories.UserRepository.UpdateColumn(tx, emailCode.UserId, "email_verified", true); err != nil {
			return err
		}
		cache.UserCache.Invalidate(emailCode.UserId)
		return repositories.EmailCodeRepository.UpdateColumn(tx, emailCode.Id, "used", true)
	})
	if err != nil {
		return "", err
	}
	return emailCode.Email, nil
}

// CheckPostStatus 用于在发表内容时检查用户状态
func (s *userService) CheckPostStatus(user *models.User) error {
	if user == nil {
		return errs.NotLogin
	}
	if user.Status != constants.StatusOk {
		return errs.UserDisabled
	}
	if user.IsForbidden() {
		return errs.ForbiddenError
	}
	observeSeconds := SysConfigService.GetInt(constants.SysConfigUserObserveSeconds, 0)
	if user.InObservationPeriod(observeSeconds) {
		return web.NewError(errs.InObservationPeriod.Code, "账号尚在观察期，观察期时长："+strconv.Itoa(observeSeconds)+"秒，请稍后再试")
	}
	return nil
}

// IncrScoreForPostTopic 发帖获积分
func (s *userService) IncrScoreForPostTopic(topic *models.Topic) {
	config := SysConfigService.GetConfig()
	if config.ScoreConfig.PostTopicScore <= 0 {
		slog.Info("请配置发帖积分")
		return
	}
	err := s.addScore(topic.UserId, config.ScoreConfig.PostTopicScore, constants.EntityTopic,
		strconv.FormatInt(topic.Id, 10), "发表话题")
	if err != nil {
		slog.Error(err.Error(), slog.Any("err", err))
	}
}

// IncrScoreForPostComment 跟帖获积分
func (s *userService) IncrScoreForPostComment(comment *models.Comment) {
	// 非话题跟帖，跳过
	if comment.EntityType != constants.EntityTopic {
		return
	}
	config := SysConfigService.GetConfig()
	if config.ScoreConfig.PostCommentScore <= 0 {
		slog.Info("请配置跟帖积分")
		return
	}
	err := s.addScore(comment.UserId, config.ScoreConfig.PostCommentScore, constants.EntityComment,
		strconv.FormatInt(comment.Id, 10), "发表跟帖")
	if err != nil {
		slog.Error(err.Error(), slog.Any("err", err))
	}
}

// IncrScore 增加分数
func (s *userService) IncrScore(userId int64, score int, sourceType, sourceId, description string) error {
	if score <= 0 {
		return errors.New("分数必须为正数")
	}
	return s.addScore(userId, score, sourceType, sourceId, description)
}

// DecrScore 减少分数
func (s *userService) DecrScore(userId int64, score int, sourceType, sourceId, description string) error {
	if score <= 0 {
		return errors.New("分数必须为正数")
	}
	return s.addScore(userId, -score, sourceType, sourceId, description)
}

// addScore 加分数，也可以加负数
func (s *userService) addScore(userId int64, score int, sourceType, sourceId, description string) error {
	if score == 0 {
		return errors.New("分数不能为0")
	}
	user := s.Get(userId)
	if user == nil {
		return errors.New("用户不存在")
	}
	if err := s.Updates(userId, map[string]interface{}{
		"score":       gorm.Expr("score + ?", score),
		"update_time": dates.NowTimestamp(),
	}); err != nil {
		return err
	}

	scoreType := constants.ScoreTypeIncr
	if score < 0 {
		scoreType = constants.ScoreTypeDecr
	}
	err := UserScoreLogService.Create(&models.UserScoreLog{
		UserId:      userId,
		SourceType:  sourceType,
		SourceId:    sourceId,
		Description: description,
		Type:        scoreType,
		Score:       score,
		CreateTime:  dates.NowTimestamp(),
	})
	if err == nil {
		cache.UserCache.Invalidate(userId)
	}
	return err
}
