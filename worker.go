package scrape

import (
	"github.com/gorilla/mux"
	"log"
	"net/url"
	"time"
)

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
	pause      chan time.Duration
}

func NewWorker(workerPool chan chan Job, queue chan<- Job, r *mux.Router, name string) Worker {
	return Worker{
		Name:       name,
		WorkerPool: workerPool,
		queue:      queue,
		pause:      make(chan time.Duration),
		router:     r,
		JobChannel: make(chan Job),
		quit:       make(chan bool)}
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
				}
			case <-w.quit:
				return
			}
		}
	}()
}

func (w Worker) Stop() {
	go func() {
		log.Println("Worker stop.", w.Name)
		w.quit <- true
	}()
}
