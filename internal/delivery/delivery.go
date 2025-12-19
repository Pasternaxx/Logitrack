package delivery

import (
	"awesomeProject/internal/order"
	"fmt"
	"sync"
	"time"
)

// функция: "горутина доставки заказа"
func StartDeliveryWorker(ch chan order.Order, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-time.After(2 * time.Second):
			fmt.Println("Временное ожидание доставки...")
		case ord := <-ch:
			fmt.Println("Order", ord.ID, "delivered")
			return
		}
	}
}
