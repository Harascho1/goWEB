package order

import (
	"database/sql"

	"github.com/Harascho1/goWEB/types"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db: db,
	}
}

func (s *Store) CreateOrder(order types.Order) (int, error) {
	res, err := s.db.Exec(
		"INSERT INTO orders (userId, total, status, address) VALUES (?,?,?,?)",
		order.UserID,
		order.Total,
		order.Status,
		order.Address,
	)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}
func (s *Store) CreateOrderItem(orderItem types.OrderItem) error {
	_, err := s.db.Exec(
		"INSERT INTO orders (orderId, productId, quantity, price) VALUES (?,?,?,?)",
		orderItem.OrderID,
		orderItem.ProductId,
		orderItem.Quantity,
		orderItem.Price,
	)
	if err != nil {
		return err
	}

	return nil
}
