package product

import (
	"github.com/xxthunderblastxx/ase-challenge/internal/pkg/model"
)

type Product struct {
	model.BaseModel
	Name             string `json:"name" gorm:"not null"`
	Description      string `json:"description"`
	StockQuantity    int    `json:"stock_quantity" gorm:"not null"`
	LowStockThresold int    `json:"low_stock_threshold" gorm:"not null"`
}
