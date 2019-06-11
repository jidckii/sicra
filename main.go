package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"path/filepath"

	"github.com/jidckii/sicra/sicra"
)

func main() {
	addError := flag.Bool("add-error", true, "Add URL to sitemap, even if response error (only for 5xx codes)")
	asyncScan := flag.Bool("async", false, "Run async requests")
	delay := flag.Int64("delay", 0, "Delay between requests in Millisecond")
	maxDepth := flag.Int("max-depth", 0, "MaxDepth limits the recursion depth of visited URLs.")
	noIndexRule := flag.String("noindex-rule", "noindex,nofollow", "Comma-separated list of parameters as a string")
	outFile := flag.String("outfile", "./sitemap.xml", "Out sitemap file")
	paralScan := flag.Int("parallel", 0, "Parallelism is the number of the maximum allowed concurrent requests")
	scrapURL := flag.String("url", "http://go-colly.org/docs", "URL for scraping")
	skipNoIndex := flag.Bool("skip-noindex", true, "Do not add link to sitemap if it contains: 'meta name=\"robots\" content=\"noindex,nofollow\"'")
	timeoutResp := flag.Int("timeout", 10, "Response timeout in second")
	uriFilter := flag.String("uri-filter", "", "Filtering on uri prefix, example: /ru-ru , allowed use regex.")
	userAgent := flag.String("user-agent", "Sicra crawler, https://github.com/jidckii/sicra", "User Agent")
	verbose := flag.Bool("v", false, "Verbose visiting URL")
	flag.Parse()

	parsedURL, err := url.Parse(*scrapURL)
	if err != nil {
		log.Fatal(err)
	}
	hostname := parsedURL.Hostname()
	baseAuth := parsedURL.User.String()

	scrape := sicra.Crawler(
		*scrapURL,
		*userAgent,
		hostname,
		baseAuth,
		*uriFilter,
		*noIndexRule,
		*paralScan,
		*maxDepth,
		*timeoutResp,
		*delay,
		*asyncScan,
		*addError,
		*skipNoIndex,
		*verbose)

	p := filepath.Dir(*outFile)
	// generate sitemap.xml
	if len(scrape.AddedURLs) > 0 {
		err = sicra.GenerateSiteMap(*outFile, scrape.AddedURLs)
		if err != nil {
			log.Fatal(err)
		}
	}
	// generate noindex.txt
	if *skipNoIndex {
		if len(scrape.NoIndexURLs) > 0 {
			err = sicra.GenerateTxt(p+"/noindex.txt", scrape.NoIndexURLs)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
	// generate errors.txt
	if len(scrape.ErrorURLs) > 0 {
		err = sicra.GenerateTxt(p+"/errors.txt", scrape.ErrorURLs)
		if err != nil {
			log.Fatal(err)
		}
	}
	// print stats
	if *verbose {
		fmt.Print(
			"Request URLs: ", scrape.AllVisitURLsCount, "\n",
			"Added URLs ", scrape.AddedURLsCount, "\n",
			"No Index URLs ", scrape.NoIndexURLsCount, "\n",
			"Response URLs ", scrape.ResponseURLsCount, "\n",
			"Error URLs ", scrape.ErrorURLsCount, "\n",
		)
	}
}
