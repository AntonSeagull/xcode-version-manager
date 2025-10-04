#!/bin/bash
set -euo pipefail

APP_NAME="xvm"
REPO="AntonSeagull/xcode-version-manager"
VERSION=${1:-"latest"}
INSTALL_PATH="/usr/local/bin"

echo "🧩 $APP_NAME installer / установщик"

ARCH=$(uname -m)
OS=$(uname -s | tr '[:upper:]' '[:lower:]')

# Normalize arch
case "$ARCH" in
  x86_64) ARCH="amd64" ;;
  arm64|aarch64) ARCH="arm64" ;;
  *) echo "❌ Unsupported arch: $ARCH / Неподдерживаемая архитектура"; exit 1 ;;
esac

# Normalize OS (важно: darwin, а не macos)
case "$OS" in
  darwin) OS="darwin" ;;
  linux)  OS="linux" ;;
  *) echo "❌ Unsupported OS: $OS / Неподдерживаемая ОС"; exit 1 ;;
esac

EXT=""
BINARY="${APP_NAME}-${OS}-${ARCH}${EXT}"

# Resolve version
if [[ "$VERSION" == "latest" ]]; then
  echo "🔍 Fetching latest release... / Получение последней версии..."
  VERSION=$(curl -fsSL "https://api.github.com/repos/$REPO/releases/latest" \
    | grep -m1 '"tag_name":' | cut -d '"' -f 4)
  if [[ -z "$VERSION" ]]; then
    echo "❌ Cannot resolve latest tag / Не удалось получить тег последнего релиза"
    exit 1
  fi
fi

URL="https://github.com/$REPO/releases/download/$VERSION/$BINARY"

echo "⬇️ Downloading $URL ..."
echo "📥 Загрузка бинарника $BINARY версии $VERSION ..."

# Проверка, что ассет существует (HTTP 200)
HTTP_CODE=$(curl -sIL "$URL" | awk '/^HTTP/{code=$2} END{print code}')
if [[ "$HTTP_CODE" != "200" ]]; then
  echo "❌ Asset not found (HTTP $HTTP_CODE): $URL"
  echo "⛔️ Файл релиза не найден. Проверь, что в релизе есть ассет с именем: $BINARY"
  exit 1
fi

# Скачиваем во временный файл
TMP="$(mktemp)"
curl -fsSL "$URL" -o "$TMP"
chmod +x "$TMP"

# На macOS иногда появляется quarantine-бит — снимем его
if [[ "$OS" == "darwin" ]]; then
  xattr -dr com.apple.quarantine "$TMP" 2>/dev/null || true
fi

echo "📦 Installing to $INSTALL_PATH/$APP_NAME ..."
echo "⚙️ Установка в $INSTALL_PATH/$APP_NAME ..."
sudo mv "$TMP" "$INSTALL_PATH/$APP_NAME"

echo
echo "✅ $APP_NAME version $VERSION installed successfully!"
echo "🎉 Успешно установлено: $APP_NAME версия $VERSION"
echo
echo "➡️ You can now run / Теперь можно запускать:"
echo "   $APP_NAME list"