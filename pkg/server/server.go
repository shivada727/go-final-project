package server

import (
	"fmt"
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
	fs := http.FileServer(http.Dir(webDir))

	mux := http.NewServeMux()

	mux.Handle("/", fs)

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
