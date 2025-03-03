package main

import (
	"log"
	"net/http"
)

func main() {

	const port = "8080"
	mux := http.NewServeMux()

	// Serve static files from the current directory
	fileServer := http.FileServer(http.Dir("."))
	mux.Handle("/", fileServer)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(server.ListenAndServe())
}
