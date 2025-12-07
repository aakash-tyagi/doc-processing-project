package main

import (
	"encoding/json"
	"net/http"
)

type Task struct {
	ID      string `json:"id"`
	Payload string `json:"payload"`
}

func HandleProcess(pool *WorkerPool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var task Task

		if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
			http.Error(w, "invalid payload", http.StatusBadRequest)
			return
		}

		pool.Submit(task)

		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte("Task queued for processing"))
	}
}
