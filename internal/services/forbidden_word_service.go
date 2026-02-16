package services

import (
	"bbs-go/internal/cache"
	"bbs-go/internal/models"
	"bbs-go/internal/models/constants"
	"bbs-go/internal/repositories"
	"regexp"
	"strings"

	"github.com/mlogclub/simple/common/strs"
	"github.com/mlogclub/simple/sqls"
	"github.com/mlogclub/simple/web/params"
)

var ForbiddenWordService = newForbiddenWordService()

func newForbiddenWordService() *forbiddenWordService {
	return &forbiddenWordService{}
}

type forbiddenWordService struct {
}

func (s *forbiddenWordService) Get(id int64) *models.ForbiddenWord {
	return repositories.ForbiddenWordRepository.Get(sqls.DB(), id)
}

func (s *forbiddenWordService) Take(where ...interface{}) *models.ForbiddenWord {
	return repositories.ForbiddenWordRepository.Take(sqls.DB(), where...)
}

func (s *forbiddenWordService) Find(cnd *sqls.Cnd) []models.ForbiddenWord {
	return repositories.ForbiddenWordRepository.Find(sqls.DB(), cnd)
}

func (s *forbiddenWordService) FindOne(cnd *sqls.Cnd) *models.ForbiddenWord {
	return repositories.ForbiddenWordRepository.FindOne(sqls.DB(), cnd)
}

func (s *forbiddenWordService) FindPageByParams(params *params.QueryParams) (list []models.ForbiddenWord, paging *sqls.Paging) {
	return repositories.ForbiddenWordRepository.FindPageByParams(sqls.DB(), params)
}

func (s *forbiddenWordService) FindPageByCnd(cnd *sqls.Cnd) (list []models.ForbiddenWord, paging *sqls.Paging) {
	return repositories.ForbiddenWordRepository.FindPageByCnd(sqls.DB(), cnd)
}

func (s *forbiddenWordService) Count(cnd *sqls.Cnd) int64 {
	return repositories.ForbiddenWordRepository.Count(sqls.DB(), cnd)
}

func (s *forbiddenWordService) Create(t *models.ForbiddenWord) error {
	if err := repositories.ForbiddenWordRepository.Create(sqls.DB(), t); err != nil {
		return err
	}
	cache.ForbiddenWordCache.Invalidate()
	return nil
}

func (s *forbiddenWordService) Update(t *models.ForbiddenWord) error {
	if err := repositories.ForbiddenWordRepository.Update(sqls.DB(), t); err != nil {
		return err
	}
	cache.ForbiddenWordCache.Invalidate()
	return nil
}

func (s *forbiddenWordService) Updates(id int64, columns map[string]interface{}) error {
	if err := repositories.ForbiddenWordRepository.Updates(sqls.DB(), id, columns); err != nil {
		return err
	}
	cache.ForbiddenWordCache.Invalidate()
	return nil
}

func (s *forbiddenWordService) UpdateColumn(id int64, name string, value interface{}) error {
	if err := repositories.ForbiddenWordRepository.UpdateColumn(sqls.DB(), id, name, value); err != nil {
		return err
	}
	cache.ForbiddenWordCache.Invalidate()
	return nil
}

func (s *forbiddenWordService) Delete(id int64) {
	repositories.ForbiddenWordRepository.Delete(sqls.DB(), id)
	cache.ForbiddenWordCache.Invalidate()
}

func (s forbiddenWordService) Check(content string) (hitWords []string) {
	if strs.IsBlank(content) {
		return
	}
	words := cache.ForbiddenWordCache.Get()
	if len(words) == 0 {
		return
	}
	for _, word := range words {
		if word.Type == constants.ForbiddenWordTypeWord {
			if strings.Contains(content, word.Word) {
				hitWords = append(hitWords, word.Word)
				break
			}
		} else if word.Type == constants.ForbiddenWordTypeRegex {
			// if matched, _ := regexp.MatchString(word.Word, content); matched {
			// 	hitWords = append(hitWords, word.Word)
			// 	break
			// }
			r, _ := regexp.Compile(word.Word)
			if r != nil {
				hits := r.FindAllString(content, 3)
				if len(hits) > 0 {
					hitWords = append(hitWords, hits...)
					break
				}
			}
		}
	}
	return
}
