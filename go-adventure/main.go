package main

import (
	"fmt"
	"net/http"

	"github.com/mbraunwarth/adventure/adventure"
	api "github.com/mbraunwarth/adventure/http"
)

func main() {
	adventure.Load()
	fmt.Printf("ListAdventures => %s\n", adventure.ListAdventureNames())

	var h api.Handler
	http.Handle("/", h)
	http.ListenAndServe(":8080", nil)
}
