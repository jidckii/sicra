# sicra
Simple crawler and sitemap generator on golang

build:  
```
go get -u github.com/jidckii/sicra/...
make
```
Run:
```
build/sicra -h

Usage of build/sicra:
  -async
        Run async requests
  -delay int
        Delay between requests in Millisecond
  -max-depth int
        MaxDepth limits the recursion depth of visited URLs.
  -outfile string
        Out sitemap file (default "./sitemap.xml")
  -parallel int
        Parallelism is the number of the maximum allowed concurrent requests
  -skip-noindex
        Do not add link to sitemap if it contains: 'meta name = "googlebot" content = "noindex"' (default true)
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
build/linux/amd64/sicra -url http://go-colly.org -uri-filter="/docs" -timeout=60 -parallel=1 -delay=1000 -v=true 
```
