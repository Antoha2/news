package repository

import (
	"log/slog"

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
	Id         int    //`reform:"id,pk"`
	Title      string //`reform:"citle"`
	Content    string //`reform:"content"`
	Categories []int
}

type SearchTerms struct {
	Offset int
	Limit  int
}
