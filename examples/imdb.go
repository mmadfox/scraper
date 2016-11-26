package main

import (
	"log"
	"net/http"
	"time"

	"github.com/mmadfox/scraper"
)

func main() {
	log.Println("Imdb scraper")
	var wc scraper.WorkerCount = 20
	s, err := scraper.New("http://www.imdb.com/trailers/", wc)
	if err != nil {
		panic(err)
	}
	p := "/title/tt{id:[0-9]+}/"
	s.Mux().HandleFunc(p, func(rw http.ResponseWriter, r *http.Request) {
		ctx := rw.(*scraper.Context)
		title := ctx.Doc.Find("h1[itemprop=name]").Text()
		log.Println(title)
	})
	go func() {
		for {
			select {
			case <-time.After(3 * time.Second):
				s.Stop()
				return
			}
		}
	}()
	s.Run()
}
