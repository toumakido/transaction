package main

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/toumakido/transaction/api/handler"
	"github.com/toumakido/transaction/api/repository"
)

func main() {
	// データベース接続を初期化
	db, err := repository.NewDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// リポジトリを初期化
	productRepo := repository.NewProductRepository(db)

	// ハンドラーを初期化
	productHandler := handler.NewProductHandler(productRepo)

	// Echoインスタンスを作成
	e := echo.New()

	// ミドルウェアを設定
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Use(handler.IDInjector()) // リクエストIDインジェクター

	// ルートを設定
	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"message":    "Transaction API Server",
			"request_id": handler.GetRequestID(c),
		})
	})

	// 商品関連のルート
	e.GET("/products/:id", productHandler.GetProduct)
	e.POST("/products/:id/process", productHandler.ProcessProduct)
	e.POST("/products/process", productHandler.ProcessProduct) // IDなしでも処理可能

	// サーバーを起動
	log.Println("Starting server on :8080")
	if err := e.Start(":8080"); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Failed to start server: %v", err)
	}
}
