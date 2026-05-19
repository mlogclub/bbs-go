package services

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"log/slog"
	"net/url"
	"strconv"
	"strings"
	"time"

	"bbs-go/internal/models"
	"bbs-go/internal/models/constants"
	"bbs-go/internal/pkg/bbsurls"
	"bbs-go/internal/pkg/uploader"

	"github.com/mlogclub/simple/sqls"
)

const (
	sitemapBatchSize   = 10000
	sitemapContentType = "application/xml; charset=utf-8"
	sitemapIndexKey    = "seo/sitemap-index.xml"
)

type seoSitemapService struct{}

var SeoSitemapService = &seoSitemapService{}

func (s *seoSitemapService) RedirectURL() string {
	return strings.TrimSpace(UploadService.ObjectURL(sitemapIndexKey))
}

type sitemapURLItem struct {
	Loc     string `xml:"loc"`
	LastMod string `xml:"lastmod,omitempty"`
}

type sitemapIndexItem struct {
	Loc     string `xml:"loc"`
	LastMod string `xml:"lastmod,omitempty"`
}

type generatedSitemapFile struct {
	Key string
	XML string
}

type sitemapUploadClient interface {
	PutObject(key string, body []byte, contentType string) (string, error)
}

type uploadServiceSitemapClient struct{}

func (c uploadServiceSitemapClient) PutObject(key string, body []byte, contentType string) (string, error) {
	return UploadService.PutObject(key, bytes.NewReader(body), &uploader.PutOptions{
		ContentType:   contentType,
		ContentLength: int64(len(body)),
	})
}

func (s *seoSitemapService) GenerateAndUpload() error {
	return s.generateAndUpload(uploadServiceSitemapClient{})
}

func (s *seoSitemapService) generateAndUpload(client sitemapUploadClient) error {
	baseURL := SysConfigService.GetBaseURL()
	if !hasAbsoluteBaseURL(baseURL) {
		err := fmt.Errorf("baseURL must start with http:// or https://: %s", baseURL)
		slog.Warn("skip sitemap generation: invalid baseURL", slog.String("baseURL", baseURL))
		return err
	}

	files, err := s.buildStaticSitemapFiles()
	if err != nil {
		return err
	}
	topicFiles, err := s.buildTopicSitemapFiles()
	if err != nil {
		return err
	}
	files = append(files, topicFiles...)
	articleFiles, err := s.buildArticleSitemapFiles()
	if err != nil {
		return err
	}
	files = append(files, articleFiles...)

	_, err = uploadSitemapFiles(client, files, todaySitemapLastMod())
	return err
}

type sitemapURLSet struct {
	XMLName xml.Name         `xml:"urlset"`
	Xmlns   string           `xml:"xmlns,attr"`
	URLs    []sitemapURLItem `xml:"url"`
}

type sitemapIndex struct {
	XMLName  xml.Name           `xml:"sitemapindex"`
	Xmlns    string             `xml:"xmlns,attr"`
	Sitemaps []sitemapIndexItem `xml:"sitemap"`
}

func buildSitemapXML(items []sitemapURLItem) string {
	urls := make([]sitemapURLItem, 0, len(items))
	for _, item := range items {
		if strings.TrimSpace(item.Loc) == "" {
			continue
		}
		urls = append(urls, item)
	}
	return marshalXML(sitemapURLSet{
		Xmlns: "http://www.sitemaps.org/schemas/sitemap/0.9",
		URLs:  urls,
	})
}

func buildContentSitemapFiles(keyPrefix string, items []sitemapURLItem, batchSize int) []generatedSitemapFile {
	if batchSize <= 0 {
		batchSize = sitemapBatchSize
	}

	files := make([]generatedSitemapFile, 0, (len(items)+batchSize-1)/batchSize)
	for start := 0; start < len(items); start += batchSize {
		end := start + batchSize
		if end > len(items) {
			end = len(items)
		}
		files = append(files, generatedSitemapFile{
			Key: keyPrefix + "-" + strconv.Itoa(len(files)+1) + ".xml",
			XML: buildSitemapXML(items[start:end]),
		})
	}
	return files
}

func (s *seoSitemapService) buildStaticSitemapFiles() ([]generatedSitemapFile, error) {
	items := []sitemapURLItem{
		{Loc: bbsurls.AbsUrl("/")},
		{Loc: bbsurls.AbsUrl("/topics")},
		{Loc: bbsurls.AbsUrl("/articles")},
		{Loc: bbsurls.AbsUrl("/about")},
		{Loc: bbsurls.AbsUrl("/links")},
	}

	var nodes []models.TopicNode
	if err := sqls.NewCnd().
		Eq("status", constants.StatusOk).
		Asc("id").
		Build(sqls.DB()).
		Find(&nodes).Error; err != nil {
		return nil, err
	}
	for _, node := range nodes {
		items = append(items, sitemapURLItem{
			Loc:     bbsurls.AbsUrl("/topics/node/" + strconv.FormatInt(node.Id, 10)),
			LastMod: seoSitemapLastMod(node.CreateTime),
		})
	}

	var tags []models.Tag
	if err := sqls.NewCnd().
		Eq("status", constants.StatusOk).
		Asc("id").
		Build(sqls.DB()).
		Find(&tags).Error; err != nil {
		return nil, err
	}
	for _, tag := range tags {
		tagID := strconv.FormatInt(tag.Id, 10)
		lastMod := seoSitemapLastMod(tag.UpdateTime)
		items = append(items,
			sitemapURLItem{Loc: bbsurls.AbsUrl("/topics/tag/" + tagID), LastMod: lastMod},
			sitemapURLItem{Loc: bbsurls.AbsUrl("/articles/tag/" + tagID), LastMod: lastMod},
		)
	}

	return buildContentSitemapFiles("seo/sitemap-static", items, sitemapBatchSize), nil
}

func (s *seoSitemapService) buildTopicSitemapFiles() ([]generatedSitemapFile, error) {
	var files []generatedSitemapFile
	var cursor int64

	for {
		var topics []models.Topic
		if err := sqls.NewCnd().
			Eq("status", constants.StatusOk).
			Gt("id", cursor).
			Asc("id").
			Limit(sitemapBatchSize).
			Build(sqls.DB()).
			Find(&topics).Error; err != nil {
			return nil, err
		}
		if len(topics) == 0 {
			break
		}

		items := make([]sitemapURLItem, 0, len(topics))
		for _, topic := range topics {
			items = append(items, sitemapURLItem{
				Loc:     bbsurls.TopicUrl(topic.Id),
				LastMod: seoSitemapLastMod(seoSitemapTopicLastModTime(topic)),
			})
		}
		files = append(files, generatedSitemapFile{
			Key: "seo/sitemap-topics-" + strconv.Itoa(len(files)+1) + ".xml",
			XML: buildSitemapXML(items),
		})
		cursor = topics[len(topics)-1].Id
	}

	return files, nil
}

func (s *seoSitemapService) buildArticleSitemapFiles() ([]generatedSitemapFile, error) {
	var files []generatedSitemapFile
	var cursor int64

	for {
		var articles []models.Article
		if err := sqls.NewCnd().
			Eq("status", constants.StatusOk).
			Gt("id", cursor).
			Asc("id").
			Limit(sitemapBatchSize).
			Build(sqls.DB()).
			Find(&articles).Error; err != nil {
			return nil, err
		}
		if len(articles) == 0 {
			break
		}

		items := make([]sitemapURLItem, 0, len(articles))
		for _, article := range articles {
			items = append(items, sitemapURLItem{
				Loc:     bbsurls.ArticleUrl(article.Id),
				LastMod: seoSitemapLastMod(seoSitemapArticleLastModTime(article)),
			})
		}
		files = append(files, generatedSitemapFile{
			Key: "seo/sitemap-articles-" + strconv.Itoa(len(files)+1) + ".xml",
			XML: buildSitemapXML(items),
		})
		cursor = articles[len(articles)-1].Id
	}

	return files, nil
}

func buildSitemapIndexXML(items []sitemapIndexItem) string {
	sitemaps := make([]sitemapIndexItem, 0, len(items))
	for _, item := range items {
		if strings.TrimSpace(item.Loc) == "" {
			continue
		}
		sitemaps = append(sitemaps, item)
	}
	return marshalXML(sitemapIndex{
		Xmlns:    "http://www.sitemaps.org/schemas/sitemap/0.9",
		Sitemaps: sitemaps,
	})
}

func marshalXML(v any) string {
	data, err := xml.MarshalIndent(v, "", "  ")
	if err != nil {
		return ""
	}

	var buf bytes.Buffer
	buf.WriteString(xml.Header)
	buf.Write(data)
	buf.WriteByte('\n')
	return buf.String()
}

func hasAbsoluteBaseURL(baseURL string) bool {
	parsed, err := url.Parse(strings.TrimSpace(baseURL))
	if err != nil {
		return false
	}
	if parsed.Opaque != "" || parsed.RawQuery != "" || parsed.Fragment != "" {
		return false
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return false
	}
	return strings.TrimSpace(parsed.Host) != ""
}

func seoSitemapTopicLastModTime(topic models.Topic) int64 {
	if topic.LastCommentTime > topic.CreateTime {
		return topic.LastCommentTime
	}
	return topic.CreateTime
}

func seoSitemapArticleLastModTime(article models.Article) int64 {
	if article.UpdateTime > article.CreateTime {
		return article.UpdateTime
	}
	return article.CreateTime
}

func seoSitemapLastMod(timestamp int64) string {
	if timestamp <= 0 {
		return ""
	}
	return time.UnixMilli(timestamp).UTC().Format("2006-01-02")
}

func todaySitemapLastMod() string {
	return time.Now().UTC().Format("2006-01-02")
}

func uploadSitemapFiles(client sitemapUploadClient, files []generatedSitemapFile, lastMod string) (string, error) {
	if client == nil {
		return "", errors.New("sitemap upload client is nil")
	}

	indexItems := make([]sitemapIndexItem, 0, len(files))
	for _, file := range files {
		key := strings.TrimSpace(file.Key)
		if key == "" {
			return "", fmt.Errorf("sitemap file key is empty")
		}
		key = strings.TrimLeft(key, "/")

		url, err := client.PutObject(key, []byte(file.XML), sitemapContentType)
		if err != nil {
			return "", err
		}
		indexItems = append(indexItems, sitemapIndexItem{
			Loc:     url,
			LastMod: lastMod,
		})
	}

	indexXML := buildSitemapIndexXML(indexItems)
	return client.PutObject(sitemapIndexKey, []byte(indexXML), sitemapContentType)
}
