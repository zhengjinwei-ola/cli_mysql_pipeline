#!/bin/bash
set -e

REMOTE_HOST=${1:?"Usage: deploy.sh <host>"}
REMOTE_PATH="/opt/cli_mysql_pipeline"

echo "=== Deploying cli_mysql_pipeline to ${REMOTE_HOST}:${REMOTE_PATH} ==="

# 先编译
bash scripts/build.sh

# 上传二进制和配置
ssh "${REMOTE_HOST}" "mkdir -p ${REMOTE_PATH} /var/log/cli_mysql_pipeline"
rsync -avz dist/ "${REMOTE_HOST}:${REMOTE_PATH}/"

# 安装 supervisor 配置（首次部署时需要）
scp scripts/supervisor.conf "${REMOTE_HOST}:/etc/supervisor/conf.d/cli_mysql_pipeline.conf"
ssh "${REMOTE_HOST}" "supervisorctl reread && supervisorctl update"

# 重启
ssh "${REMOTE_HOST}" "supervisorctl restart cli_mysql_pipeline"

echo "=== Deploy complete ==="
