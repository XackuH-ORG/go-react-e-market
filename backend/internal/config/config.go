package config

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Env  string     `yaml:"env" env:"ENV" env-default:"local"`
	HTTP HTTPConfig `yaml:"http"`
	DB   DBConfig   `yaml:"db"`
}

type HTTPConfig struct {
	Addr         string        `yaml:"addr" env:"HTTP_ADDR" env-default:"0.0.0.0:8080"`
	ReadTimeout  time.Duration `yaml:"read_timeout" env:"HTTP_READ_TIMEOUT" env-default:"4s"`
	WriteTimeout time.Duration `yaml:"write_timeout" env:"HTTP_WRITE_TIMEOUT" env-default:"4s"`
	IdleTimeout  time.Duration `yaml:"idle_timeout" env:"HTTP_IDLE_TIMEOUT" env-default:"60s"`
	ReqTimeout   time.Duration `yaml:"req_timeout" env:"HTTP_REQ_TIMEOUT" env-default:"60s"`
	User         string        `yaml:"user" env-required:"true"`
	Password     string        `yaml:"password" env-required:"true" env:"HTTP_SERVER_PASSWORD"`
}

type DBConfig struct {
	DSN    string `yaml:"dsn" env:"GOOSE_DBSTRING" env-required:"true"`
	Driver string `yaml:"driver" env:"GOOSE_DRIVER" env-default:"postgres"`
}

// MustLoad загружает конфигурацию. Паникует, если не удается загрузить конфигурацию
func MustLoad() *Config {
	// 1. Попытка загрузить .env
	// Сначала ищем .env в текущей директории
	if err := godotenv.Load(); err != nil {
		// Если не нашли в текущей рабочей директории, пробуем найти рядом с бинарником
		if execPath, err := os.Executable(); err == nil {
			execDir := filepath.Dir(execPath)
			_ = godotenv.Load(filepath.Join(execDir, ".env"))
		}
	}

	// 2. Определяем путь к файлу конфигурации
	var configPath string

	// Используем изолированный FlagSet, чтобы избежать паники/конфликтов с другими флагами
	fs := flag.NewFlagSet("config", flag.ContinueOnError)
	fs.StringVar(&configPath, "config", "", "path to config file")
	// Игнорируем ошибки парсинга флагов здесь, так как другие флаги могут быть обработаны в main
	_ = fs.Parse(os.Args[1:])

	// Если флаг не задан, проверяем переменную окружения CONFIG_PATH
	if configPath == "" {
		configPath = os.Getenv("CONFIG_PATH")
	}

	// Если путь все еще не задан, ищем файлы конфигурации по умолчанию в текущей папке и рядом с бинарником
	if configPath == "" {
		var execDir string
		if execPath, err := os.Executable(); err == nil {
			execDir = filepath.Dir(execPath)
		}

		defaults := []string{
			"config/local.yaml",
			"local.yaml",
		}
		if execDir != "" {
			defaults = append(defaults,
				filepath.Join(execDir, "config/local.yaml"),
				filepath.Join(execDir, "local.yaml"),
			)
		}

		for _, path := range defaults {
			if _, err := os.Stat(path); err == nil {
				configPath = path
				break
			}
		}
	}

	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set and no default config file (local.yaml) was found")
	}

	// Проверяем существование файла конфигурации перед чтением
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}

	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("failed to read config from %s: %s", configPath, err)
	}

	return &cfg
}
