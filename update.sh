#!/bin/bash

# 定义颜色输出（可选，让提示更清晰）
RED='\033[0;31m'
GREEN='\033[0;32m'
NC='\033[0m' # No Color

# 打印错误信息并退出
error_exit() {
  echo -e "${RED}Error: $1${NC}" >&2
  exit 1
}

# 打印普通提示
info_msg() {
  echo -e "${GREEN}$1${NC}"
}

# 更新前端
update_frontend() {
  info_msg "Updating frontend..."
  cd /root/hami || error_exit "Failed to enter /root/hami directory"
  make build-image DOCKER_IMAGE=projecthami/hami-webui-fe VERSION=dev || error_exit "Frontend build failed"
  UPDATED=1
}

# 更新后端
update_backend() {
  info_msg "Updating backend..."
  cd /root/hami/server || error_exit "Failed to enter /root/hami/server directory"
  make build-image DOCKER_IMAGE=projecthami/hami-webui-be VERSION=dev || error_exit "Backend build failed"
  UPDATED=1
}

# 重启 Kubernetes 服务
restart_service() {
  info_msg "Restarting service with k8s.yml..."
  cd /root/hami || error_exit "Failed to enter /root/hami directory"
  kubectl delete -f k8s.yml || error_exit "Failed to delete k8s resources"
  sleep 2 # 可选：等待资源清理
  kubectl create -f k8s.yml || error_exit "Failed to create k8s resources"
  info_msg "Service restarted successfully."
}

# 显示用法
usage() {
  echo "Usage: $0 {all|frontend|backend}"
  exit 1
}

# 主逻辑
UPDATED=0

case "$1" in
  all)
    update_frontend
    update_backend
    ;;
  frontend)
    update_frontend
    ;;
  backend)
    update_backend
    ;;
  "")
    echo "Error: No argument provided."
    usage
    ;;
  *)
    echo "Error: Invalid argument '$1'."
    usage
    ;;
esac

# 如果执行了任何更新，则重启服务
if [ "$UPDATED" -eq 1 ]; then
  restart_service
fi

info_msg "Update and restart process finished."
