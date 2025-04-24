# BBS-GO Makefile

# 默认目标
.PHONY: all
all: build

# 构建所有组件
.PHONY: build
build: clean build-server build-site build-admin
	@mkdir -p dist
	@cp -r server/bbs-go dist/
	@cp -r server/migrations dist/migrations
	@cp -r site/dist dist/site
	@cp -r admin/dist dist/admin

# 构建所有平台的服务器
.PHONY: build-all-platforms
build-all-platforms: clean build-site build-admin
	@echo "build macos amd64..."
	@cd server && GOOS=darwin GOARCH=amd64 go build -v -o bbs-go-macos-amd64 main.go
	@mkdir -p dist/bbs-go-macos-amd64
	@cp -r server/bbs-go-macos-amd64 dist/bbs-go-macos-amd64/bbs-go
	@cp -r server/migrations dist/bbs-go-macos-amd64/migrations
	@cp -r site/dist dist/bbs-go-macos-amd64/site
	@cp -r admin/dist dist/bbs-go-macos-amd64/admin
	@zip -r dist/bbs-go-macos-amd64.zip dist/bbs-go-macos-amd64
	@rm -rf dist/bbs-go-macos-amd64
	@echo "build macos amd64 done"

	@echo "build macos arm64..."
	@cd server && GOOS=darwin GOARCH=arm64 go build -v -o bbs-go-macos-arm64 main.go
	@mkdir -p dist/bbs-go-macos-arm64
	@cp -r server/bbs-go-macos-arm64 dist/bbs-go-macos-arm64/bbs-go
	@cp -r server/migrations dist/bbs-go-macos-arm64/migrations
	@cp -r site/dist dist/bbs-go-macos-arm64/site
	@cp -r admin/dist dist/bbs-go-macos-arm64/admin
	@zip -r dist/bbs-go-macos-arm64.zip dist/bbs-go-macos-arm64
	@rm -rf dist/bbs-go-macos-arm64
	@echo "build macos arm64 done"

	@echo "build windows amd64..."
	@cd server && GOOS=windows GOARCH=amd64 go build -v -o bbs-go-windows-amd64.exe main.go
	@mkdir -p dist/bbs-go-windows-amd64
	@cp -r server/bbs-go-windows-amd64.exe dist/bbs-go-windows-amd64/bbs-go.exe
	@cp -r server/migrations dist/bbs-go-windows-amd64/migrations
	@cp -r site/dist dist/bbs-go-windows-amd64/site
	@cp -r admin/dist dist/bbs-go-windows-amd64/admin
	@zip -r dist/bbs-go-windows-amd64.zip dist/bbs-go-windows-amd64
	@rm -rf dist/bbs-go-windows-amd64
	@echo "build windows amd64 done"

	@echo "build windows 386..."
	@cd server && GOOS=windows GOARCH=386 go build -v -o bbs-go-windows-386.exe main.go
	@mkdir -p dist/bbs-go-windows-386
	@cp -r server/bbs-go-windows-386.exe dist/bbs-go-windows-386/bbs-go.exe
	@cp -r server/migrations dist/bbs-go-windows-386/migrations
	@cp -r site/dist dist/bbs-go-windows-386/site
	@cp -r admin/dist dist/bbs-go-windows-386/admin
	@zip -r dist/bbs-go-windows-386.zip dist/bbs-go-windows-386
	@rm -rf dist/bbs-go-windows-386
	@echo "build windows 386 done"

	@echo "build linux amd64..."
	@cd server && GOOS=linux GOARCH=amd64 go build -v -o bbs-go-linux-amd64 main.go
	@mkdir -p dist/bbs-go-linux-amd64
	@cp -r server/bbs-go-linux-amd64 dist/bbs-go-linux-amd64/bbs-go
	@cp -r server/migrations dist/bbs-go-linux-amd64/migrations
	@cp -r site/dist dist/bbs-go-linux-amd64/site
	@cp -r admin/dist dist/bbs-go-linux-amd64/admin
	@zip -r dist/bbs-go-linux-amd64.zip dist/bbs-go-linux-amd64
	@rm -rf dist/bbs-go-linux-amd64
	@echo "build linux amd64 done"

	@echo "build linux 386..."
	@cd server && GOOS=linux GOARCH=386 go build -v -o bbs-go-linux-386 main.go
	@mkdir -p dist/bbs-go-linux-386
	@cp -r server/bbs-go-linux-386 dist/bbs-go-linux-386/bbs-go
	@cp -r server/migrations dist/bbs-go-linux-386/migrations
	@cp -r site/dist dist/bbs-go-linux-386/site
	@cp -r admin/dist dist/bbs-go-linux-386/admin
	@zip -r dist/bbs-go-linux-386.zip dist/bbs-go-linux-386
	@rm -rf dist/bbs-go-linux-386
	@echo "build linux 386 done"

	@echo "all done"

# 构建服务器
.PHONY: build-server
build-server:
	@echo "build server..."
	@cd server && go build -v -o bbs-go main.go

# 构建前端站点
.PHONY: build-site
build-site:
	@echo "build site..."
	@cd site && pnpm install && pnpm generate

# 构建管理后台
.PHONY: build-admin
build-admin:
	@echo "build admin..."
	@cd admin && pnpm install && pnpm build

# 清理构建产物
.PHONY: clean
clean: 
	@echo "clean server..."
	@cd server && rm -f bbs-go bbs-go-* bbs-go-*.exe

	@echo "clean site..."
	@cd site && rm -rf .nuxt .output dist

	@echo "clean admin..."
	@cd admin && rm -rf dist

	@echo "clean dist..."
	@rm -rf dist


# 帮助信息
.PHONY: help
help:
	@echo "BBS-GO Makefile 帮助信息:"
	@echo "  make build               - 构建所有组件"
	@echo "  make build-all-platforms - 构建所有平台的服务器"
	@echo "  make build-server        - 构建服务器"
	@echo "  make build-site          - 构建前端站点"
	@echo "  make build-admin         - 构建管理后台"
	@echo "  make clean               - 清理所有构建产物"
	@echo "  make help                - 显示帮助信息"