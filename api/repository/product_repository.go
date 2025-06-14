package repository

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/samber/do"
	"github.com/toumakido/transaction/api/model"
)

// ProductRepository は商品リポジトリのインターフェースです
type ProductRepository interface {
	FindByID(id string) (*model.Product, error)
	FindByIDWithLock(tx *sql.Tx, id string) (*model.Product, error)
	Update(tx *sql.Tx, product *model.Product) error
	Create(tx *sql.Tx, product *model.Product) error
	ProcessProduct(id string, stock int) (string, error)
}

// productRepository は商品リポジトリの実装です
type productRepository struct {
	db *DB
}

// NewProductRepository は新しい商品リポジトリを作成します
func NewProductRepository(injector *do.Injector) ProductRepository {
	return &productRepository{db: do.MustInvoke[*DB](injector)}
}

// FindByID は指定されたIDの商品を取得します
func (r *productRepository) FindByID(id string) (*model.Product, error) {
	var product model.Product
	query := "SELECT * FROM products WHERE id = ?"
	err := r.db.Get(&product, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find product: %w", err)
	}
	return &product, nil
}

// FindByIDWithLock はトランザクション内で指定されたIDの商品を取得し、行ロックを取得します
func (r *productRepository) FindByIDWithLock(tx *sql.Tx, id string) (*model.Product, error) {
	var product model.Product
	query := "SELECT * FROM products WHERE id = ? FOR UPDATE"
	err := tx.QueryRow(query, id).Scan(
		&product.ID,
		&product.Name,
		&product.Stock,
		&product.Price,
		&product.Version,
		&product.CreatedAt,
		&product.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find product with lock: %w", err)
	}
	return &product, nil
}

// Update は商品を更新します
func (r *productRepository) Update(tx *sql.Tx, product *model.Product) error {
	query := `
		UPDATE products 
		SET name = ?, stock = ?, price = ?, version = version + 1 
		WHERE id = ? AND version = ?
	`
	result, err := tx.Exec(query, product.Name, product.Stock, product.Price, product.ID, product.Version)
	if err != nil {
		return fmt.Errorf("failed to update product: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("optimistic lock failed: product was updated by another transaction")
	}

	// バージョンを更新
	product.Version++
	return nil
}

// Create は新しい商品を作成します
func (r *productRepository) Create(tx *sql.Tx, product *model.Product) error {
	if product.ID == "" {
		product.ID = uuid.New().String()
	}

	query := `
		INSERT INTO products (id, name, stock, price, version) 
		VALUES (?, ?, ?, ?, 1)
	`
	_, err := tx.Exec(query, product.ID, product.Name, product.Stock, product.Price)
	if err != nil {
		return fmt.Errorf("failed to create product: %w", err)
	}

	product.Version = 1
	return nil
}

// ProcessProduct は商品の在庫を処理します（トランザクション処理のサンプル）
func (r *productRepository) ProcessProduct(id string, stockChange int) (string, error) {
	// トランザクションを開始
	tx, err := r.db.Begin()
	if err != nil {
		return "", fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			// エラーが発生した場合はロールバック
			if rbErr := tx.Rollback(); rbErr != nil {
				log.Printf("failed to rollback: %v", rbErr)
			}
		}
	}()

	// 商品をロックして取得
	product, err := r.FindByIDWithLock(tx, id)
	if err != nil {
		return "", err
	}

	var msg string
	// 商品が存在しない場合は新規作成
	if product == nil {
		newProduct := &model.Product{
			ID:    id,
			Name:  "New Product",
			Stock: stockChange,
			Price: 1000.0,
		}
		if err = r.Create(tx, newProduct); err != nil {
			return "", err
		}
		msg = fmt.Sprintf("Created new product with ID: %s, Stock: %d", id, stockChange)
	} else {
		// 商品が存在する場合は在庫を更新
		product.Stock += stockChange
		if err = r.Update(tx, product); err != nil {
			return "", err
		}
		msg = fmt.Sprintf("Updated product with ID: %s, New Stock: %d", id, product.Stock)
	}

	// トランザクションをコミット
	if err = tx.Commit(); err != nil {
		return "", fmt.Errorf("failed to commit transaction: %w", err)
	}

	return msg, nil
}
