package sicra

import (
	"log"
	"net/url"
	"regexp"
	"time"

	"github.com/gocolly/colly"
)

type scrapeURL struct {
	AddedURLs         []string
	AddedURLsCount    int
	AllVisitURLsCount int
	ErrorURLsCount    int
	NoIndexURLsCount  int
	NoIndexURLs       []string
	ResponseURLsCount int
}

// Crawler takes as input the parameters to scan.
// Returns URL scanning statistics
// and a list of links for sitemap generation.
// By default, pages containing 'meta name = "googlebot" content = "noindex"' are ignored
func Crawler(
	scrapURL, userAgent, allowDomain, uriFilter string,
	paralScan, maxDepth, timeoutResp int,
	delay int64,
	asyncScan, skipNoIndex, verbose bool,
) *scrapeURL {

	scrapeURLs := new(scrapeURL)
	filter := "http(s)?://" + allowDomain
	if uriFilter != "" {
		filter = "http(s)?://" + allowDomain + uriFilter
	}
	c := colly.NewCollector(
		colly.UserAgent(userAgent),
		colly.AllowedDomains(allowDomain),
		colly.MaxDepth(maxDepth),
		colly.Async(asyncScan),
		colly.URLFilters(
			regexp.MustCompile("http(s)?://"+allowDomain+"(/)?$"),
			regexp.MustCompile(filter),
		),
	)

	c.SetRequestTimeout(time.Duration(timeoutResp) * time.Second)

	c.Limit(&colly.LimitRule{
		DomainGlob:  allowDomain,
		Delay:       time.Duration(delay) * time.Millisecond,
		RandomDelay: time.Duration(delay) * time.Millisecond,
		Parallelism: paralScan,
	})

	c.OnRequest(func(r *colly.Request) {
		if verbose {
			log.Println("Request:", r.URL.String())
		}
		scrapeURLs.AllVisitURLsCount++
	})

	c.OnError(func(er *colly.Response, err error) {
		if verbose {
			log.Println("Error:", err, er.Request.URL.String())
		}
		scrapeURLs.ErrorURLsCount++
	})

	c.OnResponse(func(re *colly.Response) {
		scrapeURLs.ResponseURLsCount++
		if verbose {
			log.Println("Response: " + re.Request.URL.String())
		}
	})

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		c.Visit(e.Request.AbsoluteURL(e.Attr("href")))
	})

	// meta name="googlebot" content="noindex"
	c.OnHTML("html", func(e *colly.HTMLElement) {
		requesturl := urlEscape(e.Request.URL.String())
		if skipNoIndex {
			metaNoindex := e.ChildAttr(`meta[name="googlebot"]`, "content")
			if metaNoindex != "noindex" {
				add(requesturl, verbose, scrapeURLs)
			} else {
				scrapeURLs.NoIndexURLsCount++
				scrapeURLs.NoIndexURLs = append(scrapeURLs.NoIndexURLs, requesturl)
				if verbose {
					log.Println("Skiped: " + requesturl)
				}
			}
		} else {
			add(requesturl, verbose, scrapeURLs)
		}
	})

	// Start scraping
	c.Visit(scrapURL)
	c.Wait()

	return scrapeURLs
}

// Add AddedURLs in struct scrapeURL
func add(url string, verbose bool, scrapeURLs *scrapeURL) {
	scrapeURLs.AddedURLsCount++
	scrapeURLs.AddedURLs = append(scrapeURLs.AddedURLs, url)
	if verbose {
		log.Println("Added: " + url)
	}
}

// Escape URL
func urlEscape(refurl string) string {
	parseURL, err := url.Parse(refurl)
	if err != nil {
		log.Fatal(err)
	}
	scheme := parseURL.Scheme
	host := parseURL.Host
	pathuri := parseURL.EscapedPath()
	query := url.QueryEscape(parseURL.RawQuery)

	escapeURL := (scheme + "://" + host + pathuri + "?" + query)
	return escapeURL
}
