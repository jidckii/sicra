package sicra

import (
	"fmt"
	"log"
	"net/url"
	"regexp"
	"strconv"
	"time"

	"github.com/gocolly/colly"
)

type scrapeURL struct {
	AddedURLs         []string
	AddedURLsCount    int
	AllVisitURLsCount int
	ErrorURLs         []string
	ErrorURLsCount    int
	NoIndexURLs       []string
	NoIndexURLsCount  int
	ResponseURLsCount int
}

// Crawler takes as input the parameters to scan.
// Returns URL scanning statistics
// and a list of links for sitemap generation.
// By default, pages containing 'meta name="robots" content="noindex,nofollow"' are ignored
func Crawler(
	scrapURL, userAgent, allowDomain, baseAuth, uriFilter, noIndexRule string,
	paralScan, maxDepth, timeoutResp int,
	delay int64,
	asyncScan, addError, skipNoIndex, verbose bool,
) *scrapeURL {

	scrapeURLs := new(scrapeURL)
	allowDomainFilter := "http(s)?://" + allowDomain
	if baseAuth != "" {
		allowDomainFilter = "http(s)?://" + baseAuth + "@" + allowDomain
	}
	filter := allowDomainFilter
	if uriFilter != "" {
		filter = allowDomainFilter + uriFilter
	}
	c := colly.NewCollector(
		colly.UserAgent(userAgent),
		colly.AllowedDomains(allowDomain),
		colly.MaxDepth(maxDepth),
		colly.Async(asyncScan),
		colly.URLFilters(
			regexp.MustCompile(allowDomainFilter+"(/)?$"),
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
		requestURL := urlEscape(er.Request.URL.String())
		r := regexp.MustCompile("^5[0-9]{1,2}$")
		statusCode := strconv.Itoa(er.StatusCode)
		strErr := fmt.Sprint(err)
		if verbose {
			log.Println("Error:", err, requestURL)
		}
		if addError {
			if r.MatchString(statusCode) {
				add(requestURL, verbose, scrapeURLs)
			}
		}
		scrapeURLs.ErrorURLsCount++
		scrapeURLs.ErrorURLs = append(scrapeURLs.ErrorURLs, statusCode+" "+strErr+" "+requestURL)
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

	// meta name="robots" content="noindex,nofollow"
	c.OnHTML("html", func(e *colly.HTMLElement) {
		requestURL := urlEscape(e.Request.URL.String())
		if skipNoIndex {
			metaNoindex := e.ChildAttr(`meta[name="robots"]`, "content")
			if metaNoindex != noIndexRule {
				add(requestURL, verbose, scrapeURLs)
			} else {
				scrapeURLs.NoIndexURLsCount++
				scrapeURLs.NoIndexURLs = append(scrapeURLs.NoIndexURLs, requestURL)
				if verbose {
					log.Println("Skiped: " + requestURL)
				}
			}
		} else {
			add(requestURL, verbose, scrapeURLs)
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
	pathuri := parseURL.Path
	query := url.QueryEscape(parseURL.RawQuery)
	if query != "" {
		query = "?" + query
	}
	fragment := url.QueryEscape(parseURL.Fragment)
	if fragment != "" {
		fragment = "#" + fragment
	}
	escapeURL := (scheme + "://" + host + pathuri + query + fragment)
	return escapeURL
}
