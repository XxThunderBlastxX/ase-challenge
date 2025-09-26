package product

type Repository interface {
	Create(*Product) error
	GetAll() ([]Product, error)
	GetByID(string) (*Product, error)
	UpdateAllColumn(string, *Product) error
	UpdateSingleColumn(string, string, any) error
	Delete(string) error
}
