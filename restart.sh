#!/bin/bash

#===============================================================================
# 服务器端部署脚本 - restart.sh
# 用途：管理 Go 应用服务的启动、停止、重启和监控
# 作者：Your Team
# 版本：v2.0
# 更新日期：2025-10-16
#===============================================================================

# ============ 配置区域 ============

# 项目部署目录（请根据实际情况修改）
PROJECT_DIR="/data/server/ddt_v4/deploy_operation_backend_server"

# 服务二进制文件名
SERVER_BIN="server"

# 需要启动的服务列表（可根据需要调整）
# 可选值：api, job, worker, grpc, websocket
SERVICES=("api" "job")

# 日志文件大小限制（字节）
# 当日志文件超过此大小时会被清空
# 默认：10MB = 10485760 字节
LOG_SIZE_LIMIT=10485760

# 监控检查间隔（秒）
WATCH_INTERVAL=60

# ============ 颜色输出 ============
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# ============ 工具函数 ============

# 打印信息
function log_info() {
    echo -e "${GREEN}[INFO]${NC} $(date '+%Y-%m-%d %H:%M:%S') - $1"
}

# 打印警告
function log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $(date '+%Y-%m-%d %H:%M:%S') - $1"
}

# 打印错误
function log_error() {
    echo -e "${RED}[ERROR]${NC} $(date '+%Y-%m-%d %H:%M:%S') - $1"
}

# 打印服务信息
function log_service() {
    echo -e "${BLUE}[SERVICE]${NC} $(date '+%Y-%m-%d %H:%M:%S') - $1"
}

# ============ 服务管理函数 ============

# 停止所有服务
function stop() {
    log_info "正在停止所有服务..."
    
    # 获取所有运行中的服务进程
    local pids=$(ps -ef | grep "${PROJECT_DIR}/${SERVER_BIN}" | grep -v grep | awk '{print $2}')
    
    if [ -z "$pids" ]; then
        log_warn "没有找到运行中的服务"
        return 0
    fi
    
    # 优雅停止：先发送 SIGTERM
    log_info "发送停止信号 (SIGTERM)..."
    echo "$pids" | xargs kill -15 2>/dev/null
    
    # 等待进程退出
    sleep 3
    
    # 检查是否还有进程存在
    local remaining=$(ps -ef | grep "${PROJECT_DIR}/${SERVER_BIN}" | grep -v grep | wc -l)
    
    if [ $remaining -gt 0 ]; then
        log_warn "部分进程未响应，强制终止 (SIGKILL)..."
        ps -ef | grep "${PROJECT_DIR}/${SERVER_BIN}" | grep -v grep | awk '{print $2}' | xargs kill -9 2>/dev/null
        sleep 1
    fi
    
    log_info "所有服务已停止"
}

# 停止指定服务
function stop_service() {
    local service_name=$1
    log_info "正在停止服务: ${service_name}"
    
    local pids=$(ps -ef | grep "${PROJECT_DIR}/${SERVER_BIN} ${service_name}" | grep -v grep | awk '{print $2}')
    
    if [ -z "$pids" ]; then
        log_warn "服务 ${service_name} 未运行"
        return 0
    fi
    
    echo "$pids" | xargs kill -15 2>/dev/null
    sleep 2
    
    # 检查是否还在运行
    if ps -ef | grep "${PROJECT_DIR}/${SERVER_BIN} ${service_name}" | grep -v grep > /dev/null; then
        log_warn "强制停止服务: ${service_name}"
        ps -ef | grep "${PROJECT_DIR}/${SERVER_BIN} ${service_name}" | grep -v grep | awk '{print $2}' | xargs kill -9 2>/dev/null
    fi
    
    log_info "服务 ${service_name} 已停止"
}

# 解压并部署
function deploy() {
    log_info "开始部署新版本..."
    
    # 检查部署包是否存在
    if [ ! -f "server.tgz" ]; then
        log_error "部署包 server.tgz 不存在"
        exit 1
    fi
    
    # 备份旧版本（如果存在）
    if [ -f "${SERVER_BIN}" ]; then
        log_info "备份旧版本..."
        local backup_dir="backup/$(date '+%Y%m%d_%H%M%S')"
        mkdir -p "$backup_dir"
        cp -f "${SERVER_BIN}" "$backup_dir/" 2>/dev/null || true
        [ -d "config" ] && cp -rf "config" "$backup_dir/" 2>/dev/null || true
        log_info "备份完成: $backup_dir"
    fi
    
    # 清理旧文件
    log_info "清理旧文件..."
    rm -rf conf config ${SERVER_BIN}
    
    # 解压新版本
    log_info "解压部署包..."
    if ! tar -zxvf server.tgz > /dev/null 2>&1; then
        log_error "解压部署包失败"
        exit 1
    fi
    
    # 设置执行权限
    chmod +x ${SERVER_BIN}
    
    log_info "部署完成"
}

# 启动所有服务
function start() {
    log_info "开始启动流程..."
    
    # 先停止旧服务
    stop
    
    # 部署新版本
    deploy
    
    # 启动服务
    run
}

# 运行服务
function run() {
    log_info "正在启动服务..."
    
    # 确保日志目录存在
    mkdir -p logs
    
    # 启动配置的所有服务
    for service in "${SERVICES[@]}"; do
        log_service "启动服务: ${service}"
        nohup ${PROJECT_DIR}/${SERVER_BIN} ${service} >> logs/${service}.log 2>&1 &
        local pid=$!
        
        # 等待服务启动
        sleep 2
        
        # 检查服务是否启动成功
        if ps -p $pid > /dev/null; then
            log_info "服务 ${service} 启动成功 (PID: $pid)"
        else
            log_error "服务 ${service} 启动失败"
        fi
    done
    
    log_info "所有服务启动完成"
}

# 重启所有服务
function restart() {
    log_info "重启所有服务"
    stop
    sleep 2
    run
}

# 查看服务状态
function status() {
    log_info "检查服务状态..."
    echo ""
    
    for service in "${SERVICES[@]}"; do
        local pid=$(ps -ef | grep "${PROJECT_DIR}/${SERVER_BIN} ${service}" | grep -v grep | awk '{print $2}')
        
        if [ -n "$pid" ]; then
            local mem=$(ps -p $pid -o rss= | awk '{printf "%.2f MB", $1/1024}')
            local cpu=$(ps -p $pid -o %cpu= | awk '{printf "%.1f%%", $1}')
            local uptime=$(ps -p $pid -o etime= | xargs)
            echo -e "${GREEN}●${NC} ${service}"
            echo "  ├─ PID: $pid"
            echo "  ├─ 内存: $mem"
            echo "  ├─ CPU: $cpu"
            echo "  └─ 运行时间: $uptime"
        else
            echo -e "${RED}○${NC} ${service} (未运行)"
        fi
        echo ""
    done
}

# 查看日志
function logs() {
    local service=$2
    local lines=${3:-50}
    
    if [ -z "$service" ]; then
        log_error "请指定服务名称"
        echo "用法: $0 logs <service_name> [lines]"
        echo "示例: $0 logs api 100"
        exit 1
    fi
    
    local log_file="logs/${service}.log"
    
    if [ ! -f "$log_file" ]; then
        log_error "日志文件不存在: $log_file"
        exit 1
    fi
    
    log_info "查看服务 ${service} 的最近 ${lines} 行日志:"
    echo ""
    tail -n $lines "$log_file"
}

# 清理日志
function clean_logs() {
    log_info "清理过大的日志文件..."
    
    for service in "${SERVICES[@]}"; do
        local log_file="logs/${service}.log"
        
        if [ -f "$log_file" ]; then
            local log_size=$(stat -c%s "$log_file" 2>/dev/null || stat -f%z "$log_file" 2>/dev/null)
            
            if [ $log_size -gt $LOG_SIZE_LIMIT ]; then
                log_warn "日志文件 ${log_file} 过大 ($(($log_size/1024/1024))MB)，正在清理..."
                # 保留最后 1000 行
                tail -n 1000 "$log_file" > "${log_file}.tmp"
                mv "${log_file}.tmp" "$log_file"
                log_info "日志文件 ${log_file} 已清理"
            fi
        fi
    done
}

# 监控模式
function watch() {
    log_info "进入监控模式 (每 ${WATCH_INTERVAL} 秒检查一次)..."
    log_info "按 Ctrl+C 退出监控"
    
    while true; do
        # 检查服务是否运行
        for service in "${SERVICES[@]}"; do
            if ! ps -ef | grep "${PROJECT_DIR}/${SERVER_BIN} ${service}" | grep -v grep > /dev/null; then
                log_error "检测到服务 ${service} 已停止，正在重启..."
                nohup ${PROJECT_DIR}/${SERVER_BIN} ${service} >> logs/${service}.log 2>&1 &
                log_info "服务 ${service} 已重启"
            fi
        done
        
        # 清理过大的日志
        clean_logs
        
        # 等待下一次检查
        sleep $WATCH_INTERVAL
    done
}

# 显示帮助信息
function show_help() {
    cat << EOF

${GREEN}服务器端部署脚本${NC}

${YELLOW}用法:${NC}
    $0 [命令] [参数]

${YELLOW}命令:${NC}
    ${GREEN}start${NC}           部署并启动所有服务
    ${GREEN}stop${NC}            停止所有服务
    ${GREEN}restart${NC}         重启所有服务
    ${GREEN}status${NC}          查看服务状态
    ${GREEN}run${NC}             仅启动服务（不部署）
    ${GREEN}watch${NC}           监控模式（自动重启异常服务）
    ${GREEN}logs${NC}            查看服务日志
    ${GREEN}clean${NC}           清理过大的日志文件
    ${GREEN}help${NC}            显示此帮助信息

${YELLOW}示例:${NC}
    $0 start              # 部署并启动所有服务
    $0 stop               # 停止所有服务
    $0 status             # 查看服务状态
    $0 logs api 100       # 查看 api 服务最近 100 行日志
    $0 watch              # 进入监控模式

${YELLOW}配置的服务:${NC}
    ${SERVICES[@]}

${YELLOW}部署目录:${NC}
    ${PROJECT_DIR}

EOF
}

# ============ 主程序入口 ============

# 切换到项目目录
cd "$PROJECT_DIR" || {
    log_error "无法切换到项目目录: $PROJECT_DIR"
    exit 1
}

# 根据参数执行相应操作
case $1 in
    start)
        start
        ;;
    stop)
        stop
        ;;
    restart)
        restart
        ;;
    run)
        run
        ;;
    status)
        status
        ;;
    watch)
        watch
        ;;
    logs)
        logs "$@"
        ;;
    clean)
        clean_logs
        ;;
    help|--help|-h)
        show_help
        ;;
    *)
        log_error "未知命令: $1"
        show_help
        exit 1
        ;;
esac

exit 0


