# FX Data Analysis

FXの売買データを管理・分析するためのフルスタックアプリケーションです。GolangのREST APIサーバーとReactのSPAで構築されています。

## 機能

- FX取引データのCRUD操作
  - 売買時期（取引日時）
  - ロット数
  - 購入レート
- リアルタイムデータ表示
- 直感的なWebインターフェース

## 技術スタック

### バックエンド
- Go 1.23以上
- Gin Web Framework
- SQLite3

### フロントエンド
- React 18
- TypeScript
- Vite
- Axios

## 必要な環境

- Go 1.21以上
- Node.js 18以上
- SQLite3

## セットアップ

### 1. バックエンドのセットアップ

```bash
# 依存パッケージのインストール
go mod tidy

# サーバーの起動
go run cmd/server/main.go
```

サーバーは `http://localhost:8080` で起動します。

### 2. フロントエンドのセットアップ

```bash
# フロントエンドディレクトリに移動
cd frontend

# 依存パッケージのインストール
npm install

# 開発サーバーの起動
npm run dev
```

フロントエンドは `http://localhost:3000` で起動します。

### 3. アプリケーションの起動

ブラウザで `http://localhost:3000` を開くと、FX取引データ管理アプリケーションが表示されます。

### 4. 1コマンドで起動（推奨）

リポジトリ直下で以下を実行すると、バックエンドとフロントエンドが同時に起動し、GUIが表示されます。

```bash
make up
```

## レート(USDJPY)の取得と保存

以下のコマンドで、`https://navi.gaitame.com/v3/info/prices/rate` から `pair=USDJPY` のデータを取得し、`data/usdjpy/日付.json` に保存します（ファイル名はローカル日付 `YYYY-MM-DD`）。

```bash
go run ./cmd/fetch-usdjpy
```

任意の日付名で保存したい場合は `-date` を指定できます。

```bash
go run ./cmd/fetch-usdjpy -date 2026-03-01
```

### シングルバイナリとしてビルド

`cmd/fetch-usdjpy` は単体の実行ファイル（シングルバイナリ）として出力できます。

```bash
make build-fetch-usdjpy
./bin/fetch-usdjpy
```

Makeを使わない場合は以下でもOKです。

```bash
mkdir -p bin
CGO_ENABLED=0 go build -trimpath -ldflags "-s -w" -o bin/fetch-usdjpy ./cmd/fetch-usdjpy
```

## API エンドポイント

### 1. 取引データの作成

新しい取引データを作成します。

**エンドポイント:** `POST /api/v1/trades`

**リクエストボディ:**
```json
{
  "trade_time": "2025-10-16T10:30:00Z",
  "lot_size": 1.5,
  "purchase_rate": 150.25
}
```

**レスポンス例:**
```json
{
  "id": 1,
  "trade_time": "2025-10-16T10:30:00Z",
  "lot_size": 1.5,
  "purchase_rate": 150.25,
  "created_at": "2025-10-16T10:30:00Z",
  "updated_at": "2025-10-16T10:30:00Z"
}
```

**curlコマンド例:**
```bash
curl -X POST http://localhost:8080/api/v1/trades \
  -H "Content-Type: application/json" \
  -d '{
    "trade_time": "2025-10-16T10:30:00Z",
    "lot_size": 1.5,
    "purchase_rate": 150.25
  }'
```

### 2. 全取引データの取得

すべての取引データを取得します（取引日時の降順）。

**エンドポイント:** `GET /api/v1/trades`

**レスポンス例:**
```json
[
  {
    "id": 1,
    "trade_time": "2025-10-16T10:30:00Z",
    "lot_size": 1.5,
    "purchase_rate": 150.25,
    "created_at": "2025-10-16T10:30:00Z",
    "updated_at": "2025-10-16T10:30:00Z"
  }
]
```

**curlコマンド例:**
```bash
curl http://localhost:8080/api/v1/trades
```

### 3. 特定の取引データの取得

IDを指定して取引データを取得します。

**エンドポイント:** `GET /api/v1/trades/:id`

**レスポンス例:**
```json
{
  "id": 1,
  "trade_time": "2025-10-16T10:30:00Z",
  "lot_size": 1.5,
  "purchase_rate": 150.25,
  "created_at": "2025-10-16T10:30:00Z",
  "updated_at": "2025-10-16T10:30:00Z"
}
```

**curlコマンド例:**
```bash
curl http://localhost:8080/api/v1/trades/1
```

### 4. 取引データの更新

IDを指定して取引データを更新します。更新したいフィールドのみを送信できます。

**エンドポイント:** `PUT /api/v1/trades/:id`

**リクエストボディ:**
```json
{
  "lot_size": 2.0,
  "purchase_rate": 151.00
}
```

**レスポンス例:**
```json
{
  "id": 1,
  "trade_time": "2025-10-16T10:30:00Z",
  "lot_size": 2.0,
  "purchase_rate": 151.00,
  "created_at": "2025-10-16T10:30:00Z",
  "updated_at": "2025-10-16T11:00:00Z"
}
```

**curlコマンド例:**
```bash
curl -X PUT http://localhost:8080/api/v1/trades/1 \
  -H "Content-Type: application/json" \
  -d '{
    "lot_size": 2.0,
    "purchase_rate": 151.00
  }'
```

### 5. 取引データの削除

IDを指定して取引データを削除します。

**エンドポイント:** `DELETE /api/v1/trades/:id`

**レスポンス例:**
```json
{
  "message": "trade deleted successfully"
}
```

**curlコマンド例:**
```bash
curl -X DELETE http://localhost:8080/api/v1/trades/1
```

## プロジェクト構造

```
fx-data-analysis/
├── cmd/
│   └── server/
│       └── main.go              # バックエンドエントリーポイント
├── internal/
│   ├── database/
│   │   └── database.go          # データベース操作
│   ├── handler/
│   │   └── handler.go           # HTTPハンドラー
│   └── model/
│       └── trade.go             # データモデル
├── frontend/
│   ├── src/
│   │   ├── api/
│   │   │   └── trades.ts        # APIクライアント
│   │   ├── components/
│   │   │   ├── TradeForm.tsx    # 取引フォーム
│   │   │   └── TradeList.tsx    # 取引一覧
│   │   ├── types/
│   │   │   └── trade.ts         # 型定義
│   │   ├── App.tsx              # メインアプリ
│   │   ├── main.tsx             # フロントエンドエントリーポイント
│   │   └── index.css            # グローバルスタイル
│   ├── index.html
│   ├── package.json
│   ├── tsconfig.json
│   ├── vite.config.ts
│   └── README.md
├── go.mod
├── go.sum
├── .gitignore
└── README.md
```

## データベース

SQLiteデータベース（`fx_trades.db`）はサーバー起動時に自動的に作成されます。

### テーブル構造: trades

| カラム名 | 型 | 説明 |
|---------|-----|------|
| id | INTEGER | 主キー（自動採番） |
| trade_time | DATETIME | 取引日時 |
| lot_size | REAL | ロット数 |
| purchase_rate | REAL | 購入レート |
| created_at | DATETIME | 作成日時 |
| updated_at | DATETIME | 更新日時 |

## スクリーンショット

### 取引データ管理画面
- 新規取引の作成フォーム
- 取引データの一覧表示
- 編集・削除機能

## 開発

### フロントエンドのビルド

```bash
cd frontend
npm run build
```

ビルド成果物は `frontend/dist` ディレクトリに生成されます。

### バックエンドのビルド

```bash
go build -o fx-server cmd/server/main.go
```

## トラブルシューティング

### CORSエラーが発生する場合
- バックエンドサーバーが起動していることを確認してください
- `cmd/server/main.go` のCORS設定で、フロントエンドのURLが許可されていることを確認してください

### データベースエラーが発生する場合
- `fx_trades.db` ファイルの権限を確認してください
- データベースファイルを削除して、サーバーを再起動すると新しいデータベースが作成されます

## ライセンス

MIT
