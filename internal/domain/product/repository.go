package product

type Repository interface {
	Create(*Product) error
	GetAll() ([]Product, error)
	GetByID(string) (*Product, error)
	Update(string, *Product) error
	Delete(string) error
}
