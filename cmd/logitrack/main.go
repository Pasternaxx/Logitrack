package main

import (
	"awesomeProject"
	http2 "awesomeProject/internal/httpapi/handler"
	"awesomeProject/internal/httpapi/middleware"
	"awesomeProject/internal/order"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

func main() {
	awesomeProject.InitConfig()

	dbHost := viper.GetString("DB_HOST")
	dbPort := viper.GetString("DB_PORT")
	dbUser := viper.GetString("DB_USER")
	dbPassword := viper.GetString("DB_PASSWORD")
	dbName := viper.GetString("DB_NAME")
	dbSSLMode := viper.GetString("DB_SSLMODE")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", dbHost, dbPort, dbUser, dbPassword, dbName, dbSSLMode)

	time.Sleep(time.Second * 3)

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Println(dbHost, dbPort, dbUser, dbPassword, dbName, dbSSLMode)
		log.Fatalln("Не удалось подключиться к БД:", err)
	}
	log.Println("Успешное подключение к БД mydb")
	defer db.Close()

	port := viper.GetString("APP_PORT")
	if port == "" {
		log.Fatalln("Порт пустой")
	}

	if viper.GetString("APP_ENV") == "prod" {
		log.Println("Продуктион")
	} else {
		log.Println("Разработка")
	}

	store := order.NewPostgreOrderStorage(db)
	service := order.NewService(store)
	h := http2.NewHandler(service)

	mux := http.NewServeMux()
	h.RegisterRouters(mux)
	log.Printf("Сервер запущен на http://localhost:%s", port)
	wrapped := middleware.LoggingMiddleware(middleware.RecoverMiddleware(mux))
	err = http.ListenAndServe(":8080", wrapped)
	if err != nil {
		log.Println("Ошибка сервера")
	}

}
