package main

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/Antoha2/news/internal/auth/internal/config"
	"github.com/Antoha2/news/internal/auth/internal/repository"
	"github.com/Antoha2/news/internal/auth/internal/services"
	transport "github.com/Antoha2/news/internal/auth/internal/transport/http"
	logger "github.com/Antoha2/news/internal/auth/pkg/logger"
	"github.com/Antoha2/news/internal/auth/pkg/logger/sl"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
)

func main() {
	Run()
}

func Run() {

	cfg := config.MustLoad()
	log := logger.SetupLogger(logger.EnvLocal)

	dbx, err := MustInitDb(cfg)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	authRep := repository.NewRepAuth(log, dbx, cfg.TokenTTL)
	authServ := services.NewServAuth(log, authRep)
	authTrans := transport.NewApi(authServ, log, cfg.HTTP.HostAddr)

	go authTrans.StartHTTP()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit
	authTrans.Stop()

}

func MustInitDb(cfg *config.Config) (*sqlx.DB, error) {

	connString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.DBConfig.User,
		cfg.DBConfig.Password,
		cfg.DBConfig.Host,
		cfg.DBConfig.Port,
		cfg.DBConfig.Dbname,
		cfg.DBConfig.Sslmode,
	)

	connConfig, err := pgx.ParseConfig(connString)
	if err != nil {
		slog.Warn("failed to parse config", sl.Err(err))
		os.Exit(1)
	}

	// Make connections
	dbx, err := sqlx.Open("pgx", stdlib.RegisterConnConfig(connConfig))
	if err != nil {
		slog.Warn("failed to create connection db", sl.Err(err))
		os.Exit(1)
	}

	err = dbx.Ping()
	if err != nil {
		slog.Warn("error to ping connection pool", sl.Err(err))
		os.Exit(1)
	}
	slog.Info(fmt.Sprintf("Подключение к базе данных на http://127.0.0.1:%v\n", cfg.DBConfig.Port))
	return dbx, nil
}
