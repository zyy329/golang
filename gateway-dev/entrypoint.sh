#!/bin/sh
set -e

TZ=${TZ:-Asia/Shanghai}
ln -sf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

mkdir -p /gateway/run/logs

# ln -sf /dev/stdout /gateway/run/logs/access.log
# ln -sf /dev/stderr /gateway/run/logs/error.log

crond

# 启动 nginx，使用 exec 让其成为 PID 1，能收到 SIGTERM
exec nginx -p /gateway/run -c /gateway/conf/nginx.conf -g "daemon off;"