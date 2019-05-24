package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"

	"github.com/jidckii/sicra/sicra"
)

func main() {
	asyncScan := flag.Bool("async", false, "Run async requests")
	delay := flag.Int64("delay", 0, "Delay between requests in second")
	maxDepth := flag.Int("max-depth", 0, "MaxDepth limits the recursion depth of visited URLs.")
	paralScan := flag.Int("parallel", 0, "Parallelism is the number of the maximum allowed concurrent requests")
	scrapURL := flag.String("url", "http://example.com", "URL for scraping")
	timeoutResp := flag.Int("timeout", 10, "Response timeout in second")
	userAgent := flag.String("user-agent", "Sicra crawler, https://github.com/jidckii/sicra", "User Agent")
	verbose := flag.Bool("v", false, "Verbose visiting URL")
	flag.Parse()

	parsedURL, err := url.Parse(*scrapURL)
	if err != nil {
		log.Fatal(err)
	}
	hostname := parsedURL.Hostname()

	if hostname == "example.com" {
		log.Println("Change default URL example.com, -url key, -h for help")
		log.Fatal()
	} else {

		counts := sicra.Crawler(
			*scrapURL,
			*userAgent,
			hostname,
			*paralScan,
			*maxDepth,
			*timeoutResp,
			*delay,
			*asyncScan,
			*verbose)

		fmt.Println("Request URLs: ", counts.AllVisitURLs)
		fmt.Println("Added URLs", counts.AddedURLs)
		fmt.Println("No Index URLs", counts.NoIndexURLs)
		fmt.Println("Response URLs", counts.ResponseURLs)
		fmt.Println("Error URLs", counts.ErrorURLs)
	}
}
