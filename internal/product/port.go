package product

type Repository interface {
	Create(*Product) error
	GetAll() ([]Product, error)
	GetByID(int64) (*Product, error)
}

type Service interface {
	Create(req CreateProductRequest) (*Product, error)
	GetAll() ([]Product, error)
	GetByID(int64) (*Product, error)
}
