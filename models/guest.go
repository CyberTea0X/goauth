package models

import (
	"database/sql"
	"errors"
	"fmt"
)

const guests_table_ddl = "" +
	"CREATE TABLE IF NOT EXISTS `guests` (" +
	"`id` int(10) unsigned NOT NULL AUTO_INCREMENT," +
	"`full_name` varchar(100) NOT NULL," +
	"PRIMARY KEY (`id`)" +
	") ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;"

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

// Creates refresh token table if it doesn't already exist
func CreateGuestTable(db *sql.DB) error {
	_, err := db.Exec(guests_table_ddl)
	return err
}
