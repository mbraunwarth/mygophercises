package sitemap

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	ns09 = "http://www.sitemaps.org/schemas/sitemap/0.9"
)

type m map[string]string

// Sitemap structure.
type Sitemap struct {
	m
	host      string
	namespace string
}

// Build returns a new empty Sitemap for the given host.
func Build(host string) (Sitemap, error) {
	readSite(host)

	return Sitemap{
		host:      host,
		namespace: ns09,
	}, nil
}

func readSite(host string) error {
	// initiate client requests for each link on the site
	resp, err := http.Get(host)
	if err != nil {
		return err
	}

	b, err := ioutil.ReadAll(resp.Body)
	fmt.Println(string(b))
	if err != nil {
		return err
	}

	// get all links from the host
	// ls, err := link.Parse(r)
	// if err != nil {
	// return link.Link{}, err
	// }
	return nil
}

// ToXML writes the corresponding XML from the sitemap.
func (s Sitemap) ToXML(file string) error {
	return nil
}
