package scraper

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"net/url"
	"strings"
)

var (
	ErrBadResource = errors.New("Bad resource")
)

func GetMD5Hash(text string) string {
	text = strings.TrimRight(text, "/")
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

type Context struct {
	Doc   *goquery.Document
	links map[string]*url.URL
	Addr  *url.URL
	Res   *http.Response
	Req   *http.Request
}

func (p *Context) Header() http.Header {
	return make(http.Header)
}

func (p *Context) Write([]byte) (int, error) {
	return 0, nil
}

func (p *Context) WriteHeader(c int) {
}

func (p *Context) Links() map[string]*url.URL {
	if len(p.links) > 0 {
		return p.links
	}
	p.Doc.Find("a[href]").Each(func(i int, s *goquery.Selection) {
		if val, ok := s.Attr("href"); ok {
			fu := fixUrl(val, p.Addr)
			if fu.Host == p.Addr.Host {
				sUrl := fu.String()
				id := GetMD5Hash(sUrl)
				if _, ok := p.links[id]; !ok {
					p.links[id] = fu
				}
			}
		}
	})
	return p.links
}

func fixUrl(u string, addr *url.URL) *url.URL {
	o, err := url.Parse(u)
	if err == nil {
		if o.Scheme == "" {
			o.Scheme = addr.Scheme
		}
		if o.Host == "" {
			o.Host = addr.Host
		}
	}
	return o
}
