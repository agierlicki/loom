package jobs

import (
	"fmt"
	"sync"
)

type InMemoryStore struct {
	mu   sync.RWMutex
	jobs map[string]Job
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		jobs: make(map[string]Job),
	}
}

func (s *InMemoryStore) Create(job Job) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.jobs[job.ID]; exists {
		return fmt.Errorf("job with id %q already exists", job.ID)
	}

	s.jobs[job.ID] = job
	return nil
}

func (s *InMemoryStore) Get(id string) (Job, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	job, exists := s.jobs[id]

	if !exists {
		return Job{}, fmt.Errorf("job with id %q not found!", id)
	}

	return job, nil
}

func (s *InMemoryStore) Update(job Job) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.jobs[job.ID]; !exists {
		return fmt.Errorf("job with id %q not found", job.ID)
	}

	s.jobs[job.ID] = job
	return nil
}

func (s *InMemoryStore) List() ([]Job, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	jobs := make([]Job, 0, len(s.jobs))
	for _, job := range s.jobs {
		jobs = append(jobs, job)
	}

	return jobs, nil
}
