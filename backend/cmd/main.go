package main

import (
	"fmt"
	"log/slog"
	"os"
)

func main() {
	cfg := config{
		addr: ":8080",
		db:   dbConfig{},
	}

	api := application{
		config: cfg,
	}

	if err := api.run(api.mount()); err != nil {
		slog.Error(fmt.Sprintf("Сервер не запустился, ошибка: %v", err))
		os.Exit(1)
	}
}
