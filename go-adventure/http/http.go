package http

import (
	"html/template"
	"log"
	"net/http"

	"github.com/mbraunwarth/adventure/adventure"
)

// Handler wraps the HTTP handler for adventures.
type Handler struct {
	s adventure.Service
}

// ServeHTTP function for Handler.
func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	n := "blue gopher"
	adv := adventure.RTFV(n)

	var arc adventure.Arc
	for _, a := range adv.Arcs {
		if a.ID == "new-york" {
			arc = a

		}
	}

	data := struct {
		Adventure string
		Arc       string
		Story     []string
		Options   []adventure.Option
	}{
		Adventure: adv.Name,
		Arc:       arc.Title,
		Story:     arc.Story,
		Options:   arc.Options,
	}

	t, err := template.New("home.html").ParseFiles("tmpl/home.html")
	if err != nil {
		log.Fatalf("error parsing templates: %s", err)
	}
	if err := t.Execute(w, data); err != nil {
		log.Fatalf("error executing template: %s", err)
	}

	//fmt.Fprintf(w, "Adventures available: %s", adventure.ListAdventureNames())
}
