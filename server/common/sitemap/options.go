package sitemap

import (
	"net/url"
	"path/filepath"
)

type Options struct {
	SitemapHost string
	SitemapPath string
	SitemapName string
}

func NewOptions(sitemapHost, sitemapPath, sitemapName string) *Options {
	return &Options{
		SitemapHost: sitemapHost,
		SitemapPath: sitemapPath,
		SitemapName: sitemapName,
	}
}

// SitemapLoc sitemap loc
func (opts *Options) SitemapLoc(ext string) string {
	base, _ := url.Parse(opts.SitemapHost)
	path := opts.SitemapPathInPublic(ext)

	for _, ref := range []string{
		path,
	} {
		base, _ = base.Parse(ref)
	}

	return base.String()
}

// SitemapPathInPublic returns path to combine sitemapsPath and Filename on website.
// It also indicates where url file path is
func (opts *Options) SitemapPathInPublic(ext string) string {
	return filepath.Join(opts.SitemapPath, opts.SitemapName) + ext
}
