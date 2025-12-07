package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {

	mux := http.NewServeMux()
	mux.HandleFunc("/process", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("processing request")
	})

	srv := &http.Server{
		Addr:    ":8082",
		Handler: mux,
	}

	log.Println("Processor Service running on port :8082")

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server error: %v", err)
	}

}
