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

type RepNewsFilter struct {
	// Id          int
	// Name        string
	// SurName     string
	// Patronymic  string
	// Age         int
	// Gender      string
	// Nationality string
	// Offset      int
	// Limit       int
}


