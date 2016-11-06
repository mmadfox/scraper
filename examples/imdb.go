package main

import (
	"github.com/mmadfox/scraper"
	"log"
	"net/http"
)

func main() {
	log.Println("Imdb scraper")
	s, err := scraper.New("http://imdb.com", 20)
	if err != nil {
		panic(err)
	}
	p := "/title/tt{id:[0-9]+}/"
	s.Mux().HandleFunc(p, func(rw http.ResponseWriter, r *http.Request) {
		ctx := rw.(*scraper.Context)
		log.Println("Got the url", ctx.Addr.String())
	})
	s.Start()
	s.Block()
}
