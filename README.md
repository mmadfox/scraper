# scraper
A fast and powerful the scraper for html web pages

## Installation
 $ go get github.com/mmadfox/scraper
 
## Examples
### http://kinogo.club
```Go
package main

import (
	"log"

	"github.com/mmadfox/scraper"
)

func main() {
	browser, err := scraper.NewBrowser("http://kinogo.club")
        if err != nil {        
                log.Fatal(err)   
        }
        browser.Visit(`/{movieName:(.*)\-[0-9]+\.html}`, func(p Page) {
                log.Println(p) 
        })
  
        browser.Run()
}

