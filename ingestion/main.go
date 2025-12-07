package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	logger := log.New(os.Stdout, "ingest: ", log.LstdFlags|log.Lmicroseconds)

	cfg := Config{
		ProcessorURL: getenv("", ""),
		WorkerCount:  5,
		QueueSize:    100,
	}

	// start worker pool
	wp := NewWorkerPool(cfg.WorkerCount, cfg.QueueSize, cfg.ProcessorURL, logger)
	wp.Start()

	// start server
	mux := http.NewServeMux()
	mux.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		handleUpload(w, r, wp, logger)
	})

	srv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	go func() {
		logger.Println("starting server on :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("server error: %v", err)
		}
	}()

}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

type Config struct {
	ProcessorURL string
	WorkerCount  int
	QueueSize    int
}
