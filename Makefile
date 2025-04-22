# BBS-GO Makefile

# 默认目标
.PHONY: all
all: build

# 构建所有组件
.PHONY: build
build: clean build-server build-site build-admin
	@mkdir -p dist
	@cp -r server/bbs-go dist/
	@cp -r site/dist dist/site
	@cp -r admin/dist dist/admin

# 构建服务器
.PHONY: build-server
build-server:
	@echo "构建服务器..."
	@cd server && go build -v -o bbs-go main.go

# 构建Linux版本的服务器
.PHONY: build-server-linux
build-server-linux:
	@echo "构建Linux版本的服务器..."
	@cd server && GOOS=linux GOARCH=amd64 go build -v -o bbs-go main.go

# 构建前端站点
.PHONY: build-site
build-site:
	@echo "构建前端站点..."
	@cd site && pnpm install && pnpm generate

# 构建管理后台
.PHONY: build-admin
build-admin:
	@echo "构建管理后台..."
	@cd admin && pnpm install && pnpm build

# 清理构建产物
.PHONY: clean
clean: 
	@echo "清理服务器构建产物..."
	@cd server && rm -f bbs-go

	@echo "清理前端站点构建产物..."
	@cd site && rm -rf .nuxt .output dist

	@echo "清理管理后台构建产物..."
	@cd admin && rm -rf dist

	@echo "清理dist目录..."
	@rm -rf dist


# 帮助信息
.PHONY: help
help:
	@echo "BBS-GO Makefile 帮助信息:"
	@echo "  make build            - 构建所有组件"
	@echo "  make build-server     - 构建服务器"
	@echo "  make build-server-linux - 构建Linux版本的服务器"
	@echo "  make build-site       - 构建前端站点"
	@echo "  make build-admin      - 构建管理后台"
	@echo "  make clean            - 清理所有构建产物"
	@echo "  make help             - 显示帮助信息"