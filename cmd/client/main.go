package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"
)

// Product は商品を表す構造体です
type Product struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Stock     int       `json:"stock"`
	Price     float64   `json:"price"`
	Version   int       `json:"version"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func main() {
	baseURL := "http://localhost:8080"

	// 既存の商品を取得
	fmt.Println("=== 既存の商品を取得 ===")
	product, err := getProduct(baseURL, "00000000-0000-0000-0000-000000000001")
	if err != nil {
		log.Fatalf("Failed to get product: %v", err)
	}
	printProduct(product)

	// 並行して同じ商品の在庫を更新（競合状態のテスト）
	fmt.Println("\n=== 並行して同じ商品の在庫を更新（競合状態のテスト） ===")
	var wg sync.WaitGroup
	for i := 1; i <= 5; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			stockChange := -10 // 在庫を10減らす
			updatedProduct, err := processProduct(baseURL, "00000000-0000-0000-0000-000000000001", stockChange)
			if err != nil {
				log.Printf("Client %d: Failed to process product: %v", i, err)
				return
			}
			log.Printf("Client %d: Updated product stock to %d (version: %d)", i, updatedProduct.Stock, updatedProduct.Version)
		}(i)
		// 少し待機して競合を発生させる
		time.Sleep(100 * time.Millisecond)
	}
	wg.Wait()

	// 更新後の商品を取得
	fmt.Println("\n=== 更新後の商品を取得 ===")
	product, err = getProduct(baseURL, "00000000-0000-0000-0000-000000000001")
	if err != nil {
		log.Fatalf("Failed to get product: %v", err)
	}
	printProduct(product)

	// 存在しない商品を処理（新規作成のテスト）
	fmt.Println("\n=== 存在しない商品を処理（新規作成のテスト） ===")
	newProductID := "test-product-id"
	newProduct, err := processProduct(baseURL, newProductID, 50)
	if err != nil {
		log.Fatalf("Failed to process new product: %v", err)
	}
	printProduct(newProduct)

	// 新しく作成した商品を取得
	fmt.Println("\n=== 新しく作成した商品を取得 ===")
	product, err = getProduct(baseURL, newProductID)
	if err != nil {
		log.Fatalf("Failed to get new product: %v", err)
	}
	printProduct(product)

	// IDなしで商品を処理（自動ID生成のテスト）
	fmt.Println("\n=== IDなしで商品を処理（自動ID生成のテスト） ===")
	autoProduct, err := processProductWithoutID(baseURL, 30)
	if err != nil {
		log.Fatalf("Failed to process auto product: %v", err)
	}
	printProduct(autoProduct)
}

// getProduct は指定されたIDの商品を取得します
func getProduct(baseURL, id string) (*Product, error) {
	resp, err := http.Get(fmt.Sprintf("%s/products/%s", baseURL, id))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get product: %s, status: %d", string(body), resp.StatusCode)
	}

	var product Product
	if err := json.NewDecoder(resp.Body).Decode(&product); err != nil {
		return nil, err
	}

	return &product, nil
}

// processProduct は商品の在庫を処理します
func processProduct(baseURL, id string, stockChange int) (*Product, error) {
	url := fmt.Sprintf("%s/products/%s/process?stock_change=%d", baseURL, id, stockChange)
	resp, err := http.Post(url, "application/json", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to process product: %s, status: %d", string(body), resp.StatusCode)
	}

	var product Product
	if err := json.NewDecoder(resp.Body).Decode(&product); err != nil {
		return nil, err
	}

	return &product, nil
}

// processProductWithoutID はIDなしで商品を処理します
func processProductWithoutID(baseURL string, stockChange int) (*Product, error) {
	url := fmt.Sprintf("%s/products/process?stock_change=%d", baseURL, stockChange)
	resp, err := http.Post(url, "application/json", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to process product: %s, status: %d", string(body), resp.StatusCode)
	}

	var product Product
	if err := json.NewDecoder(resp.Body).Decode(&product); err != nil {
		return nil, err
	}

	return &product, nil
}

// printProduct は商品情報を表示します
func printProduct(p *Product) {
	fmt.Printf("ID: %s\n", p.ID)
	fmt.Printf("Name: %s\n", p.Name)
	fmt.Printf("Stock: %d\n", p.Stock)
	fmt.Printf("Price: %.2f\n", p.Price)
	fmt.Printf("Version: %d\n", p.Version)
	fmt.Printf("Created At: %s\n", p.CreatedAt.Format(time.RFC3339))
	fmt.Printf("Updated At: %s\n", p.UpdatedAt.Format(time.RFC3339))
}
