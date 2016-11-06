package scraper

import (
	"github.com/mmadfox/scraper"
	"testing"
)

func TestGetRandomUserAgent(t *testing.T) {
	for i := 0; i < 100; i++ {
		ua := scraper.RandomUserAgent()
		if len(ua) == 0 {
			t.Errorf("empty ua %s", ua)
		}
	}
}
