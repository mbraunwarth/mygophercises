package main

import (
	"fmt"
	"log"
	"os"

	"github.com/mbraunwarth/sitemap/link"
)

func main() {
	f, err := os.Open("ex.html")
	if err != nil {
		log.Fatalf("could not open file %v: %v", f.Name(), err)
	}

	ls, err := link.Parse(f)
	if err != nil {
		log.Fatalf("could not parse file %v: %v", f.Name(), err)
	}

	fmt.Println(ls)
}
