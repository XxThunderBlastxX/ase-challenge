package product

import (
	"errors"

	"gorm.io/gorm"
)

type Service interface {
	CreateProduct(*Product) error
	GetAllProducts() ([]Product, error)
	GetProductByID(string) (*Product, error)
	UpdateProduct(string, *Product) error
	DeleteProduct(string) error

	IncermentStock(id string, quantity int) error
	DecrementStock(id string, quantity int) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}

// CreateProduct implements Service.
func (s *service) CreateProduct(product *Product) error {
	return s.repo.Create(product)
}

// DeleteProduct implements Service.
func (s *service) DeleteProduct(id string) error {
	return s.repo.Delete(id)
}

// GetAllProducts implements Service.
func (s *service) GetAllProducts() ([]Product, error) {
	return s.repo.GetAll()
}

// GetProductByID implements Service.
func (s *service) GetProductByID(id string) (*Product, error) {
	return s.repo.GetByID(id)
}

// UpdateProduct implements Service.
func (s *service) UpdateProduct(id string, product *Product) error {
	// Check if product exists
	_, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("product not found")
		}

		return err
	}

	return s.repo.Update(id, product)
}

// IncermentStock implements Service.
func (s *service) IncermentStock(id string, quantity int) error {
	p, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	p.StockQuantiy += quantity

	return s.repo.Update(id, p)
}

// DecrementStock implements Service.
func (s *service) DecrementStock(id string, quantity int) error {
	p, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	if p.StockQuantiy < quantity {
		return errors.New("insufficient stock to decrement")
	}

	p.StockQuantiy -= quantity

	return s.repo.Update(id, p)
}
