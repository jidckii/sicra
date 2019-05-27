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
	delay := flag.Int64("delay", 0, "Delay between requests in Millisecond")
	maxDepth := flag.Int("max-depth", 0, "MaxDepth limits the recursion depth of visited URLs.")
	skipNoIndex := flag.Bool("skip-noindex", true, "Do not add link to sitemap if it contains: 'meta name = \"googlebot\" content = \"noindex\"'")
	outFile := flag.String("outfile", "./sitemap.xml", "Out sitemap file")
	paralScan := flag.Int("parallel", 0, "Parallelism is the number of the maximum allowed concurrent requests")
	scrapURL := flag.String("url", "http://go-colly.org/docs", "URL for scraping")
	timeoutResp := flag.Int("timeout", 10, "Response timeout in second")
	uriFilter := flag.String("uri-filter", "", "Filtering on uri prefix, example: /ru-ru , allowed use regex.")
	userAgent := flag.String("user-agent", "Sicra crawler, https://github.com/jidckii/sicra", "User Agent")
	verbose := flag.Bool("v", false, "Verbose visiting URL")
	flag.Parse()

	parsedURL, err := url.Parse(*scrapURL)
	if err != nil {
		log.Fatal(err)
	}

	scrape := sicra.Crawler(
		*scrapURL,
		*userAgent,
		parsedURL.Hostname(),
		*uriFilter,
		*paralScan,
		*maxDepth,
		*timeoutResp,
		*delay,
		*asyncScan,
		*skipNoIndex,
		*verbose)

	sicra.GenerateSiteMap(*outFile, scrape.AddedURLs)
	if *skipNoIndex {
		sicra.GenerateNoIndex("./noindex.txt", scrape.NoIndexURLs)
	}

	fmt.Print(
		"Request URLs: ", scrape.AllVisitURLsCount, "\n",
		"Added URLs ", scrape.AddedURLsCount, "\n",
		"No Index URLs ", scrape.NoIndexURLsCount, "\n",
		"Response URLs ", scrape.ResponseURLsCount, "\n",
		"Error URLs ", scrape.ErrorURLsCount, "\n",
	)
}
