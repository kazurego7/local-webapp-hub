#!/usr/bin/env bash
set -euo pipefail

VERSION="${VERSION:-latest}"

usage() {
  cat <<USAGE
Usage: bash install.sh [-v VERSION]

Options:
  -v, --version VERSION   インストールするバージョン（既定: latest）
                          例: v0.1.3 / latest / main / <commit>

環境変数:
  VERSION                 上と同じ意味（指定がない場合は latest）
USAGE
}

while [[ $# -gt 0 ]]; do
  case "${1:-}" in
    -v|--version)
      [[ $# -ge 2 ]] || { echo "--version に値が必要です" >&2; exit 1; }
      VERSION="$2"; shift 2 ;;
    -h|--help)
      usage; exit 0 ;;
    *)
      echo "未知のオプション: $1" >&2; usage; exit 1 ;;
  esac
done

# 概要（ユーザーsystemd版、フラグなし・最小）：
# - go install でバイナリ取得（@latest）して ~/.local/bin に配置
# - ユーザーsystemdユニット (~/.config/systemd/user) を作成
# - 有効化して起動（ログイン外でも動かすために linger も有効化）
# 使い方:
#   bash install.sh

if [[ "$EUID" -eq 0 ]]; then
  echo "このスクリプトは一般ユーザーで実行してください（sudo不要）" >&2
  exit 1
fi

command -v go >/dev/null 2>&1 || { echo "go が見つかりません。インストールしてください。" >&2; exit 1; }

BIN_DIR="${GOBIN:-}"
if [[ -z "$BIN_DIR" ]]; then
  BIN_DIR="$HOME/.local/bin"
fi
mkdir -p "$BIN_DIR"

echo "[1/3] go install でバイナリ取得（出力先: $BIN_DIR, バージョン: $VERSION)"
GOBIN="$BIN_DIR" go install "github.com/kazurego7/local-webapp-hub/cmd/local-webapp-hub@${VERSION}"

SERVICE_DIR="$HOME/.config/systemd/user"
SERVICE_PATH="$SERVICE_DIR/local-webapp-hub.service"
mkdir -p "$SERVICE_DIR"

echo "[2/3] ユーザーsystemdユニットを作成: $SERVICE_PATH"
cat >"$SERVICE_PATH" <<'UNIT'
[Unit]
Description=Local Web App Hub (User)
After=default.target

[Service]
Type=simple
ExecStart=%h/.local/bin/local-webapp-hub -addr :8787
Restart=on-failure
RestartSec=2s

[Install]
WantedBy=default.target
UNIT

echo "[3/3] 有効化・起動（ログイン外でも稼働するよう linger を有効化）"
systemctl --user daemon-reload
systemctl --user enable --now local-webapp-hub.service

if command -v loginctl >/dev/null 2>&1; then
  if ! loginctl show-user "$USER" 2>/dev/null | grep -q '^Linger=yes'; then
    echo "sudo権限で linger を有効化します（ログイン外でも起動維持）"
    sudo loginctl enable-linger "$USER" || true
  fi
fi

echo
echo "インストール完了"
echo "- 実行ファイル: $BIN_DIR/local-webapp-hub"
echo "- ユニット:     $SERVICE_PATH"
echo "- 起動確認:     systemctl --user status local-webapp-hub -n 20"
echo "- 既に稼働中の更新反映: systemctl --user restart local-webapp-hub.service"
