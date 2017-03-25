package scraper

import (
	"regexp"
	"strings"
)

type Matcher func(html string) ([]string, error)

var (
	tagPattern  = "(?i)<a([^>]+)>(.+?)</a>"
	hrefPattern = "\\s*(?i)href\\s*=\\s*(\"([^\"]*\")|'[^']*'|([^'\">\\s]+))"
	tagRegexp   *regexp.Regexp
	hrefRegexp  *regexp.Regexp
)

func init() {
	tagRegexp, _ = regexp.Compile(tagPattern)
	hrefRegexp, _ = regexp.Compile(hrefPattern)
}

func Match(html string) ([]string, error) {
	links := make([]string, 0)
	tags := tagRegexp.FindAllSubmatch([]byte(html), -1)

	for _, a := range tags {
		href := hrefRegexp.FindAllSubmatch(a[1], -1)
		for _, h := range href {
			tagHref := string(h[1])
			tagHref = strings.Trim(tagHref, "\"")
			links = append(links, tagHref)
		}
	}

	return links, nil
}
