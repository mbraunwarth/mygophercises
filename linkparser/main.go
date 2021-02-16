package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/mbraunwarth/link/link"
)

func main() {
	var (
		input string
		usage = "Usage: link -input <PATH_TO_HTML>"
	)

	flag.StringVar(&input, "input", "", "specify the html input file")
	flag.Parse()

	// if input file not specified exit with error
	if input == "" {
		fmt.Fprintf(os.Stderr, "please specify an HTML document to parse\n%v\n", usage)
	}

	f, err := os.Open(input)
	if err != nil {
		log.Fatalf("could not open file %v: %v", input, err)
	}
	defer f.Close()

	ls, err := link.Parse(f)
	fmt.Println(ls)
}
