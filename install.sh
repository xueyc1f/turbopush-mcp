#!/usr/bin/env sh
# TurboPush MCP 一键安装脚本
#
# 用法:
#   curl -fsSL https://raw.githubusercontent.com/xueyc1f/turbopush-mcp/main/install.sh | sh
#
# 可选环境变量:
#   VERSION      指定版本，如 v1.1.1；默认最新 release
#   INSTALL_DIR  安装目录，默认 $HOME/.local/bin
#   SKIP_MCP_ADD 设为 1 则跳过自动执行 `claude mcp add`

set -eu

REPO="xueyc1f/turbopush-mcp"
BIN_NAME="turbo-push-mcp"
INSTALL_DIR="${INSTALL_DIR:-$HOME/.local/bin}"
VERSION="${VERSION:-}"

# ---- 检测 OS / ARCH ----
OS="$(uname -s)"
ARCH_RAW="$(uname -m)"

case "$OS" in
  Darwin) PLATFORM="apple-darwin" ;;
  Linux)  PLATFORM="unknown-linux-gnu" ;;
  MINGW*|MSYS*|CYGWIN*)
    PLATFORM="pc-windows-msvc"
    BIN_NAME="turbo-push-mcp.exe"
    ;;
  *) echo "不支持的操作系统: $OS" >&2; exit 1 ;;
esac

case "$ARCH_RAW" in
  x86_64|amd64) ARCH="x86_64" ;;
  arm64|aarch64) ARCH="aarch64" ;;
  *) echo "不支持的架构: $ARCH_RAW" >&2; exit 1 ;;
esac

ASSET="turbo-push-mcp-${ARCH}-${PLATFORM}"
[ "$PLATFORM" = "pc-windows-msvc" ] && ASSET="${ASSET}.exe"

# ---- 解析版本 ----
if [ -z "$VERSION" ]; then
  VERSION="$(curl -fsSL "https://api.github.com/repos/${REPO}/releases/latest" \
    | sed -n 's/.*"tag_name"[[:space:]]*:[[:space:]]*"\([^"]*\)".*/\1/p' | head -n1)"
  [ -z "$VERSION" ] && { echo "无法获取最新版本" >&2; exit 1; }
fi

URL="https://github.com/${REPO}/releases/download/${VERSION}/${ASSET}"
echo "→ 下载 ${ASSET} (${VERSION})"
echo "  ${URL}"

mkdir -p "$INSTALL_DIR"
TARGET="${INSTALL_DIR}/${BIN_NAME}"

if ! curl -fSL --progress-bar -o "$TARGET" "$URL"; then
  echo "下载失败: $URL" >&2
  exit 1
fi
chmod +x "$TARGET"

echo "✓ 已安装到 $TARGET"

# ---- 自动注册到 Claude Code ----
if [ "${SKIP_MCP_ADD:-0}" != "1" ] && command -v claude >/dev/null 2>&1; then
  echo "→ 注册到 Claude Code (claude mcp add turbo-push)"
  if claude mcp add turbo-push -- "$TARGET" >/dev/null 2>&1; then
    echo "✓ 已注册，重启 Claude Code 生效"
  else
    echo "! 注册失败或已存在，可手动执行: claude mcp add turbo-push -- \"$TARGET\""
  fi
else
  echo ""
  echo "下一步："
  echo "  claude mcp add turbo-push -- \"$TARGET\""
fi

# ---- PATH 提示 ----
case ":$PATH:" in
  *":$INSTALL_DIR:"*) ;;
  *)
    echo ""
    echo "提示: $INSTALL_DIR 不在 PATH 中，可在 shell 配置中加："
    echo "  export PATH=\"$INSTALL_DIR:\$PATH\""
    ;;
esac
