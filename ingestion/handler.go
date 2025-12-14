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

	logger.Printf("Incoming request: method=%s path=%s remote=%s",
		r.Method, r.URL.Path, r.RemoteAddr)

	if r.Method != http.MethodPost {
		logger.Println("Rejected request: invalid method")
		http.Error(w, "only POST", http.StatusMethodNotAllowed)
		return
	}

	b, err := io.ReadAll(r.Body)
	if err != nil {
		logger.Println("Failed to read body:", err)
		http.Error(w, "bad body", http.StatusBadRequest)
		return
	}

	var req UploadReq
	if err := json.Unmarshal(b, &req); err != nil {
		logger.Println("Invalid JSON body:", err)
		http.Error(w, "bad body", http.StatusBadRequest)
		return
	}

	logger.Printf("Upload received: doc_id=%s content_size=%d",
		req.DocId, len(req.Content))

	job := Job{
		ID:      req.DocId,
		Payload: req.Content,
		Created: time.Now(),
	}

	if err := wp.Enqueue(job); err != nil {
		logger.Printf("Queue full, enqueue failed for doc_id=%s: %v", job.ID, err)
		http.Error(w, "queue full", http.StatusTooManyRequests)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]string{
		"job_id": job.ID,
		"status": "accepted",
	})
}
