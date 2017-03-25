package scraper

import (
	"log"
	"testing"
)

func TestVisitPage(t *testing.T) {
	browser, err := NewBrowser("http://kinogo.club")
	if err != nil {
		t.Fatal(err)
	}
	browser.Visit(`/{movieName:(.*)\-[0-9]+\.html}`, func(p Page) {
		log.Println(p)
	})

	if err := browser.Run(); err != nil {
		t.Fatal(err)
	}
}
