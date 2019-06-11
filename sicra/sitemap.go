package sicra

import (
	"log"
	"os"
	"time"
)

// GenerateSiteMap will write the crawled url to given file
func GenerateSiteMap(fileName string, urls []string) error {
	err := deleteFileIfExists(fileName)
	if err != nil {
		log.Fatal(err)
	}

	fh, err := os.OpenFile(fileName, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer fh.Close()
	currentTime := time.Now().Format(time.RFC3339)
	fh.WriteString("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n")
	fh.WriteString("<urlset xmlns=\"http://www.sitemaps.org/schemas/sitemap/0.9\">\n")
	for _, loc := range urls {
		fh.WriteString("    " + "<url>\n")
		fh.WriteString("      " + "<loc>" + loc + "</loc>\n")
		fh.WriteString("      " + "<lastmod>" + currentTime + "</lastmod>\n")
		fh.WriteString("      " + "<changefreq>hourly</changefreq>\n")
		fh.WriteString("      " + "<priority>0.5</priority>\n")
		fh.WriteString("    " + "</url>\n")
	}
	fh.WriteString("</urlset> ")

	return nil
}

// GenerateTxt generate txt file for error list url or skiped  noindex
func GenerateTxt(fileName string, urls []string) error {
	err := deleteFileIfExists(fileName)
	if err != nil {
		log.Fatal(err)
	}
	fh, err := os.OpenFile(fileName, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer fh.Close()
	for _, loc := range urls {
		fh.WriteString(loc + "\n")
	}

	return nil
}

//deleteFileIfExists deletes a file if exists
func deleteFileIfExists(fileName string) error {
	//delete old file first
	if _, err := os.Stat(fileName); err == nil {
		err := os.Remove(fileName)
		if err != nil {
			log.Fatal(err)
		}
		return nil
	}
	return nil
}
