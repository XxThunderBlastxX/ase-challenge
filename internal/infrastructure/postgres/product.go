package postgres

import (
	"github.com/xxthunderblastxx/ase-challenge/internal/domain/product"
)

type productRepository struct {
	conn *ConnectionManager
}

func NewProductRepository(conn *ConnectionManager) product.Repository {
	return &productRepository{
		conn: conn,
	}
}

// Create implements product.Repository.
func (r *productRepository) Create(product *product.Product) error {
	if err := r.conn.DB.Create(product).Error; err != nil {
		return err
	}

	return nil
}

// Delete implements product.Repository.
func (r *productRepository) Delete(id string) error {
	if err := r.conn.DB.Delete(&product.Product{}, "id = ?", id).Error; err != nil {
		return err
	}

	return nil
}

// GetAll implements product.Repository.
func (r *productRepository) GetAll() ([]product.Product, error) {
	var products []product.Product

	if err := r.conn.DB.Find(&products).Error; err != nil {
		return nil, err
	}

	return products, nil
}

// GetByID implements product.Repository.
func (r *productRepository) GetByID(id string) (*product.Product, error) {
	var p product.Product

	if err := r.conn.DB.First(&p, "id = ?", id).Error; err != nil {
		return nil, err
	}

	return &p, nil
}

// UpdateAllColumn implements product.Repository.
func (r *productRepository) UpdateAllColumn(id string, product *product.Product) error {
	if err := r.conn.DB.Updates(product).Error; err != nil {
		return err
	}

	return nil
}

func (r *productRepository) UpdateSingleColumn(id string, column string, value any) error {
	if err := r.conn.DB.
		Model(&product.Product{}).
		Where("id = ?", id).
		Update(column, value).
		Error; err != nil {
		return err
	}

	return nil
}
