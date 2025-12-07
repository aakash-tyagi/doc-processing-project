package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"
)

type UploadReq struct {
	DocId   string `json:"doc_id"`
	Content string `json:"content"`
}

func handleUpload(w http.ResponseWriter, r *http.Request, wp *WorkerPool, logger *log.Logger) {
	if r.Method != http.MethodPost {
		http.Error(w, "only POST", http.StatusMethodNotAllowed)
		return
	}

	b, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "bad body", http.StatusBadRequest)
		return
	}

	var req UploadReq
	if err := json.Unmarshal(b, &req); err != nil {
		http.Error(w, "bad body", http.StatusBadRequest)
		return
	}

	job := Job{
		ID:      req.DocId,
		Payload: req.Content,
		Created: time.Now(),
	}

	if err := wp.Enqueue(job); err != nil {
		logger.Println("enqueue failed:", err)
		http.Error(w, "queue full", http.StatusTooManyRequests)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]string{"job_id": job.ID, "status": "accepted"})
}
