package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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

	pool := NewWorkerPool(5)
	pool.Start()

	logger := log.New(os.Stdout, "[PROCESS] ", log.LstdFlags|log.Lmicroseconds)

	mux := http.NewServeMux()

	mux.Handle("/metrics", promhttp.Handler())

	mux.HandleFunc("/process", HandleProcess(pool, logger))
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	srv := &http.Server{
		Addr:    ":8081",
		Handler: mux,
	}

	// Graceful shutdown
	go func() {
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
		<-stop

		log.Println("Shutting down Processor Service...")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		pool.Close()
		if err := srv.Shutdown(ctx); err != nil {
			log.Printf("Server Shutdown Error: %v", err)
		}
	}()

	log.Println("Processor Service running on port :8081")

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server error: %v", err)
	}

}
