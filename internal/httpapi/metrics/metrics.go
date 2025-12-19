package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	//http_requests_total -Текущее значение
	//rate(http_requests_total[5m]) - кол-во запросов в сек. за последние 5 минут
	//increase(http_requests_total[1h]) - кол-во запросов за последний час
	TotalRequests = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Общее количество HTTP-запросов",
		})

	// http_request_duration_seconds_bucket{le="..."} — количество запросов ≤ заданного времени
	//http_request_duration_seconds_sum — суммарное время всех запросов
	//http_request_duration_seconds_count — общее число запросов
	RequestDuration = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Время обработки запроса",
			Buckets: prometheus.DefBuckets,
		})
	//http_success_total - все успешные запросы по путям
	SuccessCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_success_total",
			Help: "Успешные запросы",
		},
		[]string{"path"},
	)
)

func init() {
	prometheus.MustRegister(TotalRequests)
	prometheus.MustRegister(RequestDuration)
	prometheus.MustRegister(SuccessCounter)
}
