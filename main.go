package main

import (
	"log"
	"path/filepath"

	"go-final-project/pkg/server"
)

func main() {

	webDir := filepath.Join(".", "web")

	port := server.ResolvePort()

	if error := server.Run(webDir, port); error != nil {
		log.Fatal(error)
	}
}
