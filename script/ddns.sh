#!/bin/bash

# 检查系统平台
OS=$(uname -s)
ARCH=$(uname -m)

# 设置GitHub仓库地址和二进制文件名称
GITHUB_REPO="https://github.com/dingdangdog/cloudflare_ddns/releases/download/latest"
BIN_NAME="cfddns"  # 你的 Go 二进制文件名称
BIN_PATH="/usr/local/bin/$BIN_NAME"  # 安装路径

# 获取最新的版本（你可以根据自己的需要修改为具体的版本号）
LATEST_RELEASE_URL="https://api.github.com/repos/dingdangdog/cloudflare_ddns/releases/latest"
LATEST_VERSION=$(curl -s $LATEST_RELEASE_URL | jq -r .tag_name)

# 获取合适平台的二进制文件URL
if [[ "$OS" == "Darwin" ]]; then
    PLATFORM="darwin"
elif [[ "$OS" == "Linux" ]]; then
    PLATFORM="linux"
else
    echo "Unsupported operating system: $OS"
    exit 1
fi

if [[ "$ARCH" == "x86_64" ]]; then
    ARCHITECTURE="amd64"
else
    echo "Unsupported architecture: $ARCH"
    exit 1
fi

#.tar.gz
EXTENSION=""
# 构建二进制下载URL
DOWNLOAD_URL="$GITHUB_REPO/$LATEST_VERSION/$BIN_NAME-$PLATFORM-$ARCHITECTURE$EXTENSION"

# 下载并解压
echo "Downloading $BIN_NAME from $DOWNLOAD_URL..."
curl -L $DOWNLOAD_URL -o "$BIN_NAME"

# 移动到系统的 bin 目录
echo "Installing $BIN_NAME to $BIN_PATH..."
sudo mv "$BIN_NAME" "$BIN_PATH"
sudo chmod +x "$BIN_PATH"

# 确认安装
if command -v $BIN_NAME &> /dev/null; then
    echo "$BIN_NAME installed successfully."
else
    echo "Installation failed."
fi
