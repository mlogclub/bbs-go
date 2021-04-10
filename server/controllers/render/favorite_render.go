package render

import (
	"bbs-go/model"
	"bbs-go/model/constants"
	"bbs-go/package/common"
	"bbs-go/package/urls"
	"bbs-go/services"
	"github.com/PuerkitoBio/goquery"
	"github.com/mlogclub/simple"
	"strings"
)

func BuildFavorite(favorite *model.Favorite) *model.FavoriteResponse {
	rsp := &model.FavoriteResponse{}
	rsp.FavoriteId = favorite.Id
	rsp.EntityType = favorite.EntityType
	rsp.CreateTime = favorite.CreateTime

	if favorite.EntityType == constants.EntityArticle {
		article := services.ArticleService.Get(favorite.EntityId)
		if article == nil || article.Status != constants.StatusOk {
			rsp.Deleted = true
		} else {
			rsp.Url = urls.ArticleUrl(article.Id)
			rsp.User = BuildUserById(article.UserId)
			rsp.Title = article.Title
			if article.ContentType == constants.ContentTypeMarkdown {
				rsp.Content = common.GetMarkdownSummary(article.Content)
			} else if article.ContentType == constants.ContentTypeHtml {
				doc, err := goquery.NewDocumentFromReader(strings.NewReader(article.Content))
				if err == nil {
					text := doc.Text()
					rsp.Content = simple.GetSummary(text, constants.SummaryLen)
				}
			}
		}
	} else {
		topic := services.TopicService.Get(favorite.EntityId)
		if topic == nil || topic.Status != constants.StatusOk {
			rsp.Deleted = true
		} else {
			rsp.Url = urls.TopicUrl(topic.Id)
			rsp.User = BuildUserById(topic.UserId)
			rsp.Title = topic.Title
			rsp.Content = common.GetMarkdownSummary(topic.Content)
		}
	}
	return rsp
}

func BuildFavorites(favorites []model.Favorite) []model.FavoriteResponse {
	if favorites == nil || len(favorites) == 0 {
		return nil
	}
	var responses []model.FavoriteResponse
	for _, favorite := range favorites {
		responses = append(responses, *BuildFavorite(&favorite))
	}
	return responses
}
