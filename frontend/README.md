# FX Data Analysis - Frontend

このディレクトリには、FX取引データを管理するためのReact SPAが含まれています。

## 技術スタック

- **React 18** - UIライブラリ
- **TypeScript** - 型安全な開発
- **Vite** - 高速ビルドツール
- **Axios** - HTTP通信ライブラリ

## セットアップ

### 依存パッケージのインストール

```bash
cd frontend
npm install
```

### 開発サーバーの起動

```bash
npm run dev
```

開発サーバーは `http://localhost:3000` で起動します。

### ビルド

```bash
npm run build
```

ビルド成果物は `dist` ディレクトリに出力されます。

## 機能

### 取引データの管理
- **新規作成** - 取引日時、ロット数、購入レートを入力して新しい取引を作成
- **一覧表示** - すべての取引データを表形式で表示
- **編集** - 既存の取引データを編集
- **削除** - 不要な取引データを削除

### API連携
- バックエンドAPIとの通信はViteのプロキシ機能を使用
- `/api` へのリクエストは自動的に `http://localhost:8080` にプロキシされます

## プロジェクト構造

```
frontend/
├── src/
│   ├── api/
│   │   └── trades.ts         # APIクライアント
│   ├── components/
│   │   ├── TradeForm.tsx     # 取引フォームコンポーネント
│   │   └── TradeList.tsx     # 取引一覧コンポーネント
│   ├── types/
│   │   └── trade.ts          # 型定義
│   ├── App.tsx               # メインアプリケーション
│   ├── main.tsx              # エントリーポイント
│   └── index.css             # グローバルスタイル
├── index.html
├── package.json
├── tsconfig.json
└── vite.config.ts
```

## 使い方

1. **バックエンドサーバーを起動**
   ```bash
   # プロジェクトルートで実行
   go run cmd/server/main.go
   ```

2. **フロントエンドサーバーを起動**
   ```bash
   cd frontend
   npm run dev
   ```

3. **ブラウザでアクセス**
   `http://localhost:3000` を開く

## 注意事項

- フロントエンドを起動する前に、必ずバックエンドAPIサーバーが起動していることを確認してください
- APIサーバーは `http://localhost:8080` で起動している必要があります
