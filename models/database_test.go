package models

import (
	"database/sql"
	"errors"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

func TestDatabase(t *testing.T) {
	config, err := ParseConfig("../config_test.toml")

	if err != nil {
		t.Error(err)
	}

	db, err := sql.Open(config.Database.Driver, config.Database.GetUrl())

	if err != nil {
		t.Error(errors.Join(errors.New("Bad connection URL"), err))
	}

	err = db.Ping()

	if err != nil {
		t.Error(errors.Join(errors.New("error connecting to the database"), err))
	}
}
