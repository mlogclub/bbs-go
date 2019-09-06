package cache

import (
	"github.com/goburrow/cache"
	"github.com/mlogclub/bbs-go/model"
	"github.com/mlogclub/bbs-go/repositories"
	"github.com/mlogclub/simple"
	"github.com/sirupsen/logrus"
	"time"
)

type tagCache struct {
	cache           cache.LoadingCache
	activeTagsCache cache.LoadingCache // 热门标签
	allTagsCache    cache.LoadingCache // 所有标签
}

var TagCache = newTagCache()

func newTagCache() *tagCache {
	return &tagCache{
		cache: cache.NewLoadingCache(
			func(key cache.Key) (value cache.Value, e error) {
				value = repositories.TagRepository.Get(simple.GetDB(), Key2Int64(key))
				return
			},
			cache.WithMaximumSize(1000),
			cache.WithExpireAfterAccess(30*time.Minute),
		),
		activeTagsCache: cache.NewLoadingCache(
			func(key cache.Key) (value cache.Value, e error) {
				dateFrom := simple.Timestamp(simple.WithTimeAsStartOfDay(time.Now()))
				rows, e := simple.GetDB().Raw("select tag_id, count(*) c from t_article_tag where create_time > ?"+
					" group by tag_id order by c desc limit 50", dateFrom).Rows()

				// rows, e := simple.GetDB().Raw("select tag_id, count(*) c from t_article_tag" +
				// 	" group by tag_id order by c desc limit 20").Rows()

				if e != nil {
					return
				}
				var tagIds []int64
				for rows.Next() {
					var tagId int64
					var c int
					err := rows.Scan(&tagId, &c)
					if err != nil {
						continue
					}
					tagIds = append(tagIds, tagId)
				}
				value = tagIds
				return
			},
			cache.WithMaximumSize(1),
			cache.WithRefreshAfterWrite(30*time.Minute),
		),
		allTagsCache: cache.NewLoadingCache(
			func(key cache.Key) (value cache.Value, e error) {
				tags, e := repositories.TagRepository.QueryCnd(simple.GetDB(), simple.NewQueryCnd("status = ?", model.TagStatusOk))
				if e != nil {
					return
				}
				value = tags
				return
			},
			cache.WithMaximumSize(1),
			cache.WithRefreshAfterWrite(1*time.Hour),
		),
	}
}

func (this *tagCache) Get(tagId int64) *model.Tag {
	val, err := this.cache.Get(tagId)
	if err != nil {
		logrus.Error(err)
		return nil
	}
	if val != nil {
		return val.(*model.Tag)
	}
	return nil
}

func (this *tagCache) GetList(tagIds []int64) (tags []model.Tag) {
	if len(tagIds) == 0 {
		return nil
	}
	for _, tagId := range tagIds {
		tag := this.Get(tagId)
		if tag != nil {
			tags = append(tags, *tag)
		}
	}
	return
}

func (this *tagCache) Invalidate(tagId int64) {
	this.cache.Invalidate(tagId)
}

func (this *tagCache) GetActiveTags() []model.Tag {
	val, err := this.activeTagsCache.Get("data")
	if err != nil {
		return nil
	}
	tagIds := val.([]int64)
	if len(tagIds) == 0 {
		return nil
	}
	var tags []model.Tag
	for _, tagId := range tagIds {
		tag := this.Get(tagId)
		if tag != nil {
			tags = append(tags, *tag)
		}
	}
	return tags
}

func (this *tagCache) GetAllTags() []model.Tag {
	val, err := this.allTagsCache.Get("data")
	if err != nil {
		return nil
	}
	return val.([]model.Tag)
}
