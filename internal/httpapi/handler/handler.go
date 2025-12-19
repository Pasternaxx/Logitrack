package handler

import (
	"awesomeProject/internal/httpapi/middleware"
	"awesomeProject/internal/order"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Handler struct {
	service *order.Service
}

func NewHandler(s *order.Service) *Handler {
	return &Handler{service: s}
}

func (h *Handler) RegisterRouters(mux *http.ServeMux) {
	mux.Handle("GET /orders", http.HandlerFunc(h.HandleOrders))
	mux.Handle("POST /orders", middleware.AuthMiddleware(http.HandlerFunc(h.HandleOrders)))
	mux.Handle("GET /orders/", http.HandlerFunc(h.orderByID))
	mux.Handle("PUT /orders/", middleware.AuthMiddleware(http.HandlerFunc(h.orderByID)))
	mux.Handle("/metrics", promhttp.Handler())
	mux.Handle("GET /health", http.HandlerFunc(h.healthHandler))
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}

func writeJSONError(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{
		"error": message,
	})
}

func (h *Handler) GetAllOrders(w http.ResponseWriter, r *http.Request) {
	orders, err := h.service.GetAll()
	if err != nil {
		log.Printf("Ошибка при получении заказов: %v", err)
		writeJSONError(w, "Метод не поддерживается", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	writeJSON(w, orders)
}

func (h *Handler) HandleCreateOrder(w http.ResponseWriter, r *http.Request) {
	var o order.Order
	if err := json.NewDecoder(r.Body).Decode(&o); err != nil {
		writeJSONError(w, "Невалидный JSON", http.StatusBadRequest)
		return
	}
	fmt.Printf("Decoded order: %+v\n", o)
	err := h.service.Save(o)
	if err != nil {
		writeJSONError(w, "Имя или статус не могут быть пустыми", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) HandleOrders(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetAllOrders(w, r)
	case http.MethodPost:
		h.HandleCreateOrder(w, r)
	default:
		writeJSONError(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) orderByID(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/orders/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("Ошибка при получении заказа: %v", err)
		writeJSONError(w, "Некорректный ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		ord, err := h.service.Get(id)
		if err != nil {
			writeJSONError(w, "Заказ не найден", http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusOK)
		writeJSON(w, ord)
	case http.MethodPut:
		var o order.Order
		if err := json.NewDecoder(r.Body).Decode(&o); err != nil {
			writeJSONError(w, "Невалидный JSON", http.StatusBadRequest)
			return
		}
		ord, err := h.service.Get(id)
		if err != nil {
			writeJSONError(w, "Заказ не найден", http.StatusNotFound)
			return
		}
		if err := h.service.UpdateStatus(id, o.Status); err != nil {
			writeJSONError(w, "Не удалось обновить статус заказа", http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
		writeJSON(w, ord)
	default:
		writeJSONError(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
