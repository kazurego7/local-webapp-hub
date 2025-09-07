# Local Web App Hub

ローカルで動作中のHTTP/HTTPSサーバ（localhost）を自動検出し、一覧表示するシンプルなHubです。
リンクをクリックすると、対象のポートへそのままアクセスします（プロキシはしません）。

## インストール / 起動（ユーザーsystemd）

最小・推奨の導入は「ユーザーsystemd」での常駐です。root不要、ログイン外でも稼働。

- クイックインストール（推奨）
  - `bash install.sh`
  - これで以下を自動実行します:
  - 確認/ログ: `systemctl --user status local-webapp-hub`

アンインストール（ユーザー）
- `systemctl --user disable --now local-webapp-hub.service`
- `rm ~/.config/systemd/user/local-webapp-hub.service`


注意: 本アプリのリンクは `http(s)://localhost:PORT/` 固定のため、リモートから閲覧すると「閲覧者PCのlocalhost」を指します。同一マシンでの利用（ブラウザ→本アプリ）が前提です。公開用途で正しいホスト名リンクにしたい場合はコード変更が必要です。

## ソースからのローカル実行（開発用途）

ビルド時にホームディレクトリ配下へ書き込めない環境でも動くよう、`GOCACHE` をカレント配下に設定しています。

```
GOCACHE=$(pwd)/.gocache go run .
```

起動後、ブラウザで `http://localhost:8787` を開いてください。
`Ctrl+C` で安全に停止できます（グレースフルシャットダウン対応）。

## 使い方

- アクセス時に常にフルスキャン（`1-65535`）を行います。

## 設定（フラグ）

- `-addr`（既定: `:8787`）
  - Hubの待受アドレス/ポート。

> メモ: Hub自身の待受ポートは走査対象から自動的に除外されます。

## 検出仕様

- 各ポートに対して、HTTP/HTTPSの`GET /`で応答の有無を判定します（HTTPクライアントのみで実施）。
- HTTPSは自己署名証明書を許容します（検出のため `InsecureSkipVerify` を使用）。
- HTMLの`<title>`が取得できた場合はそれをアプリ名として表示し、取得できない場合はURLを表示します（goquery使用）。

## 必要なGoバージョン

- Go 1.23 以上（go.modに記載。ツールチェーンは自動で1.24系へ切替済み）

## HTTP周り（整理点）

- ルーティング: `Echo` を使用（`github.com/labstack/echo/v4`）
- 終了処理: OSシグナルでグレースフルシャットダウン
