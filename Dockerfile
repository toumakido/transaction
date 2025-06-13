FROM golang:1.21-alpine

WORKDIR /app

# 依存関係のインストール
COPY go.mod go.sum ./
RUN go mod download

# アプリケーションのビルド
COPY . .
RUN go build -o main ./cmd/api

# 実行
EXPOSE 8080
CMD ["./main"]
