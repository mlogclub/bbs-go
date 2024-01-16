package render

import (
	"bbs-go/internal/models"
	"bbs-go/internal/models/constants"
	"bbs-go/internal/pkg/bbsurls"
	"bbs-go/internal/pkg/common"
	"bbs-go/internal/pkg/text"
	"bbs-go/internal/services"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func BuildFavorite(favorite *models.Favorite) *models.FavoriteResponse {
	rsp := &models.FavoriteResponse{}
	rsp.Id = favorite.Id
	rsp.EntityType = favorite.EntityType
	rsp.CreateTime = favorite.CreateTime

	if favorite.EntityType == constants.EntityArticle {
		article := services.ArticleService.Get(favorite.EntityId)
		if article == nil || article.Status != constants.StatusOk {
			rsp.Deleted = true
		} else {
			rsp.Url = bbsurls.ArticleUrl(article.Id)
			rsp.User = BuildUserInfoDefaultIfNull(article.UserId)
			rsp.Title = article.Title
			if article.ContentType == constants.ContentTypeMarkdown {
				rsp.Content = common.GetMarkdownSummary(article.Content)
			} else if article.ContentType == constants.ContentTypeHtml {
				doc, err := goquery.NewDocumentFromReader(strings.NewReader(article.Content))
				if err == nil {
					rsp.Content = text.GetSummary(doc.Text(), constants.SummaryLen)
				}
			}
		}
	} else {
		topic := services.TopicService.Get(favorite.EntityId)
		if topic == nil || topic.Status != constants.StatusOk {
			rsp.Deleted = true
		} else {
			rsp.Url = bbsurls.TopicUrl(topic.Id)
			rsp.User = BuildUserInfoDefaultIfNull(topic.UserId)
			rsp.Title = topic.Title
			rsp.Content = common.GetMarkdownSummary(topic.Content)
		}
	}
	return rsp
}

func BuildFavorites(favorites []models.Favorite) []models.FavoriteResponse {
	if len(favorites) == 0 {
		return nil
	}
	var responses []models.FavoriteResponse
	for _, favorite := range favorites {
		responses = append(responses, *BuildFavorite(&favorite))
	}
	return responses
}
