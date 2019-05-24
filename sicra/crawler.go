package sicra

import (
	"fmt"
	"log"
	"time"

	"github.com/gocolly/colly"
)

type scrapeCounts struct {
	AllVisitURLs int
	NoIndexURLs  int
	AddedURLs    int
	ResponseURLs int
	ErrorURLs    int
}

// Принимает на вход параметры для сканирования, возвращает колличесво посещённых, пропущенных к добавлению и добавленных.
func Crawler(
	scrapURL, userAgent, allowDomain string,
	paralScan, maxDepth, timeoutResp int,
	delay int64,
	asyncScan, verbose bool,
) *scrapeCounts {

	ulrCounts := scrapeCounts{}

	c := colly.NewCollector(
		colly.UserAgent(userAgent),
		colly.AllowedDomains(allowDomain),
		colly.MaxDepth(maxDepth),
		colly.Async(asyncScan),
	)
	c.SetRequestTimeout(time.Duration(timeoutResp) * time.Second)

	c.Limit(&colly.LimitRule{
		DomainGlob:  allowDomain + "/*",
		Delay:       time.Duration(delay) * time.Second,
		RandomDelay: time.Duration(delay) * time.Second,
		Parallelism: paralScan,
	})

	c.OnRequest(func(r *colly.Request) {
		if verbose {
			log.Println("Request: ", r.URL.String())
		}
		ulrCounts.AllVisitURLs++
	})

	c.OnError(func(_ *colly.Response, err error) {
		if verbose {
			log.Println("Error:", err)
		}
		ulrCounts.ErrorURLs++
	})

	c.OnResponse(func(re *colly.Response) {
		ulrCounts.ResponseURLs++
		fmt.Println("Response: " + re.Request.URL.String())
	})

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		c.Visit(e.Request.AbsoluteURL(e.Attr("href")))
	})

	// meta name="googlebot" content="noindex"
	c.OnHTML("html", func(e *colly.HTMLElement) {
		metaNoindex := e.ChildAttr(`meta[name="googlebot"]`, "content")
		if metaNoindex != "noindex" {
			ulrCounts.AddedURLs++
			if verbose {
				// log.Println("Added: " + e.Request.URL.String())
			}
		} else {
			ulrCounts.NoIndexURLs++
			if verbose {
				// log.Println("Skiped: " + e.Request.URL.String())
			}
		}
	})

	// Start scraping
	c.Visit(scrapURL)
	c.Wait()

	return &ulrCounts
}
