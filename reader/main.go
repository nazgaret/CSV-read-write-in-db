package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	end bool
)

func main() {
	//graceful stop of the app
	var gracefulStop = make(chan os.Signal)
	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)
	go func() {
		sig := <-gracefulStop
		fmt.Printf("caught sig: %+v", sig)
		fmt.Println("Wait to finish processing")
		for range time.Tick(time.Second) {
			if end {
				os.Exit(0)
			}
		}
	}()

	// setup reader
	csvIn, err := os.Open("./data/data.csv")
	if err != nil {
		log.Fatal(err)
	}
	chWrite := processCSV(csvIn)
	for {
		if len(chWrite) > 0 {
			person := <-chWrite
			callWrite(person)
		}
	}
}

//callWrite send data to writer
func callWrite(person []string) {
	//send to writer
	fmt.Println(person)
}

//processCSV read csv line by line
func processCSV(rc io.Reader) (ch chan []string) {
	ch = make(chan []string, 10)
	go func() {
		r := csv.NewReader(rc)
		if _, err := r.Read(); err != nil { //read header
			log.Fatal(err)
		}
		defer close(ch)
		for {
			rec, err := r.Read()
			if err != nil {
				if err == io.EOF {
					end = true
					break
				}
				log.Fatal(err)
			}
			ch <- rec
		}
	}()
	return
}
