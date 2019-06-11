# sicra
Simple crawler and sitemap generator on golang

build:  
```
go get -u github.com/jidckii/sicra
make
```
Run:
```
build/sicra -h

Usage of build/sicra:
  -add-error
        Add URL to sitemap, even if response error (only for 5xx codes) (default true)
  -async
        Run async requests
  -delay int
        Delay between requests in Millisecond
  -max-depth int
        MaxDepth limits the recursion depth of visited URLs.
  -noindex-rule string
        Comma-separated list of parameters as a string (default "noindex,nofollow")
  -outfile string
        Out sitemap file (default "./sitemap.xml")
  -parallel int
        Parallelism is the number of the maximum allowed concurrent requests
  -skip-noindex
        Do not add link to sitemap if it contains: 'meta name="robots" content="noindex,nofollow"' (default true)
  -timeout int
        Response timeout in second (default 10)
  -uri-filter string
        Filtering on uri prefix, example: /ru-ru , allowed use regex.
  -url string
        URL for scraping (default "http://go-colly.org/docs")
  -user-agent string
        User Agent (default "Sicra crawler, https://github.com/jidckii/sicra")
  -v    Verbose visiting URL
```

Example:  
```
build/linux/sicra -url http://go-colly.org -uri-filter="/docs" -timeout=60 -parallel=1 -delay=1000 -v=true 
```
