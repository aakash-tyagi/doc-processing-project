package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"sync"
	"time"
)

type Job struct {
	ID      string
	Payload string
	Created time.Time
}

type WorkerPool struct {
	workers  int
	queue    chan Job
	client   *http.Client
	procURL  string
	logger   *log.Logger
	wg       sync.WaitGroup
	stopCh   chan struct{}
	stopOnce sync.Once
}

func NewWorkerPool(workers, qsize int, processesorURL string, logger *log.Logger) *WorkerPool {
	return &WorkerPool{
		workers: workers,
		queue:   make(chan Job, qsize),
		client:  &http.Client{Timeout: 10 * time.Second},
		procURL: processesorURL,
		logger:  logger,
	}
}

func (w *WorkerPool) Start() {
	for i := 0; i < w.workers; i++ {
		w.wg.Add(1)
		go w.worker(i)
	}
}

func (w *WorkerPool) Stop() {
	w.stopOnce.Do(func() {
		close(w.stopCh)
		w.logger.Println("closing queue")

		time.AfterFunc(2*time.Second, func() {
			close(w.queue)
		})
		w.wg.Wait()
	})
}

func (w *WorkerPool) Enqueue(j Job) error {
	select {
	case w.queue <- j:
		return nil
	default:
		return errors.New("queue full")
	}
}

func (w *WorkerPool) worker(id int) {
	defer w.wg.Done()
	w.logger.Printf("worker-%d started", id)
	for j := range w.queue {
		// if stop signal, break
		select {
		case <-w.stopCh:
			w.logger.Printf("worker-%d stopping due to shutdown", id)
			return
		default:
		}

		// forward to processor
		body := map[string]string{"job_id": j.ID, "content": j.Payload}
		buf, _ := jsonMarshal(body)
		req, err := http.NewRequest(http.MethodPost, w.procURL, bytes.NewBuffer(buf))
		if err != nil {
			w.logger.Println("create req err:", err)
			continue
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := w.client.Do(req)
		if err != nil {
			w.logger.Println("forward err:", err)
			continue
		}

		resp.Body.Close()
		w.logger.Printf("worker-%d forwareded job%s stautus=%s", id, j.ID, resp.Status)
	}
	w.logger.Printf("worker-%d exiting", id)
}

func jsonMarshal(v any) ([]byte, error) {
	return json.Marshal(v)
}
