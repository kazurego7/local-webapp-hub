# Local Web App Hub

ローカルで動作中の HTTP/HTTPS サービス（localhost 上）を自動検出し、一覧表示するシンプルなハブです。リンクをクリックすると対象のポートへ直接アクセスします（プロキシ処理は行いません）。

### 特徴
- ローカルで稼働中の Web アプリを自動検出して一覧表示
- クリックでそのまま対象へアクセス（プロキシなし）
- 単一バイナリ＋ユーザー systemd で常駐可能（root 不要）

### 動作環境
- Linux（ユーザー systemd が利用可能な環境）
- Go 1.23 以上がインストール済み（`install.sh` が `go install` を実行します）

### インストール / 起動
- `bash install.sh`
- 起動確認: `systemctl --user status local-webapp-hub -n 20`
- アクセス: `http://localhost:8787`

### 使い方
- ブラウザで `http://localhost:8787` を開くと、現在ローカルで応答のあるポートが一覧表示されます。
- アクセス時にフルスキャン（1〜65535）を行います。Hub 自身の待受ポートは自動で除外されます。

### 設定（フラグ）
- `-addr`（既定: `:8787`）: Hub の待受アドレス/ポート
  - ユーザー systemd のユニットを編集して変更できます（例: `systemctl --user edit local-webapp-hub.service`）。

### 更新・バージョン固定インストール
- 最新へ更新: `bash install.sh`（内部で `@latest` を使用）
- バージョン指定: `bash install.sh -v vX.Y.Z` または `VERSION=vX.Y.Z bash install.sh`
- 反映（稼働中の再読み込み）: `systemctl --user restart local-webapp-hub.service`

### アンインストール
- `systemctl --user disable --now local-webapp-hub.service`
- `rm ~/.config/systemd/user/local-webapp-hub.service`
- （必要に応じて）`systemctl --user daemon-reload`

### 既知の制約 / 注意
- 生成されるリンクは「このハブにアクセスしたホスト」に合わせたスキーム相対URL（`//host:port/`）です。ブラウザは現在のスキーム（http/https）を自動適用します。
- リモート端末から閲覧しても、表示されたページのホスト（例: Tailscale の 100.x やホスト名）に追従します。

---

## プロジェクト構成
- `cmd/local-webapp-hub/` エントリポイント（main）
- `internal/server/` Echo の初期化とハンドラ、テンプレート埋め込み
- `internal/scan/` ポートスキャンと検出ロジック
- `internal/server/web/` HTML テンプレート（go:embed 対象）


## ソースからのローカル実行
ビルド時にホームディレクトリ配下へ書き込めない環境でも動くよう、`GOCACHE` をカレント配下に設定しています。

```
GOCACHE=$(pwd)/.gocache go run ./cmd/local-webapp-hub
```

起動後、ブラウザで `http://localhost:8787` を開いてください。`Ctrl+C` で安全に停止できます（グレースフルシャットダウン対応）。

## 検出仕様
- 各ポートに対して HTTP/HTTPS の `GET /` を実行し、応答の有無を判定（HTTP クライアントのみ）。
- HTTPS は自己署名証明書を許容（検出のため `InsecureSkipVerify` を使用）。
- HTML の `<title>` を取得できた場合はアプリ名として表示。取得できない場合は URL を表示（goquery 使用）。
