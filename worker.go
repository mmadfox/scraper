package scraper

import (
	"github.com/gorilla/mux"
	"net/url"
	"time"
)

type WorkerCount int

type Job struct {
	Payload *url.URL
	id      string
}

func (j Job) Id() string {
	if len(j.id) == 0 {
		j.id = GetMD5Hash(j.Payload.String())
	}
	return j.id
}

type Worker struct {
	Name       string
	queue      chan<- Job
	WorkerPool chan chan Job
	JobChannel chan Job
	router     *mux.Router
	quit       chan bool
	fetch      Fetcher
	pause      chan time.Duration
}

func newWorker(o workerOptions) Worker {
	return Worker{
		Name:       o.Name,
		WorkerPool: o.Pool,
		queue:      o.Queue,
		pause:      make(chan time.Duration),
		router:     o.Router,
		fetch:      o.Fetcher,
		JobChannel: make(chan Job),
		quit:       make(chan bool)}
}

func (w Worker) SetFetcher(fn Fetcher) {
	w.fetch = fn
}

func (w Worker) Pause(d time.Duration) {
	go func() {
		w.pause <- d
	}()
}

func (w Worker) Start() {
	go func() {
		for {
			w.WorkerPool <- w.JobChannel

			select {
			case dur := <-w.pause:
				time.Sleep(dur)
			case job := <-w.JobChannel:
				if job.Payload != nil {
					html, err := w.fetch.Fetch(job.Payload)
					if err == nil {
					}
				}
			case <-w.quit:
				return
			}
		}
	}()
}

func (w Worker) Stop() {
	go func() {
		w.quit <- true
	}()
}