.PHONY: help install build clean deploy deploy-admin deploy-user deploy-go all

# 默认目标
help:
	@echo "任务管理系统 - 打包部署脚本"
	@echo ""
	@echo "使用方法:"
	@echo "  make install        - 安装所有项目依赖"
	@echo "  make build          - 打包所有项目（admin + user + go-backend）"
	@echo "  make clean          - 清理所有构建产物"
	@echo "  make deploy         - 部署所有项目到服务器"
	@echo "  make deploy-admin   - 仅部署 admin 前端"
	@echo "  make deploy-user    - 仅部署 user 前端"
	@echo "  make deploy-go      - 仅部署 go-backend"
	@echo "  make all            - 完整流程（清理 + 安装 + 打包）"
	@echo ""

# 配置变量（可根据实际情况修改）
DEPLOY_HOST ?= 101.42.182.96
DEPLOY_USER ?= cg
DEPLOY_PATH ?= /www/wwwroot/yckj-tasks
DEPLOY_SSH_ALIAS ?= yckj
GO_BINARY_NAME ?= task-system-go
# 安装所有依赖
install: install-admin install-user install-go

install-admin:
	@echo "==> 安装 admin 前端依赖..."
	cd admin && npm install

install-user:
	@echo "==> 安装 user 前端依赖..."
	cd user && npm install

install-go:
	@echo "==> 安装 go-backend 依赖..."
	cd go-backend && go mod download

# 打包所有项目
build: build-admin build-user build-go
	@echo "==> ✅ 所有项目打包完成！"

build-admin:
	@echo "==> 打包 admin 前端..."
	cd admin && npm run build
	@echo "==> ✅ admin 打包完成 -> admin/dist"

build-user:
	@echo "==> 打包 user 前端..."
	cd user && npm run build
	@echo "==> ✅ user 打包完成 -> user/dist"

build-go:
	@echo "==> 编译 go-backend..."
	cd go-backend && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/$(GO_BINARY_NAME) ./cmd/server
	@echo "==> ✅ go-backend 编译完成 -> go-backend/bin/$(GO_BINARY_NAME)"

# 清理构建产物
clean:
	@echo "==> 清理构建产物..."
	rm -rf admin/dist
	rm -rf user/dist
	rm -f go-backend/$(GO_BINARY_NAME)
	@echo "==> ✅ 清理完成"

# 本地测试
test-admin:
	@echo "==> 启动 admin 开发服务器..."
	cd admin && npm run dev

test-user:
	@echo "==> 启动 user 开发服务器..."
	cd user && npm run dev

test-go:
	@echo "==> 启动 go-backend 服务..."
	cd go-backend && go run ./cmd/server/main.go

# 部署到服务器（需要配置 SSH）
deploy: build deploy-files restart-services
	@echo "==> ✅ 所有项目部署完成！"

deploy-files:
	@echo "==> 上传文件到服务器..."
	ssh $(DEPLOY_SSH_ALIAS) "mkdir -p $(DEPLOY_PATH)/{admin,user,go-backend}"
	rsync -avz --delete admin/dist/ $(DEPLOY_SSH_ALIAS):$(DEPLOY_PATH)/admin/
	rsync -avz --delete user/dist/ $(DEPLOY_SSH_ALIAS):$(DEPLOY_PATH)/user/
	rsync -avz go-backend/bin/$(GO_BINARY_NAME) $(DEPLOY_SSH_ALIAS):/home/cg/yckj-tasks/

deploy-admin: build-admin
	@echo "==> 部署 admin 到服务器..."
	rsync -avz --delete admin/dist/ $(DEPLOY_SSH_ALIAS):$(DEPLOY_PATH)/admin/

deploy-user: build-user
	@echo "==> 部署 user 到服务器..."
	rsync -avz --delete user/dist/ $(DEPLOY_SSH_ALIAS):$(DEPLOY_PATH)/user/

deploy-go: build-go
	@echo "==> 部署 go-backend 到服务器..."
	rsync -avz go-backend/bin/$(GO_BINARY_NAME) $(DEPLOY_SSH_ALIAS):/home/cg/yckj-tasks/
	ssh $(DEPLOY_SSH_ALIAS) "chmod +x /home/cg/yckj-tasks/$(GO_BINARY_NAME)"

restart-services:
	@echo "==> 重启服务..."
	ssh $(DEPLOY_SSH_ALIAS) "sudo supervisorctl restart task-system-go"

# 打包成压缩包（用于手动部署）
package: build
	@echo "==> 创建部署压缩包..."
	mkdir -p dist-package
	cp -r admin/dist dist-package/admin
	cp -r user/dist dist-package/user
	mkdir -p dist-package/go-backend
	cp go-backend/$(GO_BINARY_NAME) dist-package/go-backend/
	tar -czf yckj-tasks-$(shell date +%Y%m%d-%H%M%S).tar.gz dist-package/
	rm -rf dist-package
	@echo "==> ✅ 压缩包创建完成"

# 完整流程
all: clean install build
	@echo "==> ✅ 完整打包流程完成！"

# 检查环境
check:
	@echo "==> 检查构建环境..."
	@which node > /dev/null || echo "❌ Node.js 未安装"
	@which npm > /dev/null || echo "❌ npm 未安装"
	@which go > /dev/null || echo "❌ Go 未安装"
	@node --version 2>/dev/null || true
	@npm --version 2>/dev/null || true
	@go version 2>/dev/null || true
	@echo "==> ✅ 环境检查完成"
