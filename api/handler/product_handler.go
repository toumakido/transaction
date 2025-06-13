package handler

import (
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/samber/do"
	"github.com/toumakido/transaction/api/model"
	"github.com/toumakido/transaction/api/repository"
)

// ProductHandler は商品ハンドラーのインターフェースです
type ProductHandler interface {
	GetProduct(c echo.Context) error
	ProcessProduct(c echo.Context) error
}

// productHandler は商品ハンドラーの実装です
type productHandler struct {
	productRepo repository.ProductRepository
}

// NewProductHandler は新しい商品ハンドラーを作成します
func NewProductHandler(injector *do.Injector) ProductHandler {
	return &productHandler{
		productRepo: do.MustInvoke[repository.ProductRepository](injector),
	}
}

// GetProduct は指定されたIDの商品を取得します
func (h *productHandler) GetProduct(c echo.Context) error {
	// リクエストからIDを取得
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "ID is required"})
	}

	// 商品を取得
	var product *model.Product
	product, err := h.productRepo.FindByID(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	if product == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Product not found"})
	}

	return c.JSON(http.StatusOK, product)
}

// ProcessProduct は商品の在庫を処理します（トランザクション処理のサンプル）
func (h *productHandler) ProcessProduct(c echo.Context) error {
	// リクエストからIDを取得
	id := c.Param("id")
	if id == "" {
		// IDが指定されていない場合は新しいUUIDを生成
		id = uuid.New().String()
	}

	// リクエストから在庫変更量を取得
	stockChangeStr := c.QueryParam("stock_change")
	stockChange, err := strconv.Atoi(stockChangeStr)
	if err != nil || stockChangeStr == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid stock_change parameter"})
	}

	// トランザクション処理を実行
	if err := h.productRepo.ProcessProduct(id, stockChange); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// 更新後の商品を取得
	var product *model.Product
	product, err = h.productRepo.FindByID(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, product)
}
