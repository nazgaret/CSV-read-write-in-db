package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	end bool
	url = "http://localhost:9000"
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

	time.Sleep(time.Second * 5)
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
	data, err := json.Marshal(person)
	if err != nil {
		fmt.Println(err)
		return
	}
	r := bytes.NewReader(data)
	_, err = http.Post(url, "application/json", r)
	if err != nil {
		fmt.Println(err)
		return
	}
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
