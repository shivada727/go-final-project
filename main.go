package main

import (
	"go-final-project/pkg/db"
	"go-final-project/pkg/server"
	"log"
	"os"
	"path/filepath"
)

func main() {

	dbFile := os.Getenv("TODO_DBFILE")

	if dbFile == "" {
		dbFile = filepath.Join(".", "scheduler.db")
	}

	if error := db.Init(dbFile); error != nil {
		log.Fatal("db init:", error)
	}

	webDir := filepath.Join(".", "web")

	port := server.ResolvePort()

	if error := server.Run(webDir, port); error != nil {
		log.Fatal(error)
	}
}
