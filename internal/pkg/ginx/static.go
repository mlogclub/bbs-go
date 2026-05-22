package ginx

import (
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/gin-gonic/gin"
)

type noDirFS struct {
	fs http.FileSystem
}

type DirOptions struct {
	ShowList  bool
	SPA       bool
	IndexName string
}

type SPAOptions struct {
	Root             string
	EmbeddedFS       fs.FS
	EmbeddedRoot     string
	DirOptions       DirOptions
	NotFoundPrefixes []string
	NotFoundHandler  gin.HandlerFunc
}

func StaticFiles(root string) http.FileSystem {
	return noDirFS{fs: http.Dir(root)}
}

func (n noDirFS) Open(name string) (http.File, error) {
	file, err := n.fs.Open(name)
	if err != nil {
		return nil, err
	}
	info, err := file.Stat()
	if err != nil {
		_ = file.Close()
		return nil, err
	}
	if info.IsDir() {
		_ = file.Close()
		return nil, os.ErrNotExist
	}
	return file, nil
}

func HandleSPA(engine *gin.Engine, options SPAOptions) gin.HandlerFunc {
	handler := NewSPAHandler(options.Root, options.EmbeddedFS, options.EmbeddedRoot, options.DirOptions)
	engine.GET("/", handler)
	engine.HEAD("/", handler)
	engine.NoRoute(func(ctx *gin.Context) {
		for _, prefix := range options.NotFoundPrefixes {
			if strings.HasPrefix(ctx.Request.URL.Path, prefix) {
				if options.NotFoundHandler != nil {
					options.NotFoundHandler(ctx)
					return
				}
				ctx.AbortWithStatus(http.StatusNotFound)
				return
			}
		}
		handler(ctx)
	})
	return handler
}

func NewSPAHandler(root string, embeddedFS fs.FS, embeddedRoot string, options DirOptions) gin.HandlerFunc {
	if _, err := os.Stat(path.Join(root, options.IndexName)); err == nil {
		return DirHandler(http.Dir(root), options)
	}

	spaFS, err := fs.Sub(embeddedFS, embeddedRoot)
	if err != nil {
		slog.Error("failed to load embedded SPA files", slog.Any("err", err))
		return func(ctx *gin.Context) {
			ctx.AbortWithStatus(http.StatusInternalServerError)
		}
	}
	return DirHandler(http.FS(spaFS), options)
}

func DirHandler(fileSystem http.FileSystem, options DirOptions) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		name := path.Clean(ctx.Request.URL.Path)
		if name == "." || name == "/" {
			name = "/" + options.IndexName
		}

		file, err := fileSystem.Open(name)
		if err == nil {
			defer file.Close()
			info, err := file.Stat()
			if err != nil {
				ctx.AbortWithStatus(http.StatusInternalServerError)
				return
			}
			if info.IsDir() {
				if options.IndexName != "" {
					indexName := path.Join(name, options.IndexName)
					indexFile, err := fileSystem.Open(indexName)
					if err == nil {
						defer indexFile.Close()
						indexInfo, statErr := indexFile.Stat()
						if statErr != nil {
							ctx.AbortWithStatus(http.StatusInternalServerError)
							return
						}
						if !indexInfo.IsDir() {
							serveFile(ctx, indexName, indexFile, indexInfo)
							return
						}
					}
					if err != nil && !os.IsNotExist(err) {
						ctx.AbortWithStatus(http.StatusInternalServerError)
						return
					}
				}
				if !options.ShowList {
					ctx.AbortWithStatus(http.StatusNotFound)
					return
				}
			} else {
				serveFile(ctx, name, file, info)
				return
			}
		} else if !os.IsNotExist(err) {
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		if options.SPA && options.IndexName != "" {
			indexName := "/" + options.IndexName
			indexFile, err := fileSystem.Open(indexName)
			if err != nil {
				if os.IsNotExist(err) {
					ctx.AbortWithStatus(http.StatusNotFound)
					return
				}
				ctx.AbortWithStatus(http.StatusInternalServerError)
				return
			}
			defer indexFile.Close()
			indexInfo, err := indexFile.Stat()
			if err != nil {
				ctx.AbortWithStatus(http.StatusInternalServerError)
				return
			}
			serveFile(ctx, indexName, indexFile, indexInfo)
			return
		}

		ctx.AbortWithStatus(http.StatusNotFound)
	}
}

func serveFile(ctx *gin.Context, name string, file http.File, info os.FileInfo) {
	http.ServeContent(ctx.Writer, ctx.Request, name, info.ModTime(), file)
}
