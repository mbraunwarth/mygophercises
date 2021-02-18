package main

import (
	"log"

	"github.com/mbraunwarth/sitemap/sitemap"
)

func main() {
	// get host from command line (use args instead of flags)
	host := "https://www.sitemaps.org/"

	// parse host/build sitemap (naming convention?)
	s, err := sitemap.Build(host)
	if err != nil {
		log.Fatalf("could not build sitemap from host %v: %v", host, err)
	}

	// write sitemap to xml
	s.ToXML("sitemap.xml")
}
