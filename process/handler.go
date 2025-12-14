package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type Task struct {
	ID      string `json:"id"`
	Payload string `json:"payload"`
}

func HandleProcess(pool *WorkerPool, logger *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		logger.Printf("Incoming process request: method=%s path=%s remote=%s",
			r.Method, r.URL.Path, r.RemoteAddr)

		if r.Method != http.MethodPost {
			logger.Println("Rejected request: invalid method")
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var task Task
		if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
			logger.Println("Invalid payload:", err)
			http.Error(w, "invalid payload", http.StatusBadRequest)
			return
		}

		logger.Printf("Task received: id=%s payload_size=%d",
			task.ID, len(task.Payload))

		pool.Submit(task)

		logger.Printf("Task queued successfully: id=%s latency=%s",
			task.ID, time.Since(start))

		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte("Task queued for processing"))
	}
}
