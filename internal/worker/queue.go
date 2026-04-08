package worker

import "sync"

type Job struct {
	Type    string            `json:"type"`
	Payload map[string]string `json:"payload"`
}

type Queue interface {
	Enqueue(job Job)
	List() []Job
}

type InMemoryQueue struct {
	mu   sync.RWMutex
	jobs []Job
}

func NewInMemoryQueue() *InMemoryQueue {
	return &InMemoryQueue{}
}

func (q *InMemoryQueue) Enqueue(job Job) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.jobs = append(q.jobs, job)
}

func (q *InMemoryQueue) List() []Job {
	q.mu.RLock()
	defer q.mu.RUnlock()
	return append([]Job(nil), q.jobs...)
}
