package types

import (
	"time"
)

type UserStore interface {
	GetUserByEmail(email string) (*User, error)
	GetUserByID(ID int) (*User, error)
	CreateUser(User) error
}

type ProductStore interface {
	GetProducts() ([]Product, error)
	CreateProduct(Product) error
	GetProductsByIDs(productIDs []int) ([]Product, error)
	GetProductByName(name string) (*Product, error)
	UpdateProduct(Product) error
}

type OrderStore interface {
	CreateOrder(Order) (int, error)
	CreateOrderItem(OrderItem) error
}

type CartItem struct {
	ProductID int `json:"productId"`
	Quantity  int `json:"quantity"`
}

type CartCheckoutPayload struct {
	Items []CartItem `json:"items" validate:"required"`
}

type AddProductPayload struct {
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description"`
	Image       string  `json:"image"`
	Price       float64 `json:"price" validate:"required,min=0.01"`
	Quantity    int     `json:"quantity" validate:"required,min=1" `
}

type OrderItem struct {
	ID        int       `json:"id"`
	OrderID   int       `json:"orderId"`
	ProductId int       `json:"productId"`
	Quantity  int       `json:"status"`
	Price     float64   `json:"address"`
	CreatedAt time.Time `json:"createdAt"`
}

type Order struct {
	ID        int       `json:"id"`
	UserID    int       `json:"userId"`
	Total     float64   `json:"total"`
	Status    string    `json:"status"`
	Address   string    `json:"address"`
	CreatedAt time.Time `json:"createdAt"`
}

type Product struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Image       string    `json:"image"`
	Price       float64   `json:"price"`
	Quantity    int       `json:"quantity"`
	CreatedAt   time.Time `json:"createdAt"`
}

type User struct {
	ID        int       `json:"id"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	CreateAt  time.Time `json:"createAt"`
}

type RegisterUserPayload struct {
	FirstName string `json:"firstName" validate:"required"`
	LastName  string `json:"lastName" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=3,max=20"`
}

type LoginUserPayload struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}
