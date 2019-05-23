package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/gocolly/colly"
)

func main() {
	asyncScan := flag.Bool("async", false, "Run async requests")
	delay := flag.Int64("d", 0, "Delay between requests in millisecond")
	maxDepth := flag.Int("max-depth", 0, "MaxDepth limits the recursion depth of visited URLs.")
	parallelismScan := flag.Int("p", 0, "Parallelism is the number of the maximum allowed concurrent requests")
	scrapURL := flag.String("url", "http://example.com", "URL for scraping")
	userAgent := flag.String("user-agent", "Sicra crawler, https://github.com/jidckii/sicra", "User Agent")
	verbose := flag.Bool("v", true, "Verbose visiting URL")
	flag.Parse()

	parsedUrl, err := url.Parse(*scrapURL)
	if err != nil {
		log.Fatal(err)
	}
	if parsedUrl.Hostname() == "example" {
		log.Fatalln("Change default -url key")
	}

	var urlCount int = 0

	c := colly.NewCollector(
		colly.UserAgent(*userAgent),
		colly.AllowedDomains(parsedUrl.Hostname()),
		colly.MaxDepth(*maxDepth),
		colly.Async(*asyncScan),
	)

	c.Limit(&colly.LimitRule{
		DomainGlob:  parsedUrl.Hostname() + "/*",
		Delay:       time.Duration(*delay) * time.Millisecond,
		Parallelism: *parallelismScan,
	})

	if *verbose {
		c.OnRequest(func(r *colly.Request) {
			fmt.Println("Visiting", r.URL.String())
		})
	}

	// meta name="googlebot" content="noindex"

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		c.Visit(e.Request.AbsoluteURL(e.Attr("href")))
	})

	c.OnHTML("html", func(e *colly.HTMLElement) {
		metaNoindex := e.ChildAttr(`meta[name="googlebot"]`, "content")
		if metaNoindex != "noindex" {
			urlCount++
		}
	})

	// c.OnResponse(func(r *colly.Response) {
	// 	fmt.Println("Visited", r.Request.URL.String())
	// })

	// Start scraping
	c.Visit(*scrapURL)
	c.Wait()
	fmt.Println(urlCount)
}
