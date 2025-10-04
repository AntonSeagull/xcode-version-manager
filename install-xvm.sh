#!/bin/bash
set -e

APP_NAME="xvm"
REPO="AntonSeagull/xcode-version-manager"
VERSION=${1:-"latest"}
INSTALL_PATH="/usr/local/bin"

echo "🧩 $APP_NAME installer / установщик"

# Detect OS and architecture
ARCH=$(uname -m)
OS=$(uname -s | tr '[:upper:]' '[:lower:]')

# Normalize architecture
if [[ "$ARCH" == "x86_64" ]]; then
  ARCH="amd64"
elif [[ "$ARCH" == "arm64" || "$ARCH" == "aarch64" ]]; then
  ARCH="arm64"
else
  echo "❌ Unsupported architecture: $ARCH"
  echo "⛔️ Неподдерживаемая архитектура: $ARCH"
  exit 1
fi

# Normalize OS
if [[ "$OS" == "darwin" ]]; then
  OS="macos"
elif [[ "$OS" == "linux" ]]; then
  OS="linux"
else
  echo "❌ Unsupported OS: $OS"
  echo "⛔️ Неподдерживаемая ОС: $OS"
  exit 1
fi

# Determine binary name
BINARY="${APP_NAME}-${OS}-${ARCH}"

# Get latest version if not provided
if [ "$VERSION" == "latest" ]; then
  echo "🔍 Fetching latest release version... / Получение последней версии..."
  VERSION=$(curl -s https://api.github.com/repos/$REPO/releases/latest | grep tag_name | cut -d '"' -f 4)
fi

URL="https://github.com/$REPO/releases/download/$VERSION/$BINARY"

echo "⬇️ Downloading $URL ..."
echo "📥 Загрузка бинарника $BINARY версии $VERSION ..."

curl -L "$URL" -o "$APP_NAME"
chmod +x "$APP_NAME"

# Install
echo "📦 Installing to $INSTALL_PATH/$APP_NAME ..."
echo "⚙️ Установка в $INSTALL_PATH/$APP_NAME ..."
sudo mv "$APP_NAME" "$INSTALL_PATH/$APP_NAME"

echo
echo "✅ $APP_NAME version $VERSION installed successfully!"
echo "🎉 Успешно установлено: $APP_NAME версия $VERSION"
echo
echo "➡️ You can now run / Теперь можно запускать:"
echo "   $APP_NAME list"