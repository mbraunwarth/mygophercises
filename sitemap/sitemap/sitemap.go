package sitemap

import (
	"log"
	"net/http"

	"github.com/mbraunwarth/sitemap/link"
)

const (
	ns09 = "http://www.sitemaps.org/schemas/sitemap/0.9"
)

// Sitemap structure.
type Sitemap struct {
	// actual map with entries e.g. '/home' -> ['/login', '/about']
	M map[string][]string

	// host is of form hostname.domain.[org|com|de|...]
	host string

	// XML schema of the sitemap protocol
	namespace string
}

// Build returns a new empty Sitemap for the given host.
func Build(host string) (Sitemap, error) {
	s := Sitemap{
		M:         make(map[string][]string),
		host:      host,
		namespace: ns09,
	}

	// TODO when program's working for further improvement fire go routines for this
	if err := s.parseHost(); err != nil {
		return Sitemap{}, err
	}

	return s, nil
}

// parseHost makes the inital request to the sitemaps host, eventually receiving
// a response which will be read in order to process all links in its body, which
// themselves will be added to the sitemaps map field.
func (s Sitemap) parseHost() error {
	// initiate client requests for each link on the site
	resp, err := http.Get(s.host)
	if err != nil {
		return err
	}

	// remember closing the body
	defer resp.Body.Close()

	// get all links from the host
	ls, err := link.Parse(resp.Body)
	if err != nil {
		return err
	}

	for _, l := range ls {
		log.Println(l.Href)
		// is link in domain?
		//   true if of form /some/link
		// otherwise match first part against host value
		//   true if matching performs with true

		// if true add to sitemap list
	}

	return nil
}

// ToXML writes the corresponding XML from the sitemap.
func (s Sitemap) ToXML(file string) error {
	return nil
}
