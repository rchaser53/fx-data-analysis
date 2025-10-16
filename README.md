# FX Data Analysis API

FXの売買データを管理・分析するためのREST APIサーバーです。Golangとsqliteで構築されています。

## 機能

- FX取引データのCRUD操作
  - 売買時期（取引日時）
  - ロット数
  - 購入レート

## 必要な環境

- Go 1.21以上
- SQLite3

## セットアップ

### 依存パッケージのインストール

```bash
go mod tidy
```

### サーバーの起動

```bash
go run cmd/server/main.go
```

サーバーは `http://localhost:8080` で起動します。

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
│       └── main.go           # エントリーポイント
├── internal/
│   ├── database/
│   │   └── database.go       # データベース操作
│   ├── handler/
│   │   └── handler.go        # HTTPハンドラー
│   └── model/
│       └── trade.go          # データモデル
├── go.mod
├── go.sum
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

## ライセンス

MIT
