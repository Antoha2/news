package main

import (
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/Antoha2/news/internal/config"
	"github.com/Antoha2/news/internal/repository"
	"github.com/Antoha2/news/internal/service"
	"github.com/Antoha2/news/internal/transport"
	"github.com/Antoha2/news/pkg/logger"
	"github.com/Antoha2/news/pkg/logger/sl"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/stdlib"
	"gopkg.in/reform.v1"
	"gopkg.in/reform.v1/dialects/postgresql"
)

func main() {
	run()
}

func run() {

	cfg := config.MustLoad()
	slogger := logger.SetupLogger(cfg.Env)
	dbx := MustInitDb(cfg)

	rep := repository.NewRep(slogger, dbx)
	serv := service.NewServ(cfg, slogger, rep)
	authServ := service.NewServAuth(cfg, slogger, rep)
	trans := transport.NewApi(cfg, slogger, serv, authServ)

	// a := authRepository.NewRep()

	go trans.StartHTTP()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit
	trans.Stop()
}

func MustInitDb(cfg *config.Config) *reform.DB {

	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
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
	sqlDB, err := sql.Open("pgx", stdlib.RegisterConnConfig(connConfig))
	if err != nil {
		slog.Warn("failed to create connection db", sl.Err(err))
		os.Exit(1)
	}

	err = sqlDB.Ping()
	if err != nil {
		slog.Warn("error to ping connection pool", sl.Err(err))
		os.Exit(1)
	}

	logger := log.New(os.Stderr, "SQL: ", log.Flags())
	dbx := reform.NewDB(sqlDB, postgresql.Dialect, reform.NewPrintfLogger(logger.Printf))

	slog.Info(fmt.Sprintf("Подключение к базе данных на http://127.0.0.1:%v\n", cfg.DBConfig.Port))
	return dbx
}
