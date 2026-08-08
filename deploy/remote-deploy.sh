#!/usr/bin/env bash
#
# bbs-go 远程部署脚本（服务器侧执行）
#
# 前提：
#   - CI（GitHub Actions）已将最新镜像推送到 GHCR：
#     ghcr.io/aishangwuji/bbs-go:latest
#   - 服务器已有 /opt/bbs-go/docker-compose.yml（image 指向 GHCR）
#
# 用法：
#   bash /opt/bbs-go/repo/deploy/remote-deploy.sh
#
# 该脚本在服务器上执行（不在 CI 中），由部署者手动触发，保证可控可回滚。

set -euo pipefail

DEPLOY_DIR=${BBSGO_DEPLOY_DIR:-/opt/bbs-go}
COMPOSE_FILE="${DEPLOY_DIR}/docker-compose.yml"
IMAGE=ghcr.io/aishangwuji/bbs-go:latest

log() {
  echo "[deploy] $(date '+%Y-%m-%d %H:%M:%S') $*"
}

if [[ ! -f "${COMPOSE_FILE}" ]]; then
  echo "错误：未找到 ${COMPOSE_FILE}，请检查 BBSGO_DEPLOY_DIR" >&2
  exit 1
fi

if ! command -v docker >/dev/null 2>&1; then
  echo "错误：未安装 docker" >&2
  exit 1
fi

cd "${DEPLOY_DIR}"

log "开始部署 bbs-go（${IMAGE}）"

# 1. 拉取最新镜像
log "拉取镜像：${IMAGE}"
docker pull "${IMAGE}"

# 2. 校验 compose 文件
log "校验 docker compose 配置"
docker compose -f "${COMPOSE_FILE}" config --quiet

# 3. 滚动重建容器（仅当镜像变化时重建）
log "重建容器"
docker compose -f "${COMPOSE_FILE}" up -d --no-deps

# 4. 健康检查（最多等 90 秒）
log "等待健康检查（最多 90 秒）"
for i in $(seq 1 45); do
  STATUS=$(docker inspect --format '{{if .State.Health}}{{.State.Health.Status}}{{else}}{{.State.Status}}{{end}}' bbs-go 2>/dev/null || echo "unknown")
  if [[ "${STATUS}" == "healthy" ]]; then
    log "部署完成，bbs-go 状态 healthy"
    exit 0
  fi
  if [[ "${STATUS}" != "starting" && "${STATUS}" != "running" && "${STATUS}" != "unknown" ]]; then
    log "bbs-go 状态异常：${STATUS}"
    docker logs bbs-go --tail 50 || true
    exit 1
  fi
  sleep 2
done

log "健康检查超时，请查看容器日志：docker logs bbs-go --tail 200"
exit 1
