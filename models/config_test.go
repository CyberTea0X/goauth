package models

import "testing"

func TestParseConfig(t *testing.T) {
	_, err := ParseConfig("../config_test.toml")

	if err != nil {
		t.Error(err)
	}
}
