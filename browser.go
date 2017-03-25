package scraper

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/gorilla/mux"
)

type Browser struct {
	Router    *mux.Router
	Client    Fetcher
	Visits    Visiter
	FindLinks Matcher
	domain    *url.URL
	queue     chan *url.URL
	pool      chan chan *url.URL
	run       bool
	wc        int
	state     int
	sleep     int
}

func NewBrowser(domain string) (*Browser, error) {
	u, err := url.Parse(domain)
	if err != nil {
		return nil, err
	}
	wc := 12
	router := mux.NewRouter()
	queue := make(chan *url.URL, 100)
	pool := make(chan chan *url.URL, wc)
	browser := &Browser{
		Router:    router,
		Client:    NewDefaultFetcher(),
		Visits:    NewMemoryVisits(),
		FindLinks: Match,
		domain:    u,
		queue:     queue,
		pool:      pool,
		run:       false,
		wc:        wc,
		state:     1,
		sleep:     0,
	}
	return browser, nil
}

func (b *Browser) Play() {
	b.state = 1
}

func (b *Browser) Pause() {
	b.state = 0
}

func (b *Browser) Sleep(sec int) {
	b.sleep = sec
}

func (b *Browser) workerRun() {
	worker := make(chan *url.URL)

	for {
		if b.state == 0 {
			time.Sleep(time.Millisecond * 500)
			runtime.Gosched()
			continue
		}

		b.pool <- worker

		u, ok := <-worker
		if !ok {
			continue
		}

		if err := func() error {
			content, req, resp, err := b.fetchContent(u)
			if err != nil {
				return nil
			}

			rw := fakeResponseWriter{
				Request:  req,
				Response: resp,
				Content:  content,
				Addr:     u,
			}

			log.Println(u.String())
			b.Router.ServeHTTP(rw, req)

			links, err := b.FindLinks(content)
			if err != nil {
				return err
			}
			for _, link := range links {
				u, err := b.makeUrl(link)
				if err != nil {
					continue
				}
				if b.domain.Host != u.Host {
					continue
				}

				b.queue <- u
			}

			if b.sleep > 0 {
				time.Sleep(time.Duration(b.sleep) * time.Second)
			}

			return nil
		}(); err != nil {
			log.Printf("Error: %v", err)
		}
	}
}

func (b *Browser) makeUrl(u string) (*url.URL, error) {
	trg, err := url.Parse(u)
	if err != nil {
		return nil, err
	}
	if trg.Scheme == "" {
		trg.Scheme = b.domain.Scheme
	}
	if trg.Host == "" {
		trg.Host = b.domain.Host
	}
	return trg, nil
}

func (b *Browser) fetchContent(u *url.URL) (string, *http.Request, *http.Response, error) {
	req, _ := b.Client.MakeRequest(u)
	resp, err := b.Client.Fetch(req)
	if err != nil {
		return "", nil, nil, err
	}
	html, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", nil, nil, err
	}
	return string(html), req, resp, nil
}

func (b *Browser) Run() error {
	if b.run {
		return errors.New("Already running")
	}

	b.run = true
	errs := make(chan error)

	for w := 0; w <= b.wc; w++ {
		go b.workerRun()
	}

	go b.queueRun()
	go func() {
		b.queue <- b.domain
	}()

	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	return <-errs
}

func (b *Browser) queueRun() {
	for {
		select {
		case u := <-b.queue:
			if isVisited := b.Visits.Visit(u.String()); isVisited {
				continue
			}

			go func(u *url.URL) {
				worker, ok := <-b.pool
				if !ok {
					panic("Cannot get worker from pool")
				}
				worker <- u
			}(u)
		}
	}
}

func (b *Browser) Visit(path string, f func(Page)) {
	b.Router.HandleFunc(path, func(rw http.ResponseWriter, r *http.Request) {
		frw := rw.(fakeResponseWriter)
		f(Page{frw})
	})
}
