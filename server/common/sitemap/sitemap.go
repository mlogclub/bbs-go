package sitemap

import (
	"fmt"
	"strings"
	"time"

	"bbs-go/common/oss"
)

const (
	MaxSitemapLinks  = 50000
	SitemapIndexName = "sitemap.xml"
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
	if len(sm.URLs) > 0 {
		// Add current sitemap to sitemap index
		sitemapLoc := sm.opts.SitemapLoc(".xml")
		sm.IndexURLs = append(sm.IndexURLs, SitemapIndexURL{
			Loc:     sitemapLoc,
			Lastmod: time.Now(),
		})

		// Clean current sitemap urls
		sm.URLs = nil
	}

	sm.WriteToOSS()
}

// WriteToOSS write sitemap and index to aliyun oss
func (sm *Sitemap) WriteToOSS() {
	// Upload sitemap
	sitemapXml := sm.SitemapXml()
	sitemapUrl, _ := oss.PutObject(sm.opts.SitemapPathInPublic(".xml"), []byte(sitemapXml))
	fmt.Println(sitemapUrl)

	// Upload sitemap index
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
