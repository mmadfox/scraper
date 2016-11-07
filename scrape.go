<<<<<<< HEAD
package scraper
=======
package scrape
>>>>>>> c38b59f1421a599579e1ff7b808c28655add4f01

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/streamrail/concurrent-map"
<<<<<<< HEAD
	"net/http"
=======
>>>>>>> c38b59f1421a599579e1ff7b808c28655add4f01
	"net/url"
	"time"
)

const (
	PAUSE_15MIN = time.Minute * 15
	PAUSE_30MIN = time.Minute * 30
	PAUSE_60MIN = time.Minute * 60
)

<<<<<<< HEAD
type workerOptions struct {
	Name    string
	Pool    chan chan Job
	Queue   chan Job
	Router  *mux.Router
	Fetcher Fetcher
}

type Scrape struct {
	r           *mux.Router
	workerCount WorkerCount
=======
type Scrape struct {
	r           *mux.Router
	workerCount int
>>>>>>> c38b59f1421a599579e1ff7b808c28655add4f01
	queue       chan Job
	workers     []Worker
	addr        *url.URL
	pool        chan chan Job
	done        chan bool
	dup         cmap.ConcurrentMap
	run         bool
<<<<<<< HEAD
	Fetcher     Fetcher
}

func (s *Scrape) SetUserAgent(ua string) {
	s.Fetcher.SetUserAgent(ua)
}

func (s *Scrape) SetHeader(h http.Header) {
	s.Fetcher.SetHeader(h)
}

func (s *Scrape) SetReferer(r string) {
	s.Fetcher.SetReferer(r)
=======
>>>>>>> c38b59f1421a599579e1ff7b808c28655add4f01
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
<<<<<<< HEAD
	//default user agent string
	if len(s.Fetcher.UserAgent()) == 0 {
		s.Fetcher.SetUserAgent(RandomUserAgent())
	}
	s.run = true
	var wc WorkerCount
	for wc = 0; wc < s.workerCount; wc++ {
		w := newWorker(workerOptions{
			Name:    fmt.Sprintf("WorkerId %v", wc),
			Pool:    s.pool,
			Queue:   s.queue,
			Router:  s.r,
			Fetcher: s.Fetcher})
		w.Start()
		s.workers = append(s.workers, w)
	}
=======
	s.run = true
	for i := 0; i < s.workerCount; i++ {
		w := NewWorker(s.pool, s.queue, s.r, fmt.Sprintf("WorkerId: %d", i))
		w.Start()
		s.workers = append(s.workers, w)
	}

>>>>>>> c38b59f1421a599579e1ff7b808c28655add4f01
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

<<<<<<< HEAD
func New(domain string, wc WorkerCount) (*Scrape, error) {
=======
func New(domain string, workerCount int) (*Scrape, error) {
>>>>>>> c38b59f1421a599579e1ff7b808c28655add4f01
	u, err := url.Parse(domain)
	if err != nil {
		return nil, err
	}
<<<<<<< HEAD
	if wc <= 0 {
		wc = 5
	}
	q := make(chan Job)
	p := make(chan chan Job, wc)
=======
	if workerCount <= 0 {
		workerCount = 5
	}
	q := make(chan Job)
	p := make(chan chan Job, workerCount)
>>>>>>> c38b59f1421a599579e1ff7b808c28655add4f01
	go func() {
		q <- Job{Payload: u}
		return
	}()
	return &Scrape{
<<<<<<< HEAD
		Fetcher:     DefaultFetcher{},
		r:           mux.NewRouter(),
		pool:        p,
		workerCount: wc,
		workers:     make([]Worker, wc),
=======
		r:           mux.NewRouter(),
		pool:        p,
		workerCount: workerCount,
		workers:     make([]Worker, workerCount),
>>>>>>> c38b59f1421a599579e1ff7b808c28655add4f01
		queue:       q,
		dup:         cmap.New(),
		done:        make(chan bool),
		addr:        u}, nil
}
