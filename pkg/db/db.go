package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "modernc.org/sqlite"
)

var (
	DB *sql.DB
)

const schema = `
CREATE TABLE scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8) NOT NULL DEFAULT "",
    title VARCHAR(255) NOT NULL DEFAULT "",
    comment TEXT NOT NULL DEFAULT "",
    repeat VARCHAR(128) NOT NULL DEFAULT ""
);

CREATE INDEX idx_scheduler_date ON scheduler(date);
`

func Init(dbFile string) error {
	install := false

	if _, error := os.Stat(dbFile); error != nil {
		install = true
	}

	db, error := sql.Open("sqlite", dbFile)

	if error != nil {
		return fmt.Errorf("open sqlite: %w", error)
	}

	if error := db.Ping(); error != nil {
		_ = db.Close()

		return fmt.Errorf("ping sqlite: %w", error)
	}

	if _, error := db.Exec(`PRAGMA foreign_keys = ON;`); error != nil {
		_ = db.Close()

		return fmt.Errorf("set pragma foreign_keys: %w", error)
	}
	if _, error := db.Exec(`PRAGMA busy_timeout = 5000;`); error != nil {
		_ = db.Close()

		return fmt.Errorf("set pragma busy_timeout: %w", error)
	}

	if install {
		transaction, error := db.Begin()

		if error != nil {
			_ = db.Close()

			return fmt.Errorf("begin tx: %w", error)
		}

		if _, error := transaction.Exec(schema); error != nil {
			_ = transaction.Rollback()

			_ = db.Close()

			return fmt.Errorf("apply schema: %w", error)
		}

		if err := transaction.Commit(); err != nil {
			_ = db.Close()

			return fmt.Errorf("commit schema: %w", error)
		}
	}

	DB = db

	return nil
}

func Close() error {
	if DB != nil {
		return DB.Close()
	}

	return nil
}
