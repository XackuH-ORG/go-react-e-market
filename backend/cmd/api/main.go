package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	// 1. Загружаем .env файл.
	// Сначала ищем в текущей директории (Production: /opt/market/.env)
	// Если нет, ищем на уровень выше (Development: запуск из backend/cmd/api/)
	err := godotenv.Load(".env")
	if err != nil {
		err = godotenv.Load("../.env")
	}

	// Логируем результат, чтобы в journalctl было четко видно, подцепился ли файл
	if err != nil {
		log.Println("Warning: No .env file found, relying on system environment variables or defaults")
	} else {
		log.Println("Success: .env file loaded")
	}

	// 2. Читаем переменную HTTP_ADDR
	httpAddr := os.Getenv("HTTP_ADDR")
	if httpAddr == "" {
		httpAddr = ":8080" // Фолбэк, если переменная пустая
		log.Println("Warning: HTTP_ADDR not set, falling back to", httpAddr)
	}

	// 3. Настраиваем роуты
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"status": "ok", "message": "CI/CD loop closed. Backend is running."}`)
	})

	// 4. Запускаем сервер
	log.Printf("Server is starting on %s\n", httpAddr)
	if err := http.ListenAndServe(httpAddr, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
