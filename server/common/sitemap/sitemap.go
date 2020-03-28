package sitemap

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"

	"bbs-go/common/uploader"
)

type Generator struct {
	opts        *Options
	index       int
	sitemapFunc SitemapFunc
	URLs        []URL
	IndexURLs   []IndexURL
}

type SitemapFunc func(sm *Generator, sitemapLoc string)

func NewGenerator(sitemapHost, sitemapPath, sitemapName string, sitemapFunc SitemapFunc) *Generator {
	return &Generator{
		opts:        NewOptions(sitemapHost, sitemapPath, sitemapName),
		index:       1,
		sitemapFunc: sitemapFunc,
	}
}

func (sm *Generator) AddURL(url URL) {
	sm.URLs = append(sm.URLs, url)

	// logrus.Info("Sitemap add url ", url.Loc)

	if len(sm.URLs) >= MaxSitemapLinks-1 {
		sm.Finalize()
	}
}

func (sm *Generator) Finalize() {
	if len(sm.URLs) == 0 {
		return
	}
	defer func() {
		// Clean current sitemap urls
		sm.URLs = nil

		// Generator index + 1
		sm.index++
	}()

	sitemapPath := sm.opts.SitemapPathInPublic("-" + strconv.Itoa(sm.index) + SitemapXmlExt)
	sitemapLoc := write(sitemapPath, XmlContent(sm.URLs))
	if len(sitemapLoc) > 0 {
		// execute callback
		if sm.sitemapFunc != nil {
			sm.sitemapFunc(sm, sitemapLoc)
		}
		// ping search engine
		go func() {
			PingSearchEngines(sitemapLoc)
		}()
	}
}

func (sm *Generator) WriteIndex(sitemapLocs []IndexURL) string {
	sitemapPath := sm.opts.SitemapIndexPathInPublic(SitemapXmlExt)
	sitemapIndexLoc := write(sitemapPath, IndexXmlContent(sitemapLocs))
	PingSearchEngines(sitemapIndexLoc)
	return sitemapIndexLoc
}

// write sitemap and index to aliyun oss
func write(path, xml string) (sitemapUrl string) {
	if len(xml) > 0 {
		sitemapUrl, _ = uploader.PutObject(path, []byte(xml))
		logrus.Info("Upload sitemap success, " + sitemapUrl)
	}
	return
}

func XmlContent(urls []URL) string {
	if len(urls) == 0 {
		return ""
	}
	b := strings.Builder{}
	b.Write(XMLHeader)
	for _, url := range urls {
		b.WriteString(url.String())
	}
	b.Write(XMLFooter)
	return b.String()
}

func IndexXmlContent(sitemapLocs []IndexURL) string {
	if len(sitemapLocs) == 0 {
		return ""
	}
	b := strings.Builder{}
	b.Write(IndexXMLHeader)
	for _, loc := range sitemapLocs {
		b.WriteString(loc.String())
	}
	b.Write(IndexXMLFooter)
	return b.String()
}

func PingSearchEngines(sitemapLoc string) {
	engines := []string{
		"http://www.google.com/webmasters/tools/ping?sitemap=%s",
		"http://www.bing.com/webmaster/ping.aspx?siteMap=%s",
	}

	bufs := len(engines)
	does := make(chan string, bufs)
	client := http.Client{Timeout: 5 * time.Second}

	for _, engine := range engines {
		go func(baseurl string) {
			url := fmt.Sprintf(baseurl, sitemapLoc)
			println("Ping now:", url)

			resp, err := client.Get(url)
			if err != nil {
				does <- fmt.Sprintf("[E] Ping failed: %s (URL:%s)",
					err, url)
				return
			}
			defer resp.Body.Close()

			does <- fmt.Sprintf("Successful ping of `%s`", url)
		}(engine)
	}

	for i := 0; i < bufs; i++ {
		println(<-does)
	}
}
