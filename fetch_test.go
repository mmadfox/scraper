package scraper

import (
	"github.com/mmadfox/scraper"
	"net/url"
	"testing"
)

func TestFetch(t *testing.T) {
	u, _ := url.Parse("http://imdb.com")
	f := scraper.DefaultFetcher{}
	html, err := f.Fetch(u)
	if err != nil {
		t.Error(err)
	}
	if len(html) == 0 {
		t.Errorf("empty reponse")
	}
	t.Log(html)
}

func TestBuildReferer(t *testing.T) {
	u, _ := url.Parse("http://imdb.com/velue/in/path/?do=123")
	got := u.Scheme + "://" + u.Host
	want := "http://imdb.com"
	if want != got {
		t.Errorf("want %s got %s", want, got)
	}
}
