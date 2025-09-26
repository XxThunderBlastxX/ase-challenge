package product

type Service interface {
	CreateProduct(*Product) error
	GetAllProducts() ([]Product, error)
	GetProductByID(string) (*Product, error)
	UpdateProduct(string, *Product) error
	DeleteProduct(string) error
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
	return s.repo.Update(id, product)
}
