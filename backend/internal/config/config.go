package config

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Env  string `env:"ENV" env-default:"dev"`
	HTTP HTTPConfig
	DB   DBConfig
	JWT  JWTConfig
}

type HTTPConfig struct {
	Addr         string        `env:"HTTP_ADDR" env-default:"0.0.0.0:8080"`
	ReadTimeout  time.Duration `env:"HTTP_READ_TIMEOUT" env-default:"4s"`
	WriteTimeout time.Duration `env:"HTTP_WRITE_TIMEOUT" env-default:"4s"`
	IdleTimeout  time.Duration `env:"HTTP_IDLE_TIMEOUT" env-default:"60s"`
	ReqTimeout   time.Duration `env:"HTTP_REQ_TIMEOUT" env-default:"60s"`
}

type DBConfig struct {
	DSN    string `env:"GOOSE_DBSTRING" env-required:"true"`
	Driver string `env:"GOOSE_DRIVER" env-default:"postgres"`
}

type JWTConfig struct {
	Secret string        `env:"JWT_SECRET" env-required:"true"`
	TTL    time.Duration `env:"JWT_TTL" env-default:"1h"`
}

// MustLoad загружает конфигурацию из .env файла и переменных окружения.
// Паникует, если отсутствуют обязательные переменные (отмеченные как env-required).
func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")

	if configPath != "" {
		// Если задан CONFIG_PATH, пытаемся загрузить указанный файл
		if err := godotenv.Load(configPath); err != nil {
			log.Printf("warning: failed to load config from CONFIG_PATH (%s): %v", configPath, err)
		}
	} else {
		// Если не задан, ищем .env файл по иерархии
		loaded := false

		// 1. Пытаемся загрузить .env из текущей рабочей директории
		if err := godotenv.Load(); err == nil {
			loaded = true
		}

		// 2. Если не нашли, идем вверх по дереву каталогов (ищем до корня проекта - где лежит go.mod)
		if !loaded {
			dir, err := os.Getwd()
			if err == nil {
				for {
					parent := filepath.Dir(dir)
					if parent == dir { // Достигли корня файловой системы
						break
					}
					dir = parent

					envPath := filepath.Join(dir, ".env")
					if _, err := os.Stat(envPath); err == nil {
						if err := godotenv.Load(envPath); err == nil {
							loaded = true
							break
						}
					}

					// Если нашли go.mod, значит это корень проекта, дальше вверх не идем
					if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
						break
					}
				}
			}
		}

		// 3. Последняя попытка - рядом с бинарником (полезно для продакшна)
		if !loaded {
			if execPath, err := os.Executable(); err == nil {
				execDir := filepath.Dir(execPath)
				_ = godotenv.Load(filepath.Join(execDir, ".env"))
			}
		}
	}

	// Читаем переменные окружения в структуру
	// Если .env файл не был найден вообще, readEnv попытается прочитать системные переменные окружения напрямую (например, в Docker/K8s)
	var cfg Config
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		log.Fatalf("failed to read config: %s", err)
	}

	return &cfg
}
