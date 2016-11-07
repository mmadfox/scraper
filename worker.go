<<<<<<< HEAD
package scraper

import (
	"github.com/gorilla/mux"
=======
package scrape

import (
	"github.com/gorilla/mux"
	"log"
>>>>>>> c38b59f1421a599579e1ff7b808c28655add4f01
	"net/url"
	"time"
)

<<<<<<< HEAD
type WorkerCount int

=======
>>>>>>> c38b59f1421a599579e1ff7b808c28655add4f01
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
<<<<<<< HEAD
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
=======
	pause      chan time.Duration
}

func NewWorker(workerPool chan chan Job, queue chan<- Job, r *mux.Router, name string) Worker {
	return Worker{
		Name:       name,
		WorkerPool: workerPool,
		queue:      queue,
		pause:      make(chan time.Duration),
		router:     r,
>>>>>>> c38b59f1421a599579e1ff7b808c28655add4f01
		JobChannel: make(chan Job),
		quit:       make(chan bool)}
}

<<<<<<< HEAD
func (w Worker) SetFetcher(fn Fetcher) {
	w.fetch = fn
}

=======
>>>>>>> c38b59f1421a599579e1ff7b808c28655add4f01
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
<<<<<<< HEAD
				time.Sleep(dur)
			case job := <-w.JobChannel:
				if job.Payload != nil {
					html, err := w.fetch.Fetch(job.Payload)
					if err == nil {
					}
=======
				log.Println("pause")
				time.Sleep(dur)
			case job := <-w.JobChannel:
				ctx, err := Fetch(job.Payload)
				if err != nil {
					log.Println(err)
				} else {
					for _, l := range ctx.Links() {
						w.queue <- Job{Payload: l}
					}
					w.router.ServeHTTP(ctx, ctx.Req)
>>>>>>> c38b59f1421a599579e1ff7b808c28655add4f01
				}
			case <-w.quit:
				return
			}
		}
	}()
}

func (w Worker) Stop() {
	go func() {
<<<<<<< HEAD
=======
		log.Println("Worker stop.", w.Name)
>>>>>>> c38b59f1421a599579e1ff7b808c28655add4f01
		w.quit <- true
	}()
}
