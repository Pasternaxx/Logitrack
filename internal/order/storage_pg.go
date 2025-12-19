package order

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type OrderStorage interface {
	GetAll() ([]Order, error)
	Save(ord Order) error
	GetByID(id int) (*Order, error)
	UpdateStatus(id int, status string) error
}

type PostgresOrderStorage struct {
	db *sqlx.DB
}

func NewPostgreOrderStorage(db *sqlx.DB) *PostgresOrderStorage {
	return &PostgresOrderStorage{
		db: db,
	}
}

func (o *PostgresOrderStorage) GetAll() ([]Order, error) {
	var orders []Order
	err := o.db.Select(&orders, "SELECT * FROM orders")
	if err != nil {
		return orders, fmt.Errorf("не удалось получить список заказов %v", err)
	}
	return orders, nil
}

func (o *PostgresOrderStorage) Save(ord Order) error {
	_, err := o.db.Exec("INSERT INTO orders (customer_name, status, created_at) VALUES ($1, $2, NOW())", ord.CustomerName, ord.Status)
	if err != nil {
		return fmt.Errorf("не удалось сохранить заказ %v", err)
	}
	return nil
}

func (o *PostgresOrderStorage) UpdateStatus(id int, status string) error {
	_, err := o.db.Exec("UPDATE orders SET status=$2 WHERE id=$1", id, status)
	if err != nil {
		return fmt.Errorf("не удалось обновить статус %v", err)
	}
	return nil
}

func (o *PostgresOrderStorage) GetByID(id int) (*Order, error) {
	var order Order
	err := o.db.Get(&order, "SELECT * FROM orders WHERE id=$1", id)
	if err != nil {
		return &order, fmt.Errorf("заказ не найден %v", err)
	}
	return &order, nil
}
