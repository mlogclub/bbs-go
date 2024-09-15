#!/bin/sh

# 启动 bbs-go-server
${APP_HOME}/server/bbs-go &

# 启动 bbs-go-site
node ${APP_HOME}/site/.output/server/index.mjs &

# 使用 exec 使容器主进程为最后一个进程
exec "$@"