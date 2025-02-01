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

# 3. 检测操作系统和架构
OS=$(uname -s)
ARCH=$(uname -m)

echo "OS: $OS, Arch: $ARCH"

# 4. 下载二进制文件和配置文件示例
BINARY_URL="https://github.com/$REPO/releases/download/$TAG/whoisme_server"
CONFIG_URL="https://raw.githubusercontent.com/$REPO/main/ip/config_demo.json"

if [[ "$OS" == "Linux" && "$ARCH" == "x86_64" ]]; then
    BINARY_URL="https://github.com/$REPO/releases/download/$TAG/whoisme_server-linux-amd64"
elif [[ "$OS" == "Darwin" && "$ARCH" == "x86_64" ]]; then
    BINARY_URL="https://github.com/$REPO/releases/download/$TAG/whoisme_server-darwin-amd64"
elif [[ "$OS" == "Darwin" && "$ARCH" == "arm64" ]]; then
    BINARY_URL="https://github.com/$REPO/releases/download/$TAG/whoisme_server-darwin-arm64"
elif [[ "$OS" == "Linux" && "$ARCH" == "aarch64" ]]; then
    BINARY_URL="https://github.com/$REPO/releases/download/$TAG/whoisme_server-linux-arm64"
elif [[ "$OS" == "Windows" && "$ARCH" == "x86_64" ]]; then
    BINARY_URL="https://github.com/$REPO/releases/download/$TAG/whoisme_server-windows-amd64.exe"
else
    echo "Unsupported platform: $OS $ARCH"
    exit 1
fi

# 5. 设置下载的文件名
BINARY_FILE="whoisme_server"
CONFIG_FILE="config_demo.json"

# 6. 下载二进制文件
echo "downloading whoisme_server..."
curl -L -o $BINARY_FILE $BINARY_URL

# 7. 下载配置文件示例
echo "downloading config_demo.json..."
curl -L -o $CONFIG_FILE $CONFIG_URL

# 8. 更新完毕，给出提示
echo "download completed."

# 9. 赋予二进制文件执行权限
if [[ "$OS" == "Linux" || "$OS" == "Darwin" ]]; then
    chmod +x $BINARY_FILE
    echo "$BINARY_FILE has been updated."
fi