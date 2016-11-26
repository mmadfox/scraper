package scraper

import (
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/gorilla/mux"
	"github.com/streamrail/concurrent-map"
)

type workerOptions struct {
	State   int
	Name    string
	Pool    chan chan Job
	Queue   chan Job
	Router  *mux.Router
	HttpCli Fetcher
}

type Scrape struct {
	r           *mux.Router
	workerCount WorkerCount
	queue       chan Job
	addr        *url.URL
	pool        chan chan Job
	done        chan bool
	workers     []*Worker
	dup         cmap.ConcurrentMap
	run         bool
	HttpCli     Fetcher
	state       int
}

// Set the user agent string. By default, random string
func (s *Scrape) SetUserAgentString(ua string) {
	s.HttpCli.SetUserAgentString(ua)
}

// Set HTTP headers
func (s *Scrape) SetHeader(h http.Header) {
	s.HttpCli.SetHeader(h)
}

// Set HTTP Referer. By default, domain name with scheme, port, etc
func (s *Scrape) SetReferer(r string) {
	s.HttpCli.SetReferer(r)
}

// See gorillatoolkit.org/pkg/mux for details
// r := mmadfox.New("http://google.com")
// r.Mux().Host("{subdomain:[a-z]+}.domain.com").HandlerFunc(fn)
// r.Mux().Path("/video/{id:[0-9]+}/{article}/").Handler(h)
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

// Pause the queue for all workers
func (s *Scrape) Pause() {
	s.state = STATE_PAUSE
	if !s.IsRunning() {
		return
	}
	for _, w := range s.workers {
		w.Pause()
	}
}

// Start the queue for all  workers
func (s *Scrape) Start() {
	s.state = STATE_START
	if !s.IsRunning() {
		return
	}
	for _, w := range s.workers {
		w.Start()
	}
}

// Stop the queue for all workers
func (s *Scrape) Stop() {
	s.state = STATE_STOP
	if !s.IsRunning() {
		return
	}
	for _, w := range s.workers {
		<-w.Stop()
	}
	s.done <- true
}

func (s *Scrape) IsRunning() bool {
	return s.run == true
}

func (s *Scrape) Run() {
	if s.run == true {
		return
	}

	log.Printf("Run scraper by domain %s", s.addr.String())

	// default user agent string
	if len(s.HttpCli.UserAgentString()) == 0 {
		s.HttpCli.SetUserAgentString(RandomUserAgent())
	}
	// default referer
	if len(s.HttpCli.Referer()) == 0 {
		s.HttpCli.SetReferer(s.addr.String())
	}
	if s.state == 0 {
		s.state = STATE_START
	}
	s.run = true
	var wc WorkerCount
	for wc = 0; wc < s.workerCount; wc++ {
		w := newWorker(workerOptions{
			Name:    fmt.Sprintf("ID#%v", wc),
			State:   s.state,
			Pool:    s.pool,
			Queue:   s.queue,
			Router:  s.r,
			HttpCli: s.HttpCli})
		s.workers = append(s.workers, w)
		w.Do()
	}

	go s.dispatch()
	<-s.done
}

func New(domain string, wc WorkerCount) (*Scrape, error) {
	u, err := url.Parse(domain)
	if err != nil {
		return nil, err
	}
	if ok := wc.IsValid(); !ok {
		wc = 5
	}
	q := make(chan Job)
	p := make(chan chan Job, wc)
	go func() {
		q <- Job{Payload: u}
		return
	}()
	return &Scrape{
		HttpCli:     DefaultHttpCli{},
		r:           mux.NewRouter(),
		pool:        p,
		workerCount: wc,
		queue:       q,
		workers:     make([]*Worker, 0),
		dup:         cmap.New(),
		done:        make(chan bool),
		addr:        u}, nil
}
