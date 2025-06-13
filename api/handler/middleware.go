package handler

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/samber/do"
)

// RequestIDKey はリクエストIDのキーです
const RequestIDKey = "request_id"

// RequestIDGenerator はリクエストIDを生成するインターフェースです
type RequestIDGenerator interface {
	Generate() string
}

// UUIDRequestIDGenerator はUUIDを使用してリクエストIDを生成する実装です
type UUIDRequestIDGenerator struct{}

// Generate は新しいリクエストIDを生成します
func (g *UUIDRequestIDGenerator) Generate() string {
	return uuid.New().String()
}

// IDInjector はリクエストIDを注入するミドルウェアです
func IDInjector(injector *do.Injector) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// DIコンテナからRequestIDGeneratorを取得
			generator := do.MustInvoke[RequestIDGenerator](injector)

			// リクエストIDを生成
			requestID := generator.Generate()

			// コンテキストにリクエストIDを設定
			c.Set(RequestIDKey, requestID)

			// リクエストヘッダーにリクエストIDを設定
			c.Response().Header().Set("X-Request-ID", requestID)

			// 次のハンドラーを呼び出す
			return next(c)
		}
	}
}

// GetRequestID はコンテキストからリクエストIDを取得します
func GetRequestID(c echo.Context) string {
	id := c.Get(RequestIDKey)
	if id == nil {
		return ""
	}
	return id.(string)
}
