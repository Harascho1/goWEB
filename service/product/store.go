package product

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/Harascho1/goWEB/types"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) GetProducts() ([]types.Product, error) {
	rows, err := s.db.Query("SELECT * FROM product")
	if err != nil {
		return nil, err
	}

	products := make([]types.Product, 0)
	for rows.Next() {
		p, err := scanRowsIntoProduct(rows)
		if err != nil {
			return nil, err
		}

		products = append(products, *p)
	}

	return products, nil
}

func (s *Store) GetProductsByIDs(productIDs []int) ([]types.Product, error) {
	placeholders := strings.Repeat(",?", len(productIDs)-1)
	query := fmt.Sprintf("SELECT * FROM product WHERE id IN (?%s)", placeholders)

	args := make([]interface{}, len(productIDs))
	for i, v := range productIDs {
		args[i] = v
	}

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}

	products := []types.Product{}
	for rows.Next() {
		p, err := scanRowsIntoProduct(rows)
		if err != nil {
			return nil, err
		}

		products = append(products, *p)
	}

	return products, nil
}

func (s *Store) UpdateProduct(product types.Product) error {
	_, err := s.db.Exec(
		"UPDATE product SET name = ?, price = ?, image = ?, description = ?, quantity = ? WHERE id = ?",
		product.Name,
		product.Price,
		product.Image,
		product.Description,
		product.Quantity,
		product.ID,
	)
	return err

}

func (s *Store) GetProductByName(name string) (*types.Product, error) {

	rows, err := s.db.Query(
		"SELECT * FROM product WHERE name = ?",
		name,
	)
	if err != nil {
		return nil, err
	}

	p := new(types.Product)
	for rows.Next() {
		p, err = scanRowsIntoProduct(rows)
		if err != nil {
			return nil, err
		}
	}

	if p.ID == 0 {
		return nil, fmt.Errorf("product not found")
	}

	return p, nil

}

func scanRowsIntoProduct(row *sql.Rows) (*types.Product, error) {
	p := new(types.Product)

	err := row.Scan(
		&p.ID,
		&p.Name,
		&p.Description,
		&p.Image,
		&p.Price,
		&p.Quantity,
		&p.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (s *Store) CreateProduct(product types.Product) error {
	_, err := s.db.Exec(
		"INSERT INTO product (name, description, image, price, quantity) VALUES (?,?,?,?,?)",
		product.Name,
		product.Description,
		product.Image,
		product.Price,
		product.Quantity,
	)
	if err != nil {
		return err
	}
	return nil
}
