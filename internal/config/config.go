package config

import (
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Env      string `env:"ENV" env-default:"local"`
	HTTP     HTTPConfig
	DBConfig DBConfig
	TokenTTL time.Duration `env:"TOKEN_TTL" env-default:"1h"`
}

type DBConfig struct {
	User                  string `env:"DB_USER" env-default:"user"`
	Password              string `env:"DB_PASSWORD" env-default:"user"`
	Host                  string `env:"DB_HOST" env-default:"localhost"`
	Port                  string `env:"DB_PORT" env-default:"5432"`
	Dbname                string `env:"DB_DBNAME" env-default:"newsdb"`
	Sslmode               string `env:"DB_SSLMODE" env-default:"disable"`
	DefaultPropertyLimit  int    `env:"DB_LIMIT" env-default:"100"`
	DefaultPropertyOffset int    `env:"DB_OFFSET" env-default:"0"`
}

type HTTPConfig struct {
	HostPort string `env:"HTTP_PORT" env-default:"3000"`
}

// загрузка конфига из .env
func MustLoad() *Config {

	if err := godotenv.Load(); err != nil {
		panic("No .env file found" + err.Error())
	}
	var cfg Config
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		panic("cannot read config: " + err.Error())
	}

	return &cfg
}
