package data

import (
	"database/sql"
	"github.com/google/wire"
	_ "github.com/lib/pq"
)

var ProviderSet = wire.NewSet(NewData, NewTaskRepo)

type Data struct {
	*sql.DB
}

func NewData(sqlConn string) (*Data, error) {
	db, err := sql.Open("postgres", sqlConn)
	return &Data{db}, err
}
