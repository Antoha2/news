package config

import (
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env            string     `yaml:"env" env-default:"local"`
	StoragePath    string     `yaml:"storage_path" env-required:"true"`
	GRPC           GRPCConfig `yaml:"grpc"`
	HTTP           HTTPConfig `yaml:"http"`
	DBConfig       DBConfig   `yaml:"postgres"`
	MigrationsPath string
	TokenTTL       time.Duration `yaml:"token_ttl" env-default:"1h"`
}

type GRPCConfig struct {
	Port    int           `yaml:"grpcport"`
	Timeout time.Duration `yaml:"grpctimeout"`
}

type DBConfig struct {
	User     string `yaml:"dbuser"`
	Password string `yaml:"dbpassword"`
	Host     string `yaml:"dbhost"`
	Port     int    `yaml:"dbport"`
	Dbname   string `yaml:"dbname"`
	Sslmode  string `yaml:"dbsslmode"`
}

type ContextKey string
type HTTPConfig struct {
	HostAddr   int        `yaml:"httpport"`
	ContextKey ContextKey `yaml:"httpcontextkey"`
}

// const HTTPAddr = ":8180"
const HashSalt = "antoha"
const HashMinLength = 3

func MustLoad() *Config {
	configPath := fetchConfigPath()
	if configPath == "" {
		panic("config path is empty")
	}

	// check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic("config file does not exist: " + configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic("config path is empty: " + err.Error())
	}

	return &cfg
}

func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}

// func GetConfig() *Config {
// 	return &Config{
// 		DB: DBConfig{
// 			Dbname:   "auth",
// 			User:     "user",
// 			Password: "user",
// 			Host:     "localhost", //"postgres", //
// 			Port:     5432,
// 			Sslmode:  "",
// 		},
// 		HTTP: HTTPConfig{
// 			HostAddr:   os.Getenv("HOST_ADDR"),
// 			ContextKey: "History",
// 		},
// 	}
// }

// type Config struct {
// 	Env            string     `yaml:"env" env-default:"local"`
// 	StoragePath    string     `yaml:"storage_path" env-required:"true"`
// 	GRPC           GRPCConfig `yaml:"grpc"`
// 	MigrationsPath string
// 	TokenTTL       time.Duration `yaml:"token_ttl" env-default:"48h"`
// 	DB             DBConfig      `yaml:"DB"`
// 	HTTP           HTTPConfig
// }

// type GRPCConfig struct {
// 	Port    int           `yaml:"port"`
// 	Timeout time.Duration `yaml:"timeout"`
// }

// // type Config struct {
// // 	HTTP HTTPConfig
// // 	DB   DBConfig
// // }

// type ContextKey string

// type HTTPConfig struct {
// 	HostAddr   string
// 	ContextKey ContextKey
// }

// func MustLoad() *Config {
// 	configPath := fetchConfigPath()
// 	if configPath == "" {
// 		panic("config path is empty")
// 	}

// 	// check if file exists
// 	if _, err := os.Stat(configPath); os.IsNotExist(err) {
// 		panic("config file does not exist: " + configPath)
// 	}

// 	var cfg Config

// 	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
// 		panic("config path is empty: " + err.Error())
// 	}

// 	cfg.HTTP.HostAddr = os.Getenv("HOST_ADDR")
// 	cfg.HTTP.ContextKey = "history"

// 	return &cfg

// }

// func fetchConfigPath() string {
// 	var res string

// 	flag.StringVar(&res, "config", "", "path to config file")
// 	flag.Parse()

// 	if res == "" {
// 		res = os.Getenv("CONFIG_PATH")
// 	}

// 	return res
// }
