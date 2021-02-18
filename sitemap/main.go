package main

import (
	"github.com/mbraunwarth/sitemap/sitemap"
)

func main() {
	// get host from command line (use args instead of flags)
	host := "https://www.sitemaps.org/"

	// parse host/build sitemap (naming convention?)
	s := sitemap.Build(host)

	// write sitemap to xml
	s.ToXML("sitemap.xml")
}
