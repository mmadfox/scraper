package scraper

import (
	"errors"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var (
	ErrBadRequest  = errors.New("Bad request")
	ErrMakeRequest = errors.New("Make request error")
)

type FactoryRequest func(u *url.URL) (*http.Request, error)

type Fetcher interface {
	MakeRequest(*url.URL) (*http.Request, error)
	Fetch(*http.Request) (*http.Response, error)
}

type defaultFetcher struct {
	c         *http.Client
	fr        FactoryRequest
	lastUrl   *url.URL
	userAgent string
}

func isBadRequest(resp *http.Response) bool {
	ct := resp.Header.Get("Content-Type")
	if !strings.Contains(ct, "text/html") || resp.StatusCode != 200 {
		return true
	}
	return false
}

func (f *defaultFetcher) MakeRequest(u *url.URL) (req *http.Request, err error) {
	if req, err = f.fr(u); err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Connection", "keep-alive")
	if f.lastUrl != nil {
		req.Header.Set("Referer", f.lastUrl.String())
	}
	req.Header.Set("User-Agent", f.userAgent)
	f.lastUrl = u
	return req, err
}

func (f *defaultFetcher) Fetch(req *http.Request) (resp *http.Response, err error) {
	resp, err = f.c.Do(req)
	if err != nil {
		return nil, ErrBadRequest
	}
	if ok := isBadRequest(resp); ok {
		return nil, ErrBadRequest
	}
	return resp, nil
}

func NewFetcher(cli *http.Client, fr FactoryRequest) Fetcher {
	return &defaultFetcher{
		c:         cli,
		fr:        fr,
		userAgent: RandomUserAgent(),
	}
}

func NewDefaultFetcher() Fetcher {
	cli := &http.Client{
		Timeout: time.Duration(10 * time.Second),
	}
	fetcher := NewFetcher(cli, func(u *url.URL) (*http.Request, error) {
		req, err := http.NewRequest("GET", u.String(), nil)
		return req, err
	})
	return fetcher
}

func IsBadRequest(err error) bool {
	return err == ErrBadRequest
}
