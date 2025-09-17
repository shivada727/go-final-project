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

	if _, err := os.Stat(dbFile); err != nil {
		install = true
	}

	db, err := sql.Open("sqlite", dbFile)

	if err != nil {
		return fmt.Errorf("open sqlite: %w", err)
	}

	if err := db.Ping(); err != nil {
		_ = db.Close()

		return fmt.Errorf("ping sqlite: %w", err)
	}

	if _, err := db.Exec(`PRAGMA foreign_keys = ON;`); err != nil {
		_ = db.Close()

		return fmt.Errorf("set pragma foreign_keys: %w", err)
	}
	if _, err := db.Exec(`PRAGMA busy_timeout = 5000;`); err != nil {
		_ = db.Close()

		return fmt.Errorf("set pragma busy_timeout: %w", err)
	}

	if install {
		transaction, err := db.Begin()

		if err != nil {
			_ = db.Close()
			return fmt.Errorf("begin tx: %w", err)
		}

		if _, err := transaction.Exec(schema); err != nil {
			_ = transaction.Rollback()

			_ = db.Close()

			return fmt.Errorf("apply schema: %w", err)
		}

		if err := transaction.Commit(); err != nil {
			_ = db.Close()

			return fmt.Errorf("commit schema: %w", err)
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
