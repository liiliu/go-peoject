#!/bin/bash

#===============================================================================
# 生产环境部署脚本 - yun-prod.sh
# 用途：编译、打包并部署 Go 应用到生产服务器
# 作者：Your Team
# 版本：v2.0
# 更新日期：2025-10-16
#===============================================================================

# ============ 配置区域 ============

# 服务器配置（生产环境）
SERVER_HOST="your-prod-server-ip"  # 修改为生产服务器 IP
SERVER_PORT=""                      # 生产服务器 SSH 端口
SERVER_USER="root"
SERVER_PATH="/data/server/your_project_prod"  # 修改为生产部署目录

# 编译配置
BUILD_OUTPUT="server"
BUILD_DIR="deploy"

# 环境配置文件
CONFIG_ENV="prod"

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# ============ 工具函数 ============

# 构建 SSH/SCP 端口参数
function get_port_param() {
    local port=$1
    if [ -z "$port" ] || [ "$port" = "22" ]; then
        echo ""
    else
        echo "-P $port"
    fi
}

function log_info() {
    echo -e "${GREEN}[INFO]${NC} $(date '+%Y-%m-%d %H:%M:%S') - $1"
}

function log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $(date '+%Y-%m-%d %H:%M:%S') - $1"
}

function log_error() {
    echo -e "${RED}[ERROR]${NC} $(date '+%Y-%m-%d %H:%M:%S') - $1"
}

function log_step() {
    echo -e "${BLUE}[STEP]${NC} $1"
}

function check_command() {
    if ! command -v $1 &> /dev/null; then
        log_error "命令 $1 未找到，请先安装"
        exit 1
    fi
}

function cleanup() {
    log_info "清理临时文件..."
    cd ..
    rm -rf "${BUILD_DIR}"
    log_info "清理完成"
}

function handle_error() {
    log_error "部署过程中出现错误"
    cleanup
    exit 1
}

trap handle_error ERR

# ============ 主程序 ============

echo ""
echo "========================================="
echo "  ⚠️  部署到生产环境 (PRODUCTION)  ⚠️"
echo "========================================="
echo ""
log_warn "警告：即将部署到生产环境！"
echo ""
echo "请确认以下检查项："
echo "  ✓ 已修改 config/env/prod.toml 中的所有配置"
echo "  ✓ 数据库连接信息正确"
echo "  ✓ Redis 密码已设置"
echo "  ✓ JWT Secret 已修改为强密钥"
echo "  ✓ debug 模式已关闭 (debug = false)"
echo "  ✓ autoMigrate 已关闭 (autoMigrate = false)"
echo "  ✓ 本脚本中的 SERVER_HOST 和 SERVER_PATH 已配置"
echo ""
read -p "确认部署到生产环境？(输入 YES 继续): " confirm

if [ "$confirm" != "YES" ]; then
    echo ""
    log_warn "已取消部署"
    exit 0
fi

echo ""
log_info "开始部署到生产环境..."
log_info "目标环境: 生产环境 (${CONFIG_ENV})"

# 检查配置文件
if [ ! -f "config/env/prod.toml" ]; then
    log_error "生产环境配置文件不存在: config/env/prod.toml"
    exit 1
fi

# 检查配置文件中是否还有默认值
log_step "检查配置安全性"
if grep -q "CHANGE-THIS" "config/env/prod.toml"; then
    log_error "生产配置中仍有未修改的默认值 (CHANGE-THIS)"
    log_error "请先修改 config/env/prod.toml 中的所有敏感配置"
    exit 1
fi

if grep -q "your-" "config/env/prod.toml"; then
    log_warn "配置文件中发现 'your-' 开头的值，请确认是否已修改"
    read -p "继续？(y/n): " cont
    if [ "$cont" != "y" ]; then
        exit 0
    fi
fi

log_info "配置检查通过"

# 检查必要的命令
log_step "1/8 检查依赖"
check_command "go"
check_command "tar"
check_command "scp"
check_command "ssh"
log_info "依赖检查通过"

# 清理旧的构建目录
log_step "2/8 清理旧构建"
if [ -d "${BUILD_DIR}" ]; then
    log_warn "删除旧的构建目录: ${BUILD_DIR}"
    rm -rf "${BUILD_DIR}"
fi

# 创建构建目录
log_step "3/8 创建构建目录"
mkdir -p "${BUILD_DIR}/config"
log_info "构建目录已创建: ${BUILD_DIR}"

# 编译 Go 程序（Linux 64位）
log_step "4/8 编译 Go 程序"
log_info "编译目标: Linux/amd64"
log_info "输出文件: ${BUILD_DIR}/${BUILD_OUTPUT}"

if ! env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o "${BUILD_DIR}/${BUILD_OUTPUT}" -a -installsuffix cgo -ldflags="-s -w" .; then
    log_error "编译失败"
    exit 1
fi

log_info "编译成功 (大小: $(du -h "${BUILD_DIR}/${BUILD_OUTPUT}" | cut -f1))"

# 复制配置文件
log_step "5/8 复制配置文件"
cp -f "config/env/${CONFIG_ENV}.toml" "${BUILD_DIR}/config/config.toml"
log_info "配置文件已复制: ${CONFIG_ENV}.toml -> config.toml"

# 打包
log_step "6/8 打包部署文件"
cd "${BUILD_DIR}"
if ! tar -czf server.tgz *; then
    log_error "打包失败"
    cd ..
    exit 1
fi
log_info "打包完成: server.tgz (大小: $(du -h server.tgz | cut -f1))"

# 上传到服务器
log_step "7/8 上传到服务器"
log_info "目标服务器: ${SERVER_USER}@${SERVER_HOST}:${SERVER_PORT}"
log_info "目标路径: ${SERVER_PATH}"

if ! scp $(get_port_param "${SERVER_PORT}") server.tgz "${SERVER_USER}@${SERVER_HOST}:${SERVER_PATH}/server.tgz"; then
    log_error "上传失败"
    cd ..
    exit 1
fi
log_info "上传成功"

# 远程执行重启脚本
log_step "8/8 远程重启服务"
log_info "正在远程执行重启脚本..."

if ! ssh $(get_port_param "${SERVER_PORT}") "${SERVER_USER}@${SERVER_HOST}" "cd ${SERVER_PATH} && bash restart.sh start"; then
    log_error "远程重启失败"
    cd ..
    exit 1
fi

log_info "远程重启完成"

# 清理临时文件
cd ..
rm -rf "${BUILD_DIR}"

echo ""
log_info "==============================================="
log_info "🎉 生产环境部署成功完成！"
log_info "服务器: ${SERVER_HOST}"
log_info "环境: 生产环境 (${CONFIG_ENV})"
log_info "==============================================="
echo ""
