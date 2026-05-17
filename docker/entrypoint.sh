#!/bin/sh
set -eu

export BBSGO_SERVER_URL="${BBSGO_SERVER_URL:-http://127.0.0.1:8082}"

mkdir -p /app/data /app/logs /app/res/uploads
if [ ! -f /app/data/bbs-go.yaml ]; then
	cp /app/defaults/bbs-go.yaml /app/data/bbs-go.yaml
fi

./bbs-go &
api_pid=$!

node scripts/serve-ssr.mjs &
web_pid=$!

terminate() {
	kill "$api_pid" "$web_pid" 2>/dev/null || true
	wait "$api_pid" "$web_pid" 2>/dev/null || true
}

trap 'terminate; exit 143' INT TERM

while kill -0 "$api_pid" 2>/dev/null && kill -0 "$web_pid" 2>/dev/null; do
	sleep 1
done

status=0
if ! kill -0 "$api_pid" 2>/dev/null; then
	wait "$api_pid" || status=$?
else
	wait "$web_pid" || status=$?
fi

terminate
exit "$status"
