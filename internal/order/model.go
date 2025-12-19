package order

import "time"

type Order struct {
	ID           int       `db:"id" json:"id" redis:"id"`
	CustomerName string    `db:"customer_name" json:"customer_name" redis:"customer_name"`
	Status       string    `db:"status" json:"status" redis:"status"`
	CreatedAt    time.Time `db:"created_at" json:"created_at" redis:"created_at"`
}

func NewOrder(id int, name string, status string) *Order {
	return &Order{
		ID:           id,
		CustomerName: name,
		Status:       status,
	}
}

func (o Order) IsDelivered() bool {
	return o.Status == "delivered"
}
