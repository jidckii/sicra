# sicra
Simple crawler and sitemap generator in golang

build:  
```
go get -u github.com/gocolly/colly/...

go get github.com/jidckii/sicra/...
cd $GOPATH/src/github.com/jidckii/sicra
go build -o build/sicra main.go
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
build/sicra -url http://go-colly.org/ -timeout=60 -parallel=1 -delay=60 -v=true
```
