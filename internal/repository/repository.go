package repository

import (
	"log/slog"

	"github.com/pkg/errors"
	"gopkg.in/reform.v1"
)

type RepImpl struct {
	log *slog.Logger
	DB  *reform.DB
}

func NewRep(log *slog.Logger, dbx *reform.DB) *RepImpl {
	return &RepImpl{
		log: log,
		DB:  dbx,
	}
}

type RepNews struct {
	Id         int    `reform:"id,pk"`
	Title      string `reform:"title"`
	Content    string `reform:"content"`
	Categories []int
}

type SearchTerms struct {
	Offset int
	Limit  int
}

var (
	ErrUserExists   = errors.New("user already exists")
	ErrUserNotFound = errors.New("user not found")
	ErrAppNotFound  = errors.New("app not found")
)
