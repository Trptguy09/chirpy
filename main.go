package main

import (
	"net/http"
)

func main() {
	ServerMux := http.NewServeMux()

	srv := http.Server{
		Addr:    ":8080",
		Handler: ServerMux,
	}
	srv.ListenAndServe()

}
