#!/bin/bash

# 1. GitHub Repo 信息
REPO="dingdangdog/cloudflare_ddns"
TAG="v0.1.0"  # 这里可以动态获取最新版本，如果需要

# 2. 获取最新发布的版本（如果想动态获取版本）
latest_version=$(curl --silent "https://api.github.com/repos/$REPO/releases/latest" | jq -r .tag_name)

if [ "$latest_version" != "$TAG" ]; then
    echo "new version: $latest_version"
    TAG="$latest_version"
else
    echo "lastest version: $TAG"
fi

# 3. 下载二进制文件和配置文件示例
BINARY_URL="https://github.com/$REPO/releases/download/$TAG/whoisme_server"
CONFIG_URL="https://raw.githubusercontent.com/$REPO/main/ip/config_demo.json"

# 4. 设置下载的文件名
BINARY_FILE="whoisme_server"
CONFIG_FILE="config_demo.json"

# 5. 下载二进制文件
echo "downloading whoisme_server..."
curl -L -o $BINARY_FILE $BINARY_URL

# 6. 下载配置文件示例
echo "downloading config_demo.json..."
curl -L -o $CONFIG_FILE $CONFIG_URL

# 7. 更新完毕，给出提示
echo "download completed."

# 8. 赋予二进制文件执行权限
chmod +x $BINARY_FILE
echo "$BINARY_FILE has been updated."
