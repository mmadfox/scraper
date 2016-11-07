package scraper

import (
	"github.com/mmadfox/scraper"
	"net/url"
	"testing"
)

func TestFetch(t *testing.T) {
	u, _ := url.Parse("http://www.apple.com/us/shop/goto/account")
	f := scraper.DefaultFetcher{}
	resp, _, err := f.Fetch(u)
	if err != nil {
		t.Error(err)
	}
	t.Log(resp.StatusCode)
}

func TestBuildReferer(t *testing.T) {
	u, _ := url.Parse("http://imdb.com/velue/in/path/?do=123")
	got := u.Scheme + "://" + u.Host
	want := "http://imdb.com"
	if want != got {
		t.Errorf("want %s got %s", want, got)
	}
}
