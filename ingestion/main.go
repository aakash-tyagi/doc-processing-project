package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	jobsReceived = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "jobs_received_total",
			Help: "Total number of jobs received",
		},
	)

	jobsProcessed = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "jobs_processed_total",
			Help: "Total number of jobs processed",
		},
	)

	jobProcessingDuration = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "job_processing_duration_seconds",
			Help:    "Time taken to process jobs",
			Buckets: prometheus.DefBuckets,
		},
	)
)

func main() {
	logger := log.New(os.Stdout, "ingest: ", log.LstdFlags|log.Lmicroseconds)

	prometheus.MustRegister(
		jobsReceived,
		jobsProcessed,
		jobProcessingDuration,
	)

	cfg := Config{
		ProcessorURL: getenv("PROCESSOR_URL", "url to process"),
		WorkerCount:  5,
		QueueSize:    100,
	}

	// start worker pool
	wp := NewWorkerPool(cfg.WorkerCount, cfg.QueueSize, cfg.ProcessorURL, logger)
	wp.Start()

	// start server
	mux := http.NewServeMux()

	mux.Handle("/metrics", promhttp.Handler())

	mux.HandleFunc("/ingest", func(w http.ResponseWriter, r *http.Request) {
		handleUpload(w, r, wp, logger)
	})

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
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

	// graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
	logger.Println("shutting down")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	srv.Shutdown(ctx)
	wp.Stop()

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
