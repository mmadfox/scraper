package scraper

import "testing"

func TestUa(t *testing.T) {
	for v := 0; v < 100; v++ {
		userAgent := RandomUserAgent()
		if len(userAgent) == 0 {
			t.Fatalf("RandomUserAgent() = %v, want not empty", userAgent)
		}
	}
}
