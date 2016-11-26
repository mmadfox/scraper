package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/mmadfox/scraper"
)

type Movie struct {
	Name   string `json: "name"`
	Poster string `json: "poster"`
}

func main() {
	log.Println("Trakt Tv scraper")
	storage := make([]Movie, 0)
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
		name := ctx.Doc.Find(".mobile-title").Find("h1").Text()
		poster, _ := ctx.Doc.Find("img.real").Attr("data-original")
		storage = append(storage, Movie{
			Name:   name,
			Poster: poster,
		})
	})
	s.SetHeader(h)
	go func() {
		for {
			select {
			case <-time.After(30 * time.Second):
				s.Stop()
				return
			}
		}
	}()
	go func() {
		for {
			select {
			case <-time.After(10 * time.Second):
				s.Pause()
				return
			}
		}
	}()
	s.Run()
	js, _ := json.Marshal(storage)
	log.Print(string(js))
}
