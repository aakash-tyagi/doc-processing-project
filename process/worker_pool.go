package main

import (
	"log"
	"time"
)

type Processor struct {
	WorkerID int
}

func NewProcessor(id int) *Processor {
	return &Processor{WorkerID: id}
}

func (p *Processor) Process(task Task) {
	log.Printf("[worker %d] Started task %s", p.WorkerID, task.ID)
	time.Sleep(2 * time.Second)

	log.Printf("[Worker %d] Finished task: %s", p.WorkerID, task.ID)
}

type WorkerPool struct {
	tasks      chan Task
	workerSize int
}

func NewWorkerPool(workerSize int) *WorkerPool {
	return &WorkerPool{
		tasks:      make(chan Task, 100),
		workerSize: workerSize,
	}
}

func (wp *WorkerPool) Submit(task Task) {
	wp.tasks <- task
}

func (wp *WorkerPool) Close() {
	close(wp.tasks)
}
