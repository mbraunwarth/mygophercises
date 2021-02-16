package link

import (
	"log"
	"os"

	"golang.org/x/net/html"
)

// Link structure assembles the href attribute and the underlying text of a link tag.
type Link struct {
	Href string
	Text string
}

// Parse takes an HTML file and parses its content such that it returns a list of all
// link tags with the corresponding href attribute and its text.
func Parse(file *os.File) ([]Link, error) {
	// as the parser expects an io.Reader, just passing the os.File
	doc, err := html.Parse(file)
	if err != nil {
		return nil, err
	}
	log.Println(doc)
	return nil, nil
}
