package user

// Repository port — data layer contract, Three adapters implement this: postgres, mongodb, inmemory.
type Repository interface {
	Create(u *User) error
	FindAll() ([]User, error)
	FindByID(id int) (*User, error)
	FindByActive(active bool) ([]User, error)
	Count() (int, error)
	Update(u *User) error
	Delete(id int) error
}

// Service port — business logic contract.
type Service interface {
	Create(input CreateUserInput) (*User, error)
	GetAll() ([]User, error)
	GetByID(id int) (*User, error)
	GetByActive(active bool) ([]User, error)
	GetCount() (int, error)
	Update(id int, input UpdateUserInput) (*User, error)
	Delete(id int) error
}
