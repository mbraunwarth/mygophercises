package sitemap

import (
	"log"
	"net/http"

	"github.com/mbraunwarth/sitemap/link"
)

const (
	ns09 = "http://www.sitemaps.org/schemas/sitemap/0.9"
)

// TODO restructure Sitemap struct

// Sitemap structure.
type Sitemap struct {
	m         map[string]string
	host      string
	namespace string
}

// TODO Builder does not return any errors, think of 'readSite' first

// Build returns a new empty Sitemap for the given host.
func Build(host string) (Sitemap, error) {
	// TODO when program's working for further improvement fire go routines for this
	readSite(host)

	return Sitemap{
		host:      host,
		namespace: ns09,
	}, nil
}

// TODO new name! make it receiver method on Sitemap for access to host, map etc.
func readSite(host string) error {
	// initiate client requests for each link on the site
	resp, err := http.Get(host)
	if err != nil {
		return err
	}

	// remember closing the body
	b := resp.Body
	defer b.Close()

	// get all links from the host
	ls, err := link.Parse(b)
	if err != nil {
		return err
	}

	for _, l := range ls {
		log.Println(l.Href)
		parseLink(l.Href)
	}

	return nil
}

// TODO new function necessary? if so, use receiver method instead
func parseLink(link string) {
	// is link in domain?
	//   true if of form /some/link
	// otherwise match first part against host value
	//   true if matching performs with true

	// if true add to sitemap list
}

// ToXML writes the corresponding XML from the sitemap.
func (s Sitemap) ToXML(file string) error {
	return nil
}
