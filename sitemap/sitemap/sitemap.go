package sitemap

import (
	"log"
	"net/http"
	"regexp"

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

	// validate links and add to list for current site - starting with host
	s.M = parseLinks(ls, s.host, make(map[string][]string, 0))

	return nil
}

// parseLinks validates links and add to list for current site.
func parseLinks(ls []link.Link, s string, m map[string][]string) map[string][]string {
	// matches in-house links
	validLink := regexp.MustCompile(`^([a-zA-Z]+\/?)+\.[a-zA-Z]*`)
	// matches links prefixed with the domain
	validHostLink := regexp.MustCompile(s)

	var list []string

	for _, l := range ls {
		// link not already parsed?
		for k := range m {
			if k != l.Href {
				// is link in domain?
				loc := validHostLink.FindIndex([]byte(l.Href))
				if (loc != nil && loc[0] == 0) || validLink.MatchString(l.Href) {
					// if so add to sitemap list if not already added
					if !contains(list, l.Href) {
						list = append(list, l.Href)
					}
				}
			}
		}
	}

	m[s] = list

	return m
}

// ToXML writes the corresponding XML from the sitemap.
func (s Sitemap) ToXML(file string) error {
	log.Println(s.M)
	return nil
}

// contains checks if a string is present in a slice.
func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}
