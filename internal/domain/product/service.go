package product

import (
	"errors"

	apperrors "github.com/xxthunderblastxx/ase-challenge/internal/pkg/errors"
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
	// Validate required fields
	if product.Name == "" {
		return apperrors.NewMissingRequiredDataError("name")
	}

	if product.StockQuantity < 0 {
		return apperrors.NewInvalidInputError("stock quantity cannot be negative")
	}

	err := s.repo.Create(product)
	if err != nil {
		// Handle duplicate entry errors (if name should be unique)
		// This depends on your database constraints
		return apperrors.NewDatabaseError("failed to create product: " + err.Error())
	}

	return nil
}

// DeleteProduct implements Service.
func (s *service) DeleteProduct(id string) error {
	if id == "" {
		return apperrors.NewMissingRequiredDataError("id")
	}

	// Check if product exists first
	_, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperrors.NewProductNotFoundError(id)
		}
		return apperrors.NewDatabaseError("failed to check product existence: " + err.Error())
	}

	err = s.repo.Delete(id)
	if err != nil {
		return apperrors.NewDatabaseError("failed to delete product: " + err.Error())
	}

	return nil
}

// GetAllProducts implements Service.
func (s *service) GetAllProducts() ([]Product, error) {
	products, err := s.repo.GetAll()
	if err != nil {
		return nil, apperrors.NewDatabaseError("failed to retrieve products: " + err.Error())
	}

	return products, nil
}

// GetProductByID implements Service.
func (s *service) GetProductByID(id string) (*Product, error) {
	if id == "" {
		return nil, apperrors.NewMissingRequiredDataError("id")
	}

	product, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NewProductNotFoundError(id)
		}
		return nil, apperrors.NewDatabaseError("failed to retrieve product: " + err.Error())
	}

	return product, nil
}

// UpdateProduct implements Service.
func (s *service) UpdateProduct(id string, product *Product) error {
	if id == "" {
		return apperrors.NewMissingRequiredDataError("id")
	}

	// Validate required fields
	if product.Name == "" {
		return apperrors.NewMissingRequiredDataError("name")
	}

	if product.StockQuantity < 0 {
		return apperrors.NewInvalidInputError("stock quantity cannot be negative")
	}

	// Check if product exists
	_, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperrors.NewProductNotFoundError(id)
		}
		return apperrors.NewDatabaseError("failed to check product existence: " + err.Error())
	}

	err = s.repo.Update(id, product)
	if err != nil {
		return apperrors.NewDatabaseError("failed to update product: " + err.Error())
	}

	return nil
}

// IncrementStock implements Service.
func (s *service) IncermentStock(id string, quantity int) error {
	if id == "" {
		return apperrors.NewMissingRequiredDataError("id")
	}

	if quantity <= 0 {
		return apperrors.NewInvalidInputError("increment quantity must be greater than 0")
	}

	p, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperrors.NewProductNotFoundError(id)
		}
		return apperrors.NewDatabaseError("failed to retrieve product: " + err.Error())
	}

	p.StockQuantity += quantity

	err = s.repo.Update(id, p)
	if err != nil {
		return apperrors.NewDatabaseError("failed to update product stock: " + err.Error())
	}

	return nil
}

// DecrementStock implements Service.
func (s *service) DecrementStock(id string, quantity int) error {
	if id == "" {
		return apperrors.NewMissingRequiredDataError("id")
	}

	if quantity <= 0 {
		return apperrors.NewInvalidInputError("decrement quantity must be greater than 0")
	}

	p, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperrors.NewProductNotFoundError(id)
		}
		return apperrors.NewDatabaseError("failed to retrieve product: " + err.Error())
	}

	if p.StockQuantity < quantity {
		return apperrors.NewInsufficientStockError(p.StockQuantity, quantity)
	}

	p.StockQuantity -= quantity

	err = s.repo.Update(id, p)
	if err != nil {
		return apperrors.NewDatabaseError("failed to update product stock: " + err.Error())
	}

	return nil
}
