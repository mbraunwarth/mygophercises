package link

import (
	"fmt"
	"os"
	"strings"

	"golang.org/x/net/html"
)

// Link structure assembles the href attribute and the underlying text of a link tag.
type Link struct {
	Href string
	Text string
}

func (l Link) String() string {
	return fmt.Sprintf("{%v - %v}", l.Href, l.Text)
}

// Parse takes an HTML file and parses its content such that it returns a list of all
// link tags with the corresponding href attribute and its text.
func Parse(file *os.File) ([]Link, error) {
	// as the parser expects an io.Reader, just passing the os.File
	root, err := html.Parse(file)
	if err != nil {
		return nil, err
	}

	var links []Link
	var traverse func(n *html.Node)
	traverse = func(n *html.Node) {
		// get href attribute if node is a link tag
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key == "href" {
					t := stripText(n)
					links = append(links, Link{Href: a.Val, Text: t})
				}
			}
		}

		// depth first traversel
		if n.FirstChild != nil {
			traverse(n.FirstChild)
		}
		if n.NextSibling != nil {
			traverse(n.NextSibling)
		}
	}
	traverse(root)

	return links, nil
}

func stripText(n *html.Node) string {
	var txt []string
	if n.Type == html.TextNode {
		t := n.Data
		txt = append(txt, t)
	}
	// depth first traversel
	if n.FirstChild != nil {
		txt = append(txt, stripText(n.FirstChild))
	}
	if n.NextSibling != nil {
		txt = append(txt, stripText(n.NextSibling))
	}

	return strings.Join(txt, "")
}
