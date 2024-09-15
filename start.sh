#!/bin/sh

export APP_HOME=/app/bbs-go

# 启动 bbs-go-server
cd ${APP_HOME}/server
${APP_HOME}/server/bbs-go &

# 启动 bbs-go-site
cd ${APP_HOME}/site
node ${APP_HOME}/site/.output/server/index.mjs &

# 保持容器运行
wait
