package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
	"time"
)

const configPath = "config/local.yaml"

type Config struct {
	Env        string `yaml:"env"`
	HTTPServer `yaml:"http_server"`
	DbServer   `yaml:"db"`
}

type HTTPServer struct {
	// Host     string `yaml:"host" env-default: "localhost"`
	Port        string        `yaml:"port"             env-default: "8080"`
	Timeout     time.Duration `yaml:"timeout"          env-default: "4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout"     env-default: "60s"`
}

// Думаю, что к любой базе понадобятся одни и теже параметры подключения, поэтому решил определить их здесь
type DbServer struct {
	Host     string `yaml:"host"     env-default: "db"`
	Port     string `yaml:"port"     env-default: "5432"`
	Username string `yaml:"username" env-default: "postgres"`
	Password string `yaml:"password" env-default: "postgres"`
	DBName   string `yaml:"dbname" env-default: "postgres"`
	SSLMode  string `yaml:"sslmode"  env-default: "disable"`
}

func MustLoad() *Config {

	// Проверим, что файл конфига существует
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("Файл конфигурации не существует: %s", configPath)
	}

	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("Не удалось прочитать конфигурационный файл: %s", err)
	}

	// Когда база в контейнере, то имя контейнера будет db. Когда мы обращаемся к ней с локальной машины, то пробрасываем
	// ее интерфейс на localhost, поэтому подменим путь
	if cfg.Env == "local" {
		cfg.DbServer.Host = "localhost"
		cfg.DbServer.Port = "5437"
	}

	return &cfg
}
