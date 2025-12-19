package awesomeProject

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

func InitConfig() {
	viper.SetDefault("port", "8080")
	if err := godotenv.Load("variable_environment.env"); err != nil {
		log.Println("Не удалось загрузить .env файл, используем переменные окружения системы")
	}
	viper.AutomaticEnv()
}
