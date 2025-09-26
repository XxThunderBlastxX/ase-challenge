package product

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	ID           string `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Name         string `json:"name" gorm:"not null"`
	Description  string `json:"description"`
	StockQuantiy int    `json:"stock_quantity" gorm:"not null"`
}
