package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	// Пытаемся загрузить из ../.env first (development)
	err := godotenv.Load("../.env")
	if err != nil {
		// Если не найден, то пытаемся загрузить из текущей директории (production)
		err = godotenv.Load(".env")
	}
	if err != nil {
		fmt.Println("Warning: .env file not found, using default values")
	}

	// Получение HTTP_ADDR из переменных окружения, с дефолтным значением ":8080"
	httpAddr := os.Getenv("HTTP_ADDR")
	if httpAddr == "" {
		httpAddr = ":8080"
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `{"status": "ok", "message": "CI/CD loop closed. Backend is running."}`)
	})
	fmt.Printf("Server is running on %s\n", httpAddr)
	http.ListenAndServe(httpAddr, nil)
}
