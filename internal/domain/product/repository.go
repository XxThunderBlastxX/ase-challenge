package product

type Repository interface {
	Create(*Product) (Product, error)
	GetAll() ([]Product, error)
	GetByID(string) (Product, error)
	Update(string, *Product) (Product, error)
	Delete(string) error
}
