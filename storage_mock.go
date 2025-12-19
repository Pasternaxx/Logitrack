package awesomeProject

import (
	"awesomeProject/internal/order"
	"errors"
	"fmt"
)

type OrderStoragexc interface {
	GetAll() []order.Order
	Save(ord order.Order) (order.Order, error)
	GetByID(id int) (*order.Order, error)
	Update(order order.Order) error
}

type OrderStorageMock struct {
	Orders []order.Order
}

func (o *OrderStorageMock) Save(ord order.Order) (order.Order, error) {
	if ord.ID == 0 {
		return ord, errors.New("ID не может быть 0")
	}
	o.Orders = append(o.Orders, ord)
	defer fmt.Println("запрос завершён")
	return ord, nil
}
func (o *OrderStorageMock) GetAll() []order.Order {
	return o.Orders
}

func (o *OrderStorageMock) GetByID(id int) (*order.Order, error) {
	for i := range o.Orders {
		if o.Orders[i].ID == id {
			return &o.Orders[i], nil
		}
	}
	return nil, fmt.Errorf("order %d: %w", id, errors.New("заказ не найден"))
}

func (o *OrderStorageMock) Update(order order.Order) error {
	for i := range o.Orders {
		if o.Orders[i].ID == order.ID {
			o.Orders[i] = order
			return nil
		}
	}
	return errors.New("заказ не найден")
}
