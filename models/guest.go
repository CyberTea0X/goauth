package models

import (
	"database/sql"
	"errors"
	"fmt"
)

type Guest struct {
	Id       int64
	FullName string
}

// Calls rows.Next() and Scans the row into the guest struct
func (g *Guest) FromRow(rows *sql.Rows) (*Guest, error) {
	exists := rows.Next()
	if !exists {
		return g, errors.New(fmt.Sprintf("Can't scan user from row"))
	}

	err := rows.Scan(&g.Id, &g.FullName)
	if err != nil {
		return g, errors.New(fmt.Sprintf("Can't scan guest from row: %s", err.Error()))
	}

	return g, nil
}

// Inserts row that identifies token into the database (not token string)
func (g *Guest) InsertToDb(db *sql.DB) (int64, error) {
	const query = "INSERT INTO guests (full_name) VALUES (?)"
	res, err := db.Exec(query, &g.FullName)
	if err != nil {
		return 0, err
	}
	lastInsertId, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return lastInsertId, err
}
