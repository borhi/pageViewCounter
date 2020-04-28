package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func main() {
	port := flag.Int("port", 3000, "service port")
	flag.Parse()

	commandChan := make(chan chan uint64)
	var counter uint64 = 0
	go func() {
		for replyChan := range commandChan {
			counter++
			replyChan <- counter
		}
	}()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		replyChan := make(chan uint64)
		commandChan <- replyChan
		counter := <-replyChan

		if _, err := w.Write([]byte(strconv.FormatUint(counter, 10))); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Print(err.Error())
			return
		}
	})

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil))
}
