package main

import (
	"log"
	"os"

	"grpc-file-service/cmd/server"
)

func main() {
	dbConnStr := os.Getenv("DB_CONN_STR") // Получаем строку подключения к БД из переменной окружения
	if dbConnStr == "" {
		log.Fatal("DB_CONN_STR environment variable is required")
	}

	// Инициализация приложения
	app, err := server.InitApp(dbConnStr, "storage")
	if err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}

	// Запуск приложения на порту 1337
	server.Run(app, 1337)
}
