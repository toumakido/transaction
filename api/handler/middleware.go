package handler

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// RequestIDKey はリクエストIDのキーです
const RequestIDKey = "request_id"

// IDInjector はリクエストIDを注入するミドルウェアです
func IDInjector() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// 新しいUUIDを生成
			requestID := uuid.New().String()

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
