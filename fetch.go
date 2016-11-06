package scraper

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

var (
	ErrBadRequest = errors.New("Bad request")
)

type Fetcher interface {
	SetHeader(http.Header)
	SetUserAgent(string)
	SetReferer(string)
	UserAgent() string
	Referer() string
	Fetch(*url.URL) (*Context, error)
}

type DefaultFetcher struct {
	header    http.Header
	userAgent string
	referer   string
}

func (f DefaultFetcher) SetHeader(h http.Header) {
	f.header = h
}

func (f DefaultFetcher) SetUserAgent(u string) {
	f.userAgent = u
}

func (f DefaultFetcher) SetReferer(r string) {
	f.referer = r
}

func (f DefaultFetcher) UserAgent() string {
	return f.userAgent
}

func (f DefaultFetcher) Referer() string {
	return f.referer
}

func (f DefaultFetcher) Fetch(u *url.URL) (string, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return "", err
	}
	req.Header = f.header
	if len(f.userAgent) > 0 {
		req.Header.Add("User-Agent", f.userAgent)
	}
	if len(f.referer) > 0 {
		req.Header.Add("Referer", f.referer)
	}
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	if !strings.Contains(res.Header.Get("Content-Type"), "text/html") {
		return "", ErrBadRequest
	}
	if !strings.Contains(res.Header.Get("Content-Type"), "text/html") {
		return "", ErrBadRequest
	}
	if res.StatusCode != 200 {
		return "", ErrBadRequest
	}
	htmlData, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	return string(htmlData), nil
}
