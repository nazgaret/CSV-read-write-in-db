package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	// setup reader
	csvIn, err := os.Open("./data/data.csv")
	if err != nil {
		log.Fatal(err)
	}
	r := csv.NewReader(csvIn)

	// handle header
	rec, err := r.Read()
	if err != nil {
		log.Fatal(err)
	}
	rec = append(rec, "score")

	callWrite(rec)

	for {
		rec, err = r.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal(err)
		}

		callWrite(rec)

	}
}
func callWrite(person []string) {
	fmt.Println(person)
}
