package order

import "errors"

type Service struct {
	storage OrderStorage
}

func NewService(o *PostgresOrderStorage) *Service {
	return &Service{
		storage: o,
	}
}

func (s *Service) Save(ord Order) error {
	// Валидация полей, кроме ID
	if err := validatorSave(ord); err != nil {
		return err
	}
	return s.storage.Save(ord)
}

func (s *Service) GetAll() ([]Order, error) {
	return s.storage.GetAll()
}

func (s *Service) Get(id int) (*Order, error) {
	return s.storage.GetByID(id)
}

func (s *Service) UpdateStatus(id int, newStatus string) error {
	if err := validatorUpdate(newStatus); err != nil {
		return err
	}
	return s.storage.UpdateStatus(id, newStatus)
}

func validatorSave(order Order) error {
	if order.CustomerName == "" || order.Status == "" {
		return errors.New("не указано имя или статус заказа")
	}
	return nil
}
func validatorUpdate(newStatus string) error {
	var orderStatuses = []string{"created", "shipped", "delivered", "cancelled"}
	for _, stat := range orderStatuses {
		if stat == newStatus {
			return nil
		}
	}
	return errors.New("неверный статус")
}
