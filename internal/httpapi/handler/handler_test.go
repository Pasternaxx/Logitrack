package handler

import (
	"awesomeProject/internal/order"
	"bytes"
	"encoding/json"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

var h = NHandler()

func NHandler() *Handler {
	testDb, err := sqlx.Connect("postgres", "host=localhost port=5433 user=postgres password=123 dbname=test_db sslmode=disable")
	if err != nil {
		log.Fatalln("Не удалось подключиться к БД:", err)
	}

	store := order.NewPostgreOrderStorage(testDb)

	service := order.NewService(store)
	h := NewHandler(service)
	return h
}

func TestHandleOrders(t *testing.T) {
	req := httptest.NewRequest("GET", "/orders", nil)
	w := httptest.NewRecorder()

	h.HandleOrders(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Ожидали 200, получили %d", resp.StatusCode)
	}
}

func TestHandleCreateOrder(t *testing.T) {
	sliceOrders := []order.Order{
		{CustomerName: "Josh", Status: "created"},
		{CustomerName: "", Status: "created"},
		{CustomerName: "Bella", Status: ""}}
	for i, ord := range sliceOrders {
		body, err := json.Marshal(ord)
		if err != nil {
			t.Fatalf("не удалось сериализовать Order: %v", err)
		}
		req := httptest.NewRequest("POST", "/orders", bytes.NewReader(body))
		w := httptest.NewRecorder()

		h.HandleCreateOrder(w, req)
		resp := w.Result()
		if i == 0 && resp.StatusCode != http.StatusCreated {
			t.Errorf("Ожидали 201, получили %d", resp.StatusCode)
		}
		if i == 1 && resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Ожидали 400, пустое имя или статус, получили %d", resp.StatusCode)
		}
		if i == 2 && resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Ожидали 400, пустое имя или статус, получили %d", resp.StatusCode)
		}
	}
}

func TestOrderByID(t *testing.T) {
	req := []*http.Request{
		httptest.NewRequest("GET", "/orders/1", nil),
		httptest.NewRequest("GET", "/orders/1098", nil),
		httptest.NewRequest("GET", "/orders/", nil),
	}
	for i, request := range req {
		w := httptest.NewRecorder()
		h.orderByID(w, request)
		resp := w.Result()
		if i == 0 && resp.StatusCode != http.StatusOK {
			t.Errorf("Ожидали 200, получили %d", resp.StatusCode)
		}
		if i == 1 && resp.StatusCode != http.StatusNotFound {
			t.Errorf("Ожидали 404, получили %d", resp.StatusCode)
		}
		if i == 2 && resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Ожидали 400, получили %d", resp.StatusCode)
		}
	}
}

func TestPutOrderByID(t *testing.T) {
	sliceStatus := []order.Order{
		{Status: "cancelled"},
		{Status: "shippedsa"},
		{Status: ""},
		{CustomerName: "Jjodnfd", Status: "shipped"},
	}
	for i, stat := range sliceStatus {
		body, err := json.Marshal(stat)
		if err != nil {
			t.Fatalf("не удалось сериализовать Order: %v", err)
		}
		req := httptest.NewRequest("PUT", "/orders/1", bytes.NewReader(body))
		w := httptest.NewRecorder()

		h.orderByID(w, req)
		resp := w.Result()

		if i == 0 && resp.StatusCode != http.StatusOK {
			t.Errorf("Ожидали 200, получили %d", resp.StatusCode)
		}
		if i == 1 && resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Ожидали 400, получили %d", resp.StatusCode)
		}
		if i == 2 && resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Ожидали 400, получили %d", resp.StatusCode)
		}
		if i == 3 && resp.StatusCode != http.StatusOK {
			t.Errorf("Ожидали 200, получили %d", resp.StatusCode)
		}
	}
}
