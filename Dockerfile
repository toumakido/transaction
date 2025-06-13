# ビルドステージ
FROM golang:1.24.2-alpine AS builder

WORKDIR /app

# 依存関係のインストール
COPY go.mod go.sum ./
RUN go mod download

# アプリケーションのビルド
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/api

# 実行ステージ
FROM gcr.io/distroless/static:nonroot

WORKDIR /

# ビルドステージからバイナリをコピー
COPY --from=builder /app/main /main

# 非rootユーザーとして実行
USER nonroot:nonroot

# 実行
EXPOSE 8080
ENTRYPOINT ["/main"]
