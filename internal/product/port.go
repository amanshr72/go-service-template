package product

type Repository interface {
	Create(*Product) error
	GetAll() ([]Product, error)
	GetByID(int) (*Product, error)
}

type Service interface {
	Create(name string, price float64) (*Product, error)
	GetAll() ([]Product, error)
	GetByID(int) (*Product, error)
}
