package middleware

import (
	"net/http"
	"time"

	"github.com/sirupsen/logrus"

	"awesomeProject/internal/httpapi/metrics"

	"github.com/prometheus/client_golang/prometheus"
)

// loggingResponseWriter — оборачиваем ResponseWriter, чтобы перехватить статус
type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

// WriteHeader — переопределяем, чтобы запомнить статус
func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		times := time.Now()
		timer := prometheus.NewTimer(metrics.RequestDuration)

		// Создаём обёртку над ResponseWriter
		lrw := &loggingResponseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK, // значение по умолчанию
		}
		next.ServeHTTP(lrw, r)
		timer.ObserveDuration()
		metrics.TotalRequests.Inc()
		if lrw.statusCode >= 200 && lrw.statusCode < 300 {
			metrics.SuccessCounter.WithLabelValues(r.URL.Path).Inc()
		}

		fields := logrus.Fields{
			"time":   times,
			"method": r.Method,
			"path":   r.URL.Path,
			"status": lrw.statusCode,
		}
		if lrw.statusCode >= 400 {
			logrus.WithFields(fields).Error("HTTP запрос с ошибкой")
		} else {
			logrus.WithFields(fields).Info("HTTP запрос выполнен")
		}
	})
}

func RecoverMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				logrus.WithFields(logrus.Fields{
					"error": "internal error",
				})
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if authHeader := r.Header.Get("Authorization"); authHeader != "Bearersecret123" {
			http.Error(w, "Ошибка авторизации", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
