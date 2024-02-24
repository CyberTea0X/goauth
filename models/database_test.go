package models

import "testing"

func DatabaseTest(t *testing.T) {
	config, err := ParseConfig("../config_test.toml")

	if err != nil {
		t.Error(err)
	}

	db, err := SetupDatabase(&config.Database)

	if err != nil {
		t.Error(err)
	}

	err = db.Ping()

	if err != nil {
		t.Error(err)
	}
}
