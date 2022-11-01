package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
)

var quizFilename string

func main() {
	flag.StringVar(&quizFilename, "f", "problems.csv", "path for quiz file")
	flag.Parse()

	f, err := os.Open(quizFilename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	r := csv.NewReader(f)
	records, err := r.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	correct := 0
	quiz_size := len(records)

	for _, entry := range records {
		var answer string
		fmt.Println(entry[0])
		fmt.Scanln(&answer)

		if entry[1] == answer {
			correct += 1
		}
	}

	fmt.Printf("%d correct out of %d\n", correct, quiz_size)
}
