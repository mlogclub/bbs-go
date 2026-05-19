package services

import (
	"encoding/xml"
	"errors"
	"strconv"
	"strings"
	"testing"
)

func TestBuildSitemapXML_EscapesURLs(t *testing.T) {
	xml := buildSitemapXML([]sitemapURLItem{{Loc: "https://example.com/topics/tag/1?name=a&b=c", LastMod: "2026-05-18"}})
	for _, expected := range []string{
		`<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">`,
		`<loc>https://example.com/topics/tag/1?name=a&amp;b=c</loc>`,
		`<lastmod>2026-05-18</lastmod>`,
	} {
		if !strings.Contains(xml, expected) {
			t.Fatalf("expected sitemap XML to contain %q, got:\n%s", expected, xml)
		}
	}
}

func TestBuildSitemapIndexXML_UsesUploadedChildURLs(t *testing.T) {
	xml := buildSitemapIndexXML([]sitemapIndexItem{{Loc: "https://cdn.example.com/seo/sitemap-static.xml", LastMod: "2026-05-18"}, {Loc: "https://cdn.example.com/seo/sitemap-topics-1.xml", LastMod: "2026-05-18"}})
	for _, expected := range []string{
		`<sitemapindex xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">`,
		`<loc>https://cdn.example.com/seo/sitemap-static.xml</loc>`,
		`<loc>https://cdn.example.com/seo/sitemap-topics-1.xml</loc>`,
	} {
		if !strings.Contains(xml, expected) {
			t.Fatalf("expected sitemap index XML to contain %q, got:\n%s", expected, xml)
		}
	}
}

type fakeSitemapUploader struct {
	failOnKey   string
	uploads     map[string]string
	contentType map[string]string
}

func (f *fakeSitemapUploader) PutObject(key string, body []byte, contentType string) (string, error) {
	if f.failOnKey == key {
		return "", errors.New("upload failed")
	}
	if f.uploads == nil {
		f.uploads = map[string]string{}
	}
	if f.contentType == nil {
		f.contentType = map[string]string{}
	}
	f.uploads[key] = string(body)
	f.contentType[key] = contentType
	return "https://cdn.example.com/" + key, nil
}

func TestUploadSitemapFiles_DoesNotReturnIndexWhenChildUploadFails(t *testing.T) {
	uploader := &fakeSitemapUploader{failOnKey: "seo/sitemap-static-1.xml"}
	_, err := uploadSitemapFiles(uploader, []generatedSitemapFile{{Key: "seo/sitemap-static-1.xml", XML: "<urlset></urlset>"}}, "2026-05-18")
	if err == nil {
		t.Fatal("expected upload error")
	}
	if _, ok := uploader.uploads["seo/sitemap-index.xml"]; ok {
		t.Fatal("sitemap index must not upload when a child sitemap upload fails")
	}
}

func TestUploadSitemapFiles_UsesXMLContentType(t *testing.T) {
	uploader := &fakeSitemapUploader{}
	_, err := uploadSitemapFiles(uploader, []generatedSitemapFile{{Key: "seo/sitemap-static-1.xml", XML: "<urlset></urlset>"}}, "2026-05-18")
	if err != nil {
		t.Fatalf("expected upload to succeed, got %v", err)
	}

	for _, key := range []string{"seo/sitemap-static-1.xml", "seo/sitemap-index.xml"} {
		if got := uploader.contentType[key]; got != "application/xml; charset=utf-8" {
			t.Fatalf("expected %s content type to be application/xml; charset=utf-8, got %q", key, got)
		}
	}
}

func TestUploadSitemapFiles_NilClientReturnsError(t *testing.T) {
	_, err := uploadSitemapFiles(nil, []generatedSitemapFile{{Key: "seo/sitemap-static-1.xml", XML: "<urlset></urlset>"}}, "2026-05-18")
	if err == nil {
		t.Fatal("expected nil client error")
	}
}

func TestUploadSitemapFiles_EmptyKeyReturnsError(t *testing.T) {
	uploader := &fakeSitemapUploader{}
	_, err := uploadSitemapFiles(uploader, []generatedSitemapFile{{Key: " \t\n", XML: "<urlset></urlset>"}}, "2026-05-18")
	if err == nil {
		t.Fatal("expected empty key error")
	}
	if len(uploader.uploads) != 0 {
		t.Fatalf("expected no uploads for empty key, got %d", len(uploader.uploads))
	}
}

func TestUploadSitemapFiles_UsesFixedKeysForChildrenAndIndex(t *testing.T) {
	uploader := &fakeSitemapUploader{}
	indexURL, err := uploadSitemapFiles(uploader, []generatedSitemapFile{
		{Key: "seo/sitemap-static-1.xml", XML: "<urlset></urlset>"},
		{Key: "seo/sitemap-topics-1.xml", XML: "<urlset></urlset>"},
	}, "2026-05-18")
	if err != nil {
		t.Fatalf("expected upload to succeed, got %v", err)
	}
	if indexURL != "https://cdn.example.com/seo/sitemap-index.xml" {
		t.Fatalf("indexURL=%q", indexURL)
	}
	for _, key := range []string{
		"seo/sitemap-static-1.xml",
		"seo/sitemap-topics-1.xml",
		"seo/sitemap-index.xml",
	} {
		if _, ok := uploader.uploads[key]; !ok {
			t.Fatalf("expected upload key %q", key)
		}
	}
}

func TestBuildContentSitemapFiles_SplitsByBatchSize(t *testing.T) {
	items := make([]sitemapURLItem, 0, 3)
	for i := 1; i <= 3; i++ {
		items = append(items, sitemapURLItem{Loc: "https://example.com/topic/" + strconv.Itoa(i), LastMod: "2026-05-18"})
	}
	files := buildContentSitemapFiles("seo/sitemap-topics", items, 2)
	if len(files) != 2 {
		t.Fatalf("len(files)=%d want 2", len(files))
	}
	if files[0].Key != "seo/sitemap-topics-1.xml" {
		t.Fatalf("first key=%q", files[0].Key)
	}
	if files[1].Key != "seo/sitemap-topics-2.xml" {
		t.Fatalf("second key=%q", files[1].Key)
	}
}

func TestBuildSitemapContentFiles_DoesNotExceedBatchSize(t *testing.T) {
	items := make([]sitemapURLItem, 0, 3)
	for i := 1; i <= 3; i++ {
		items = append(items, sitemapURLItem{Loc: "https://example.com/topic/" + strconv.Itoa(i)})
	}
	files := buildContentSitemapFiles("seo/sitemap-topics", items, 2)
	for _, file := range files {
		var urlset sitemapURLSet
		if err := xml.Unmarshal([]byte(file.XML), &urlset); err != nil {
			t.Fatalf("failed to unmarshal %s: %v", file.Key, err)
		}
		if len(urlset.URLs) > 2 {
			t.Fatalf("%s has %d urls, want at most 2", file.Key, len(urlset.URLs))
		}
	}
}

func TestHasAbsoluteBaseURL(t *testing.T) {
	tests := []struct {
		name    string
		baseURL string
		want    bool
	}{
		{name: "https host", baseURL: "https://example.com", want: true},
		{name: "http host with path", baseURL: "http://example.com/path", want: true},
		{name: "relative root", baseURL: "/", want: false},
		{name: "empty", baseURL: "", want: false},
		{name: "https missing host", baseURL: "https://", want: false},
		{name: "https empty host with slash path", baseURL: "https:///", want: false},
		{name: "unsupported scheme", baseURL: "ftp://example.com", want: false},
		{name: "query", baseURL: "https://example.com?x=1", want: false},
		{name: "fragment", baseURL: "https://example.com#frag", want: false},
		{name: "opaque", baseURL: "https:example.com", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := hasAbsoluteBaseURL(tt.baseURL); got != tt.want {
				t.Fatalf("hasAbsoluteBaseURL(%q)=%v want %v", tt.baseURL, got, tt.want)
			}
		})
	}
}
