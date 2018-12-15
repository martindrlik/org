package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"runtime/trace"
	"time"
)

var (
	url        = flag.String("url", "", "")
	concurrent = flag.Int("concurrent", 10, "")
	requests   = flag.Int("requests", 10, "")
)

func main() {
	f, err := os.Create("trace.out")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Fatal(err)
		}
	}()
	if err := trace.Start(f); err != nil {
		log.Fatal(err)
	}
	defer trace.Stop()
	flag.Parse()
	do(5) // warmup
	ch := make(chan int)
	for i := 0; i < *concurrent; i++ {
		go func() {
			ch <- do(*requests)
		}()
	}
	numOK := 0
	start := time.Now()
	for i := 0; i < *concurrent; i++ {
		numOK += <-ch
	}
	d := time.Now().Sub(start)
	total := (*concurrent) * (*requests)
	fmt.Printf("%d requests, %d OK, concurrent %d, %v => %v\n",
		total, numOK, *concurrent, d, float64(total)/d.Seconds())
}

func do(times int) int {
	numOK := 0
	for i := 0; i < times; i++ {
		res, err := http.Get(*url)
		if err != nil {
			log.Fatal(err)
		}
		defer res.Body.Close()
		if _, err := ioutil.ReadAll(res.Body); err != nil {
			log.Fatal(err)
		}
		if res.StatusCode == http.StatusOK {
			numOK++
		}
	}
	return numOK
}
