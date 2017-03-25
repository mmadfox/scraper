package scraper

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

type fc struct {
	in   string
	out  error
	out2 error
}

var (
	cases = []fc{
		fc{"/index.html", nil, nil},
		fc{"/picture.png", ErrBadRequest, nil},
		fc{"/timeout", ErrBadRequest, nil},
		fc{"//////////", ErrBadRequest, nil},
		fc{"/index-500.html", ErrBadRequest, nil},
		fc{"/errors", ErrBadRequest, ErrMakeRequest},
	}
)

type fetcherHandler struct {
}

func (h *fetcherHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch u := r.URL.String(); u {
	case "/index.html":
		w.Header().Add("Content-Type", "text/html")
		w.WriteHeader(200)
		break
	case "/picture.png":
		w.Header().Add("Content-Type", "image/png")
		break
	case "/timeout":
		time.Sleep(time.Second * 2)
		break
	case "/index-500.html":
		w.WriteHeader(500)
		break
	}

	w.Write([]byte("response"))
}

func testCases(f Fetcher, u string, t *testing.T, cs []fc) {
	for _, c := range cs {
		u, err := url.Parse(u + c.in)
		if err != nil {
			t.Fatal(err)
		}
		req, err := f.MakeRequest(u)

		if err != c.out2 {
			t.Errorf("fetcher.MakeRequest(%q) => (%q, %q), want %q", u, req, err, c.out2)
		}
		if req == nil {
			continue
		}
		resp, err := f.Fetch(req)
		if err != c.out {
			t.Errorf("fetcher.Fetch(%q) => (%q, %q), want %q", req, resp, err, c.out)
		}
	}
}

func TestFetcherIsBadRequest(t *testing.T) {
	var errs = []struct {
		in  error
		out bool
	}{
		{ErrBadRequest, true},
		{errors.New("error"), false},
		{nil, false},
	}

	for _, e := range errs {
		f := IsBadRequest(e.in)
		if f != e.out {
			t.Errorf("IsBadRequest(%q) => %q, want %q", e.in, f, e.out)
		}
	}
}

func TestFetcher(t *testing.T) {
	h := &fetcherHandler{}
	server := httptest.NewServer(h)
	defer server.Close()

	cli := &http.Client{
		Timeout: time.Duration(1 * time.Second),
	}
	fetcher := NewFetcher(cli, func(u *url.URL) (*http.Request, error) {
		if u.Path == "/errors" {
			return nil, ErrMakeRequest
		}
		req, err := http.NewRequest("GET", u.String(), nil)
		if err != nil {
			return nil, err
		}
		req.Header.Add("User-Agent", "test")
		req.Header.Add("Referer", "test")
		return req, err
	})

	testCases(fetcher, server.URL, t, cases)
}

func TestFetcherDefault(t *testing.T) {
	h := &fetcherHandler{}
	server := httptest.NewServer(h)
	defer server.Close()

	fetcher := NewDefaultFetcher()
	c := cases[:len(cases)-1]

	testCases(fetcher, server.URL, t, c)
}
