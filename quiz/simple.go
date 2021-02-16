package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
)

var problemsPath string

// user can customize the problems file via cli flag
// if not set by the user, it defaults to problems.csv
// a pair of a question and its corresponding answer is furthermore referred to as a problem
func main() {
	// parse flags and determine the problems file
	flag.StringVar(&problemsPath, "problems", "problems.csv", "CSV file with problems and their answeres separated by commas")
	flag.Parse()

	// preparing csv file for reading
	problemsFile, err := os.Open(problemsPath)
	if err != nil {
		log.Fatalf("error opening file: %s", err)
	}
	defer problemsFile.Close()

	// parse csv content to records variable resulting in a 2D slice
	r := csv.NewReader(problemsFile)
	records, err := r.ReadAll()
	if err != nil {
		log.Fatalf("error reading csv: %s", err)
	}

	// fill problems map resulting in the form Question -> Answer [String -> String]
	problems := make(map[string]string, len(records))
	for _, row := range records {
		problems[row[0]] = row[1]
	}

	// for each question store if the user answered correct
	correct := 0

	// ask question(s) and take answer from input (single word/number)
	// no response to user till all questions has been answered,
	// but check for correctness and store either right or wrong answered
	for question, answer := range problems {
		// resetting user input to empty string
		input := ""

		fmt.Printf("%s = ", question)
		if _, err := fmt.Scanln(&input); err != nil {
			log.Fatalf("error scanning user input: %s", err)
		}

		if input == answer {
			correct++
		}
	}

	// output total number of questions and those which were answered correctly
	fmt.Printf("%d answered questions\n", len(problems))
	fmt.Printf("%d questions answered correct\n", correct)
}
