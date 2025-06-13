package model

import (
	"time"
)

// Product は商品を表す構造体です
type Product struct {
	ID        string    `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
	Stock     int       `db:"stock" json:"stock"`
	Price     float64   `db:"price" json:"price"`
	Version   int       `db:"version" json:"version"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

// ProductRequest は商品の更新リクエストを表す構造体です
type ProductRequest struct {
	Name  string  `json:"name"`
	Stock int     `json:"stock"`
	Price float64 `json:"price"`
}
