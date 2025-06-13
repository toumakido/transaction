# Transaction API

トランザクション処理のサンプルAPIサーバーです。商品の在庫管理を例にした、トランザクション処理の実装例を提供します。

## 技術スタック

- **言語**: Go 1.24.2
- **フレームワーク**: Echo v4
- **データベース**: MySQL 8.0
- **依存性注入**: github.com/samber/do
- **コンテナ化**: Docker, Docker Compose
- **デプロイ**: Distroless イメージ

## 機能

- 商品情報の取得
- 商品在庫の更新（トランザクション処理）
- 新規商品の作成
- リクエストIDの自動生成と追跡

## セットアップ

### 前提条件

- Docker
- Docker Compose

### 起動方法

```bash
# リポジトリをクローン
git clone https://github.com/toumakido/transaction.git
cd transaction

# Docker Composeでアプリケーションを起動
docker compose up -d
```

アプリケーションは http://localhost:8080 で利用可能になります。

## API エンドポイント

### 商品情報の取得

```
GET /products/:id
```

**レスポンス例**:

```json
{
  "id": "00000000-0000-0000-0000-000000000001",
  "name": "Sample Product",
  "stock": 100,
  "price": 1000,
  "version": 1,
  "created_at": "2025-06-13T12:00:00Z",
  "updated_at": "2025-06-13T12:00:00Z"
}
```

### 商品在庫の更新

```
POST /products/:id/process?stock_change=10
```

**パラメータ**:
- `stock_change`: 在庫の変更量（正の値で増加、負の値で減少）

**レスポンス例**:

```json
{
  "id": "00000000-0000-0000-0000-000000000001",
  "name": "Sample Product",
  "stock": 110,
  "price": 1000,
  "version": 2,
  "created_at": "2025-06-13T12:00:00Z",
  "updated_at": "2025-06-13T12:05:00Z"
}
```

### 新規商品の作成（IDを指定）

```
POST /products/:id/process?stock_change=50
```

### 新規商品の作成（ID自動生成）

```
POST /products/process?stock_change=30
```

## プロジェクト構造

```
transaction/
├── api/
│   ├── handler/
│   │   ├── middleware.go     # ミドルウェア（リクエストIDインジェクターなど）
│   │   └── product_handler.go # 商品関連のハンドラー
│   ├── model/
│   │   └── product.go        # 商品モデル
│   └── repository/
│       ├── db.go             # データベース接続
│       └── product_repository.go # 商品リポジトリ
├── cmd/
│   ├── api/
│   │   └── main.go           # APIサーバーのエントリーポイント
│   └── client/
│       └── main.go           # テスト用クライアント
├── db/
│   └── migrations/
│       └── 01_init.sql       # 初期データベーススキーマ
├── docker-compose.yml        # Docker Compose設定
├── Dockerfile                # Dockerビルド設定
├── go.mod                    # Goモジュール定義
└── go.sum                    # Goモジュールのチェックサム
```

## 依存性注入（DI）

このプロジェクトでは、`github.com/samber/do`ライブラリを使用して依存性注入を実装しています。主な利点は以下の通りです：

- コンポーネント間の依存関係を明示的に管理
- テストやモック化が容易
- 型安全なDI実装
- グローバルインジェクターによる自動的な依存関係の解決

## コンテナ化

アプリケーションはマルチステージビルドを使用してコンテナ化されています：

1. **ビルドステージ**: `golang:1.24.2-alpine`イメージを使用してアプリケーションをビルド
2. **実行ステージ**: `gcr.io/distroless/static:nonroot`イメージを使用して実行

これにより、以下の利点があります：

- 軽量なコンテナイメージ
- セキュリティの向上（最小限の攻撃対象領域）
- 非rootユーザーでの実行

## クライアントの使用方法

テスト用クライアントを使用して、APIの動作を確認できます：

```bash
# クライアントを実行
go run cmd/client/main.go
```

クライアントは以下のテストを実行します：

1. 既存の商品を取得
2. 並行して同じ商品の在庫を更新（競合状態のテスト）
3. 更新後の商品を取得
4. 存在しない商品を処理（新規作成のテスト）
5. 新しく作成した商品を取得
6. IDなしで商品を処理（自動ID生成のテスト）
