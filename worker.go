package scraper

import (
	"log"
	"net/url"
	"runtime"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gorilla/mux"
)

const (
	STATE_START int = iota
	STATE_STOP
	STATE_PAUSE
)

type WorkerCount int

func (w WorkerCount) IsValid() bool {
	return w > 0 && w < 200
}

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
	Name         string
	queue        chan<- Job
	WorkerPool   chan chan Job
	JobChannel   chan Job
	router       *mux.Router
	httpcli      Fetcher
	stopChannel  chan bool
	stateChannel chan int
	state        int
}

func newWorker(o workerOptions) *Worker {
	return &Worker{
		Name:         o.Name,
		WorkerPool:   o.Pool,
		queue:        o.Queue,
		state:        o.State,
		stateChannel: make(chan int),
		router:       o.Router,
		httpcli:      o.HttpCli,
		stopChannel:  make(chan bool),
		JobChannel:   make(chan Job),
	}
}

func (w *Worker) SetHttpCli(f Fetcher) {
	w.httpcli = f
}

func (w *Worker) Pause() {
	go func() {
		w.stateChannel <- STATE_PAUSE
	}()
}

func (w *Worker) Stop() chan bool {
	go func() {
		w.stateChannel <- STATE_STOP
	}()
	return w.stopChannel
}

func (w *Worker) Start() {
	go func() {
		w.stateChannel <- STATE_START
	}()
}

func (w *Worker) Do() {
	log.Println("Do worker", w.Name)
	state := w.state

	go func() {
		for {

			select {
			case state = <-w.stateChannel:
				switch state {
				case STATE_PAUSE:
					log.Println("State pause")
				case STATE_STOP:
					log.Println("State stop")
				case STATE_START:
					log.Println("State start")
				}
			default:
				if state == STATE_PAUSE {
					time.Sleep(time.Millisecond * 250)
					runtime.Gosched()
					continue
				}

				if state == STATE_STOP {
					log.Printf("Stop worker %s", w.Name)
					w.stopChannel <- true
					close(w.stopChannel)
					return
				}

				w.WorkerPool <- w.JobChannel

				job, ok := <-w.JobChannel

				if !ok {
					continue
				}

				log.Printf("Scan url: %s, worker: %s",
					job.Payload.String(), w.Name)

				resp, req, err := w.httpcli.Fetch(job.Payload)
				if err != nil {
					log.Printf("Error: %v", err)
					continue
				}

				ctx := &Context{
					Addr:  job.Payload,
					Res:   resp,
					Req:   req,
					links: make(map[string]*url.URL, 0),
				}

				doc, err := goquery.NewDocumentFromResponse(resp)
				if err != nil {
					log.Printf("Error: %v")
					continue
				}
				ctx.Doc = doc
				w.router.ServeHTTP(ctx, req)
				for _, link := range ctx.Links() {
					w.queue <- Job{Payload: link}
				}
			}
		}
	}()
}
