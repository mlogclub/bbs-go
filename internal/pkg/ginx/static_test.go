package ginx

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestDirHandlerWithSPA(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "index.html"), []byte("<html>spa</html>"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "manifest"), []byte("manifest content"), 0o644); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		path       string
		accept     string
		wantStatus int
		wantBody   string
	}{
		{path: "/assets/not-exists.css", accept: "text/css,*/*;q=0.1", wantStatus: http.StatusNotFound},
		{path: "/images/not-exists.png", accept: "image/avif,image/webp,image/png,image/*,*/*;q=0.8", wantStatus: http.StatusNotFound},
		{path: "/fonts/not-exists", accept: "*/*", wantStatus: http.StatusNotFound},
		{path: "/manifest", accept: "*/*", wantStatus: http.StatusOK, wantBody: "manifest content"},
		{path: "/topic/123", accept: "text/html,application/xhtml+xml", wantStatus: http.StatusOK, wantBody: "<html>spa</html>"},
	}
	handler := DirHandler(http.Dir(root), DirOptions{
		ShowList:  false,
		SPA:       true,
		IndexName: "index.html",
	})

	for _, tt := range tests {
		rec := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(rec)
		ctx.Request = httptest.NewRequest(http.MethodGet, tt.path, nil)
		ctx.Request.Header.Set("Accept", tt.accept)

		handler(ctx)

		if rec.Code != tt.wantStatus {
			t.Fatalf("%s status=%d want %d; body=%q", tt.path, rec.Code, tt.wantStatus, rec.Body.String())
		}
		if tt.wantBody != "" && strings.TrimSpace(rec.Body.String()) != tt.wantBody {
			t.Fatalf("%s body=%q want %q", tt.path, strings.TrimSpace(rec.Body.String()), tt.wantBody)
		}
	}
}

func TestStaticFilesDoesNotOpenDirectories(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "file.txt"), []byte("content"), 0o644); err != nil {
		t.Fatal(err)
	}

	fileSystem := StaticFiles(root)
	if file, err := fileSystem.Open("/file.txt"); err != nil {
		t.Fatalf("Open(file.txt) err=%v", err)
	} else {
		_ = file.Close()
	}

	if file, err := fileSystem.Open("/"); err == nil {
		_ = file.Close()
		t.Fatal("Open(/) succeeded, want error")
	} else if !os.IsNotExist(err) {
		t.Fatalf("Open(/) err=%v, want not exist", err)
	}
}

func TestHandleSPABuildsHandlerOnceAndKeepsNotFoundPrefixes(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "index.html"), []byte("<html>spa</html>"), 0o644); err != nil {
		t.Fatal(err)
	}

	engine := gin.New()
	HandleSPA(engine, SPAOptions{
		Root: root,
		DirOptions: DirOptions{
			ShowList:  false,
			SPA:       true,
			IndexName: "index.html",
		},
		NotFoundPrefixes: []string{"/api/"},
	})

	tests := []struct {
		path       string
		wantStatus int
		wantBody   string
	}{
		{path: "/", wantStatus: http.StatusOK, wantBody: "<html>spa</html>"},
		{path: "/topic/123", wantStatus: http.StatusOK, wantBody: "<html>spa</html>"},
		{path: "/api/not-exists", wantStatus: http.StatusNotFound},
	}

	for _, tt := range tests {
		rec := httptest.NewRecorder()
		engine.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, tt.path, nil))

		if rec.Code != tt.wantStatus {
			t.Fatalf("%s status=%d want %d", tt.path, rec.Code, tt.wantStatus)
		}
		if tt.wantBody != "" && strings.TrimSpace(rec.Body.String()) != tt.wantBody {
			t.Fatalf("%s body=%q want %q", tt.path, strings.TrimSpace(rec.Body.String()), tt.wantBody)
		}
	}
}

func TestSPARouteManifestStatus(t *testing.T) {
	root := t.TempDir()
	for name, body := range map[string]string{"index.html": "<html>spa</html>", ".route-patterns.json": `["^/$", "^/topic/[^/]+/?$", "^/user/signin/callback/.*$"]`} {
		if err := os.WriteFile(filepath.Join(root, name), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	app := gin.New()
	HandleSPA(app, SPAOptions{Root: root, DirOptions: DirOptions{SPA: true, IndexName: "index.html"}})
	for _, tc := range []struct {
		path   string
		status int
	}{
		{"/", 200}, {"/topic/123", 200}, {"/missing", 404}, {"/assets/missing.js", 404},
		{"/user/signin/callback/token.with.dots", 200},
	} {
		for _, method := range []string{http.MethodGet, http.MethodHead} {
			rec := httptest.NewRecorder()
			app.ServeHTTP(rec, httptest.NewRequest(method, tc.path, nil))
			if rec.Code != tc.status {
				t.Fatalf("%s %s: got %d want %d", method, tc.path, rec.Code, tc.status)
			}
			if method == http.MethodHead && rec.Body.Len() != 0 {
				t.Fatal("HEAD must not have a response body")
			}
			if method == http.MethodGet && tc.path == "/missing" && rec.Body.String() != "<html>spa</html>" {
				t.Fatal("missing route must retain the SPA error page shell")
			}
		}
	}
}
