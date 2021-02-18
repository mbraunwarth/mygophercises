package sitemap

const (
	ns09 = "http://www.sitemaps.org/schemas/sitemap/0.9"
)

// Sitemap structure.
type Sitemap struct {
	host      string
	namespace string
}

// Build returns a new empty Sitemap for the given host.
func Build(host string) Sitemap {
	return Sitemap{
		host:      host,
		namespace: ns09,
	}
}

//type Site??
