package sitemap

import (
	"fmt"
	"strings"
	"time"

	"bbs-go/common/oss"
)

type Sitemap struct {
	opts      *Options
	URLs      []URL
	IndexURLs []SitemapIndexURL
}

func NewSitemap(sitemapHost, sitemapPath, sitemapName string) *Sitemap {
	return &Sitemap{
		opts: NewOptions(sitemapHost, sitemapPath, sitemapName),
	}
}

func (sm *Sitemap) Add(url URL) {
	sm.URLs = append(sm.URLs, url)
}

func (sm *Sitemap) Write() {
	// Clean current sitemap urls
	defer func() {
		sm.URLs = nil
	}()

	// Add current sitemap to sitemap index
	if len(sm.URLs) > 0 {
		sitemapLoc := sm.opts.SitemapLoc(".xml")
		sm.IndexURLs = append(sm.IndexURLs, SitemapIndexURL{
			Loc:     sitemapLoc,
			Lastmod: time.Now(),
		})
	}

	sm.WriteToOSS()
}

// WriteToOSS write sitemap and index to aliyun oss
func (sm *Sitemap) WriteToOSS() {
	// Upload sitemap
	sitemapXml := sm.SitemapXml()
	sitemapUrl, _ := oss.PutObject(sm.opts.SitemapPathInPublic(SitemapXmlExt), []byte(sitemapXml))
	fmt.Println(sitemapUrl)

	// Upload sitemap index
	sitemapIndexXml := sm.SitemapIndexXml()
	sitemapIndexUrl, _ := oss.PutObject(sm.opts.SitemapIndexPathInPublic(SitemapXmlExt), []byte(sitemapIndexXml))
	fmt.Println(sitemapIndexUrl)
}

func (sm *Sitemap) SitemapXml() string {
	if len(sm.URLs) == 0 {
		return ""
	}
	b := strings.Builder{}
	b.Write(XMLHeader)
	for _, url := range sm.URLs {
		b.WriteString(url.String())
	}
	b.Write(XMLFooter)
	return b.String()
}

func (sm *Sitemap) SitemapIndexXml() string {
	if len(sm.IndexURLs) == 0 {
		return ""
	}
	b := strings.Builder{}
	b.Write(IndexXMLHeader)
	for _, url := range sm.IndexURLs {
		b.WriteString(url.String())
	}
	b.Write(IndexXMLFooter)
	return b.String()
}
