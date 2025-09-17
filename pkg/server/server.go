package server

import (
	"fmt"
	"go-final-project/pkg/api"
	"net/http"
	"os"
)

const defaultPort = "7540"

func ResolvePort() string {
	port := os.Getenv("TODO_PORT")

	if port == "" {
		port = defaultPort
	}

	return port
}

func NewHandler(webDir string) (http.Handler, error) {
	if _, err := os.Stat(webDir); err != nil {
		return nil, fmt.Errorf("web dir not found: %s (%w)", webDir, err)
	}

	fs := http.FileServer(http.Dir(webDir))
	mux := http.NewServeMux()
	mux.Handle("/", fs)

	api.Init(mux)

	return mux, nil
}

func Run(webDir string, port string) error {
	handler, err := NewHandler(webDir)

	if err != nil {
		return err
	}

	addr := ":" + port

	fmt.Printf("Server is listening on http://localhost%s/ (web dir: %s)\n", addr, webDir)

	return http.ListenAndServe(addr, handler)
}
