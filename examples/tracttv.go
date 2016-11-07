package main

import (
	"github.com/mmadfox/scraper"
	"log"
	"net/http"
	"time"
)

func main() {
	log.Println("Trakt Tv scraper")
	var wc scraper.WorkerCount = 5
	h := http.Header{}
	h.Add("Referer", "https://trakt.tv")

	s, err := scraper.New("https://trakt.tv/movies/trending", wc)
	if err != nil {
		panic(err)
	}
	p := `/movies/{movieName:(.*)\-[0-9]+}`
	s.Mux().HandleFunc(p, func(rw http.ResponseWriter, r *http.Request) {
		ctx := rw.(*scraper.Context)
		log.Println("Got the url", ctx.Addr.String())
	})
	s.SetHeader(h)
	go func() {
		for {
			select {
			case <-time.After(1 * time.Minute):
				s.StopAndClose()
				return
			}
		}
	}()
	s.Start()
	s.Block()
	log.Println("Stop")
}
