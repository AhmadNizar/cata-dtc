package worker

import (
	"log"
	"sync"

	"github.com/robfig/cron/v3"
)

type Scheduler struct {
	cron   *cron.Cron
	jobs   map[string]cron.EntryID
	mu     sync.RWMutex
	logger *log.Logger
}

func NewScheduler() *Scheduler {
	return &Scheduler{
		cron:   cron.New(cron.WithSeconds()),
		jobs:   make(map[string]cron.EntryID),
		logger: log.Default(),
	}
}

func (s *Scheduler) AddJob(name string, spec string, job func()) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Remove existing job if it exists
	if entryID, exists := s.jobs[name]; exists {
		s.cron.Remove(entryID)
		s.logger.Printf("Removed existing job: %s", name)
	}

	entryID, err := s.cron.AddFunc(spec, job)
	if err != nil {
		return err
	}

	s.jobs[name] = entryID
	s.logger.Printf("Added job: %s with schedule: %s", name, spec)
	return nil
}

func (s *Scheduler) RemoveJob(name string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if entryID, exists := s.jobs[name]; exists {
		s.cron.Remove(entryID)
		delete(s.jobs, name)
		s.logger.Printf("Removed job: %s", name)
	}
}

func (s *Scheduler) Start() {
	s.cron.Start()
	s.logger.Println("Scheduler started")
}

func (s *Scheduler) Stop() {
	s.cron.Stop()
	s.logger.Println("Scheduler stopped")
}

func (s *Scheduler) IsRunning() bool {
	return len(s.cron.Entries()) > 0
}

func (s *Scheduler) ListJobs() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	jobs := make([]string, 0, len(s.jobs))
	for name := range s.jobs {
		jobs = append(jobs, name)
	}
	return jobs
}