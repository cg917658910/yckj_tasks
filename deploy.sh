#!/bin/bash

# 任务管理系统 - 打包部署脚本
# 使用方法: ./deploy.sh [选项]

set -e  # 遇到错误立即退出

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 配置变量（可根据实际情况修改）
DEPLOY_HOST="${DEPLOY_HOST:-your-server.com}"
DEPLOY_USER="${DEPLOY_USER:-ubuntu}"
DEPLOY_PATH="${DEPLOY_PATH:-/var/www/yckj-tasks}"
GO_BINARY_NAME="${GO_BINARY_NAME:-task-system-go}"

# 项目根目录
PROJECT_ROOT="$(cd "$(dirname "$0")" && pwd)"

# 日志函数
log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 检查环境
check_env() {
    log_info "检查构建环境..."
    
    if ! command -v node &> /dev/null; then
        log_error "Node.js 未安装"
        exit 1
    fi
    
    if ! command -v npm &> /dev/null; then
        log_error "npm 未安装"
        exit 1
    fi
    
    if ! command -v go &> /dev/null; then
        log_error "Go 未安装"
        exit 1
    fi
    
    log_info "Node.js: $(node --version)"
    log_info "npm: $(npm --version)"
    log_info "Go: $(go version)"
}

# 安装依赖
install_deps() {
    log_info "安装所有项目依赖..."
    
    log_info "安装 admin 前端依赖..."
    cd "$PROJECT_ROOT/admin" && npm install
    
    log_info "安装 user 前端依赖..."
    cd "$PROJECT_ROOT/user" && npm install
    
    log_info "安装 go-backend 依赖..."
    cd "$PROJECT_ROOT/go-backend" && go mod download
    
    log_info "✅ 依赖安装完成"
}

# 打包 admin
build_admin() {
    log_info "打包 admin 前端..."
    cd "$PROJECT_ROOT/admin"
    npm run build
    log_info "✅ admin 打包完成 -> admin/dist"
}

# 打包 user
build_user() {
    log_info "打包 user 前端..."
    cd "$PROJECT_ROOT/user"
    npm run build
    log_info "✅ user 打包完成 -> user/dist"
}

# 编译 go-backend
build_go() {
    log_info "编译 go-backend..."
    cd "$PROJECT_ROOT/go-backend"
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o "$GO_BINARY_NAME" ./cmd/server
    log_info "✅ go-backend 编译完成 -> go-backend/$GO_BINARY_NAME"
}

# 打包所有项目
build_all() {
    log_info "开始打包所有项目..."
    build_admin
    build_user
    build_go
    log_info "✅ 所有项目打包完成！"
}

# 清理构建产物
clean() {
    log_info "清理构建产物..."
    rm -rf "$PROJECT_ROOT/admin/dist"
    rm -rf "$PROJECT_ROOT/user/dist"
    rm -f "$PROJECT_ROOT/go-backend/$GO_BINARY_NAME"
    log_info "✅ 清理完成"
}

# 创建压缩包
package() {
    log_info "创建部署压缩包..."
    
    build_all
    
    cd "$PROJECT_ROOT"
    TIMESTAMP=$(date +%Y%m%d-%H%M%S)
    PACKAGE_NAME="yckj-tasks-$TIMESTAMP.tar.gz"
    
    mkdir -p dist-package
    cp -r admin/dist dist-package/admin
    cp -r user/dist dist-package/user
    mkdir -p dist-package/go-backend
    cp "go-backend/$GO_BINARY_NAME" dist-package/go-backend/
    
    tar -czf "$PACKAGE_NAME" dist-package/
    rm -rf dist-package
    
    log_info "✅ 压缩包创建完成: $PACKAGE_NAME"
}

# 部署到服务器
deploy() {
    log_info "部署到服务器 $DEPLOY_USER@$DEPLOY_HOST..."
    
    build_all
    
    log_info "创建远程目录..."
    ssh "$DEPLOY_USER@$DEPLOY_HOST" "mkdir -p $DEPLOY_PATH/{admin,user,go-backend}"
    
    log_info "上传 admin..."
    rsync -avz --delete "$PROJECT_ROOT/admin/dist/" "$DEPLOY_USER@$DEPLOY_HOST:$DEPLOY_PATH/admin/"
    
    log_info "上传 user..."
    rsync -avz --delete "$PROJECT_ROOT/user/dist/" "$DEPLOY_USER@$DEPLOY_HOST:$DEPLOY_PATH/user/"
    
    log_info "上传 go-backend..."
    rsync -avz "$PROJECT_ROOT/go-backend/$GO_BINARY_NAME" "$DEPLOY_USER@$DEPLOY_HOST:$DEPLOY_PATH/go-backend/"
    ssh "$DEPLOY_USER@$DEPLOY_HOST" "chmod +x $DEPLOY_PATH/go-backend/$GO_BINARY_NAME"
    
    log_info "重启服务..."
    ssh "$DEPLOY_USER@$DEPLOY_HOST" "systemctl restart task-system-go || supervisorctl restart task-system-go || true" || log_warn "服务重启失败，可能需要手动重启"
    
    log_info "✅ 部署完成！"
}

# 显示帮助
show_help() {
    cat << HELP
任务管理系统 - 打包部署脚本

使用方法:
  ./deploy.sh [选项]

选项:
  check           检查构建环境
  install         安装所有项目依赖
  build           打包所有项目（admin + user + go-backend）
  build-admin     仅打包 admin 前端
  build-user      仅打包 user 前端
  build-go        仅编译 go-backend
  clean           清理所有构建产物
  package         创建部署压缩包
  deploy          部署到服务器
  help            显示此帮助信息

环境变量:
  DEPLOY_HOST     部署服务器地址（默认: your-server.com）
  DEPLOY_USER     SSH 用户名（默认: ubuntu）
  DEPLOY_PATH     部署路径（默认: /var/www/yckj-tasks）
  GO_BINARY_NAME  Go 二进制文件名（默认: task-system-go）

示例:
  ./deploy.sh build
  ./deploy.sh package
  DEPLOY_HOST=192.168.1.100 ./deploy.sh deploy

HELP
}

# 主函数
main() {
    case "${1:-help}" in
        check)
            check_env
            ;;
        install)
            check_env
            install_deps
            ;;
        build)
            check_env
            build_all
            ;;
        build-admin)
            check_env
            build_admin
            ;;
        build-user)
            check_env
            build_user
            ;;
        build-go)
            check_env
            build_go
            ;;
        clean)
            clean
            ;;
        package)
            check_env
            package
            ;;
        deploy)
            check_env
            deploy
            ;;
        help|--help|-h)
            show_help
            ;;
        *)
            log_error "未知选项: $1"
            show_help
            exit 1
            ;;
    esac
}

main "$@"
