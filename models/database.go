package models

import (
	"database/sql"
	"errors"
	_ "github.com/go-sql-driver/mysql"
)

const refreshTableDDL = "" +
	"CREATE TABLE IF NOT EXISTS `refresh_tokens` (" +
	"`id` int(10) unsigned NOT NULL AUTO_INCREMENT," +
	"`device_id` int(10) unsigned NOT NULL," +
	"`user_id` int(10) unsigned NOT NULL," +
	"`expires_at` int(10) unsigned NOT NULL," +
	"PRIMARY KEY (`id`)," +
	"UNIQUE KEY `refresh_tokens_device_id_IDX` (`device_id`,`user_id`,`expires_at`) USING BTREE" +
	") ENGINE=InnoDB AUTO_INCREMENT=26 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;"

func SetupDatabase(config *DatabaseConfig) (*sql.DB, error) {
	db, err := sql.Open(config.Driver, config.GetUrl())

	if err != nil {
		return nil, errors.Join(errors.New("error connecting to the database"), err)
	}

	if _, err = db.Exec(refreshTableDDL); err != nil {
		return nil, errors.Join(errors.New("Failed to create refresh token table: "), err)
	}

	return db, nil
}

func TruncateDatabase(db *sql.DB) error {
	if _, err := db.Exec("TRUNCATE TABLE refresh_tokens"); err != nil {
		return errors.Join(errors.New("error clearing refresh token table"), err)
	}
	return nil
}
