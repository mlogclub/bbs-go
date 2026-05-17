package services

import (
	"bbs-go/internal/models"
	"bbs-go/internal/models/constants"
	"bbs-go/internal/models/req"
	"bbs-go/internal/pkg/locales"
	"bbs-go/internal/repositories"
	"errors"
	"strconv"
	"strings"

	"github.com/mlogclub/simple/common/dates"
	"github.com/mlogclub/simple/common/strs"
	"github.com/mlogclub/simple/sqls"
	"github.com/mlogclub/simple/web/params"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var VoteService = newVoteService()

func newVoteService() *voteService {
	return &voteService{}
}

type voteService struct {
}

func (s *voteService) Get(id int64) *models.Vote {
	if id <= 0 {
		return nil
	}
	return repositories.VoteRepository.Get(sqls.DB(), id)
}

func (s *voteService) Take(where ...interface{}) *models.Vote {
	return repositories.VoteRepository.Take(sqls.DB(), where...)
}

func (s *voteService) Find(cnd *sqls.Cnd) []models.Vote {
	return repositories.VoteRepository.Find(sqls.DB(), cnd)
}

func (s *voteService) FindOne(cnd *sqls.Cnd) *models.Vote {
	return repositories.VoteRepository.FindOne(sqls.DB(), cnd)
}

func (s *voteService) FindPageByParams(params *params.QueryParams) (list []models.Vote, paging *sqls.Paging) {
	return repositories.VoteRepository.FindPageByParams(sqls.DB(), params)
}

func (s *voteService) FindPageByCnd(cnd *sqls.Cnd) (list []models.Vote, paging *sqls.Paging) {
	return repositories.VoteRepository.FindPageByCnd(sqls.DB(), cnd)
}

func (s *voteService) Count(cnd *sqls.Cnd) int64 {
	return repositories.VoteRepository.Count(sqls.DB(), cnd)
}

func (s *voteService) Create(t *models.Vote) error {
	return repositories.VoteRepository.Create(sqls.DB(), t)
}

func (s *voteService) Update(t *models.Vote) error {
	return repositories.VoteRepository.Update(sqls.DB(), t)
}

func (s *voteService) Updates(id int64, columns map[string]interface{}) error {
	return repositories.VoteRepository.Updates(sqls.DB(), id, columns)
}

func (s *voteService) UpdateColumn(id int64, name string, value interface{}) error {
	return repositories.VoteRepository.UpdateColumn(sqls.DB(), id, name, value)
}

func (s *voteService) Delete(id int64) {
	repositories.VoteRepository.Delete(sqls.DB(), id)
}

func (s *voteService) CheckCreateForm(form *req.CreateVoteForm) error {
	if form == nil {
		return nil
	}
	form.Title = strings.TrimSpace(form.Title)
	if strs.IsBlank(form.Title) {
		return errors.New(locales.Get("vote.title_required"))
	}
	if strs.RuneLen(form.Title) > 128 {
		return errors.New(locales.Get("vote.title_too_long"))
	}
	if len(form.Options) < 2 {
		return errors.New(locales.Get("vote.option_min"))
	}
	if len(form.Options) > 20 {
		return errors.New(locales.Get("vote.option_max"))
	}

	for i := range form.Options {
		form.Options[i].Content = strings.TrimSpace(form.Options[i].Content)
		if strs.IsBlank(form.Options[i].Content) {
			return errors.New(locales.Get("vote.option_required"))
		}
		if strs.RuneLen(form.Options[i].Content) > 256 {
			return errors.New(locales.Get("vote.option_too_long"))
		}
	}

	if form.ExpiredAt <= dates.NowTimestamp() {
		return errors.New(locales.Get("vote.expired_at_future"))
	}

	switch form.Type {
	case constants.VoteTypeSingle:
		form.VoteNum = 1
	case constants.VoteTypeMultiple:
		if form.VoteNum <= 0 {
			return errors.New(locales.Get("vote.multiple_num_required"))
		}
		if form.VoteNum > len(form.Options) {
			return errors.New(locales.Get("vote.multiple_num_too_large"))
		}
	default:
		return errors.New(locales.Get("vote.type_not_supported"))
	}
	return nil
}

func (s *voteService) CreateWithOptionsTx(ctx *sqls.TxContext, topicId, userId int64, form *req.CreateVoteForm, now int64) (*models.Vote, error) {
	if form == nil {
		return nil, nil
	}
	if err := s.CheckCreateForm(form); err != nil {
		return nil, err
	}

	vote := &models.Vote{
		Type:        form.Type,
		Title:       form.Title,
		ExpiredAt:   form.ExpiredAt,
		TopicId:     topicId,
		UserId:      userId,
		VoteNum:     form.VoteNum,
		OptionCount: len(form.Options),
		CreateTime:  now,
	}
	if err := repositories.VoteRepository.Create(ctx.Tx, vote); err != nil {
		return nil, err
	}

	for i, option := range form.Options {
		item := &models.VoteOption{
			VoteId:     vote.Id,
			Content:    option.Content,
			SortNo:     i + 1,
			CreateTime: now,
		}
		if err := repositories.VoteOptionRepository.Create(ctx.Tx, item); err != nil {
			return nil, err
		}
	}

	return vote, nil
}

func (s *voteService) Cast(userId int64, form req.VoteCastForm) error {
	if form.VoteId <= 0 {
		return errors.New(locales.Get("vote.vote_id_required"))
	}
	if len(form.OptionIds) == 0 {
		return errors.New(locales.Get("vote.select_option_required"))
	}

	selected := make([]int64, 0, len(form.OptionIds))
	selectedSet := make(map[int64]bool, len(form.OptionIds))
	for _, optionId := range form.OptionIds {
		if optionId <= 0 {
			return errors.New(locales.Get("vote.option_invalid"))
		}
		if selectedSet[optionId] {
			continue
		}
		selectedSet[optionId] = true
		selected = append(selected, optionId)
	}
	if len(selected) == 0 {
		return errors.New(locales.Get("vote.select_option_required"))
	}

	return sqls.WithTransaction(func(ctx *sqls.TxContext) error {
		tx := ctx.Tx
		vote := &models.Vote{}
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(vote, "id = ?", form.VoteId).Error; err != nil {
			return errors.New(locales.Get("vote.not_found"))
		}
		if dates.NowTimestamp() > vote.ExpiredAt {
			return errors.New(locales.Get("vote.expired"))
		}

		exists := repositories.VoteRecordRepository.Take(tx, "user_id = ? AND vote_id = ?", userId, vote.Id)
		if exists != nil {
			return errors.New(locales.Get("vote.already_voted"))
		}

		options := repositories.VoteOptionRepository.Find(tx,
			sqls.NewCnd().Eq("vote_id", vote.Id).Asc("sort_no").Asc("id"),
		)
		if len(options) == 0 {
			return errors.New(locales.Get("vote.option_not_found"))
		}

		optionMap := make(map[int64]bool, len(options))
		for _, option := range options {
			optionMap[option.Id] = true
		}
		for _, optionId := range selected {
			if !optionMap[optionId] {
				return errors.New(locales.Get("vote.option_invalid"))
			}
		}

		switch vote.Type {
		case constants.VoteTypeSingle:
			if len(selected) != 1 {
				return errors.New(locales.Get("vote.single_only_one"))
			}
		case constants.VoteTypeMultiple:
			if len(selected) > vote.VoteNum {
				return errors.New(locales.Get("vote.multiple_over_limit"))
			}
		default:
			return errors.New(locales.Get("vote.type_not_supported"))
		}

		now := dates.NowTimestamp()
		record := &models.VoteRecord{
			UserId:     userId,
			VoteId:     vote.Id,
			OptionIds:  s.JoinOptionIds(selected),
			CreateTime: now,
		}
		if err := repositories.VoteRecordRepository.Create(tx, record); err != nil {
			return err
		}

		if err := tx.Model(&models.Vote{}).Where("id = ?", vote.Id).
			UpdateColumn("vote_count", gorm.Expr("vote_count + 1")).Error; err != nil {
			return err
		}

		for _, optionId := range selected {
			if err := tx.Model(&models.VoteOption{}).Where("id = ?", optionId).
				UpdateColumn("vote_count", gorm.Expr("vote_count + 1")).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *voteService) JoinOptionIds(optionIds []int64) string {
	if len(optionIds) == 0 {
		return ""
	}
	parts := make([]string, 0, len(optionIds))
	for _, optionId := range optionIds {
		parts = append(parts, strconv.FormatInt(optionId, 10))
	}
	return strings.Join(parts, ",")
}

func (s *voteService) ParseOptionIds(optionIds string) []int64 {
	if strs.IsBlank(optionIds) {
		return nil
	}
	arr := strings.Split(optionIds, ",")
	ret := make([]int64, 0, len(arr))
	for _, item := range arr {
		id, err := strconv.ParseInt(strings.TrimSpace(item), 10, 64)
		if err == nil && id > 0 {
			ret = append(ret, id)
		}
	}
	return ret
}
