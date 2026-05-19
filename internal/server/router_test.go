package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kataras/iris/v12"
)

func TestSitemapRedirectRouteReturnsNotFoundWhenEmpty(t *testing.T) {
	app := newSitemapRedirectTestApp("")

	req := httptest.NewRequest(http.MethodGet, "/sitemap.xml", nil)
	rec := httptest.NewRecorder()
	app.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status=%d want %d", rec.Code, http.StatusNotFound)
	}
}

func TestSitemapRedirectRouteReturnsFoundWhenConfigured(t *testing.T) {
	app := newSitemapRedirectTestApp(" https://cdn.example.com/seo/sitemap-index.xml ")

	req := httptest.NewRequest(http.MethodGet, "/sitemap.xml", nil)
	rec := httptest.NewRecorder()
	app.ServeHTTP(rec, req)

	if rec.Code != http.StatusFound {
		t.Fatalf("status=%d want %d", rec.Code, http.StatusFound)
	}
	if got := rec.Header().Get("Location"); got != "https://cdn.example.com/seo/sitemap-index.xml" {
		t.Fatalf("Location=%q", got)
	}
}

func newSitemapRedirectTestApp(redirectURL string) *iris.Application {
	app := iris.New()
	app.Get("/sitemap.xml", func(ctx iris.Context) {
		writeSitemapRedirect(ctx, redirectURL)
	})
	if err := app.Build(); err != nil {
		panic(err)
	}
	return app
}
