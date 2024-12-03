package server

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Run(app *App, port int) {
	// Канал для остановки приложения
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Канал ошибок для gRPC сервера
	errChan := make(chan error, 1)

	// Запуск gRPC сервера
	go func() {
		log.Printf("Starting gRPC server on port %d...", port)
		errChan <- app.GRPCServer.Start(port)
	}()

	select {
	case <-stop:
		log.Println("Shutdown signal received...")
		shutdown(app)
	case err := <-errChan:
		log.Fatalf("gRPC server error: %v", err)
	}
}

func shutdown(app *App) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Останавливаем gRPC сервер
	if err := app.GRPCServer.Stop(); err != nil {
		log.Printf("Error stopping gRPC server: %v", err)
	}

	// Закрываем соединение с базой данных
	if err := app.DB.Close(); err != nil {
		log.Printf("Error closing DB: %v", err)
	} else {
		log.Println("Database connection closed.")
	}

	log.Println("Application shut down gracefully.")
}
