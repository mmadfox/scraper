package scrape

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/streamrail/concurrent-map"
	"net/url"
	"time"
)

const (
	PAUSE_15MIN = time.Minute * 15
	PAUSE_30MIN = time.Minute * 30
	PAUSE_60MIN = time.Minute * 60
)

type Scrape struct {
	r           *mux.Router
	workerCount int
	queue       chan Job
	workers     []Worker
	addr        *url.URL
	pool        chan chan Job
	done        chan bool
	dup         cmap.ConcurrentMap
	run         bool
}

func (s *Scrape) Mux() *mux.Router {
	return s.r
}

func (s *Scrape) dispatch() {
	for {
		select {
		case job := <-s.queue:
			if _, ok := s.dup.Get(job.Id()); !ok {
				s.dup.Set(job.Id(), true)
				go func(job Job) {
					jobChannel := <-s.pool
					jobChannel <- job
				}(job)
			}
		}
	}
}

func (s *Scrape) Pause(d time.Duration) {
	for _, w := range s.workers {
		w.Pause(d)
	}
}

func (s *Scrape) Start() *Scrape {
	if s.run == true {
		return nil
	}
	s.run = true
	for i := 0; i < s.workerCount; i++ {
		w := NewWorker(s.pool, s.queue, s.r, fmt.Sprintf("WorkerId: %d", i))
		w.Start()
		s.workers = append(s.workers, w)
	}

	go s.dispatch()
	return s
}

func (s *Scrape) Stop() {
	s.run = false
	for _, w := range s.workers {
		w.Stop()
	}
	close(s.queue)
	close(s.pool)
}

func (s *Scrape) Close() {
	go func() {
		s.done <- true
	}()
}

func (s *Scrape) Block() {
	<-s.done
}

func New(domain string, workerCount int) (*Scrape, error) {
	u, err := url.Parse(domain)
	if err != nil {
		return nil, err
	}
	if workerCount <= 0 {
		workerCount = 5
	}
	q := make(chan Job)
	p := make(chan chan Job, workerCount)
	go func() {
		q <- Job{Payload: u}
		return
	}()
	return &Scrape{
		r:           mux.NewRouter(),
		pool:        p,
		workerCount: workerCount,
		workers:     make([]Worker, workerCount),
		queue:       q,
		dup:         cmap.New(),
		done:        make(chan bool),
		addr:        u}, nil
}
