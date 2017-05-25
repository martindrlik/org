package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

var addr = flag.String("addr", ":8080", "")

var ch = make(chan *bytes.Buffer, 1000)

func main() {
	flag.Parse()
	http.HandleFunc("/", handler)
	go func() {
		for b := range ch {
			fmt.Println(b)
		}
	}()
	log.Fatal(http.ListenAndServe(*addr, nil))
}

type Request struct {
	RegistrationIds []string `json:"registration_ids"`
}

func readRequest(r io.Reader) (*Request, error) {
	req := new(Request)
	dec := json.NewDecoder(r)
	err := dec.Decode(&req)
	if err != nil && err != io.EOF {
		return nil, err
	}
	return req, nil
}

type Response struct {
	MulticastId  uint64   `json:"multicast_id"`
	Success      int      `json:"success"`
	Failure      int      `json:"failure"`
	CanonicalIds int      `json:"canonical_ids"`
	Results      []Result `json:"results"`
}

type Result struct {
	MessageId string `json:"message_id"`
	Error     string `json:"error"`
}

func handler(w http.ResponseWriter, r *http.Request) {
	buf := new(bytes.Buffer)
	fmt.Fprintf(buf, "\nfakegcm: %s\n", time.Now().UTC().Format(time.StampMilli))
	r.Write(buf)

	sp := bytes.SplitN(buf.Bytes(), []byte("\r\n\r\n"), 2)
	body := bytes.NewBuffer(sp[1])
	req, err := readRequest(body)
	if err != nil {
		fmt.Fprintf(buf, "fakegcm: %v\n", err)
		ch <- buf
		return
	}
	ch <- buf
	n := len(req.RegistrationIds)
	res := Response{
		MulticastId: 2371663165171299815,
		Results:     make([]Result, n),
	}
	for i := 0; i < n; i++ {
		res.Results[i] = Result{MessageId: "0:5219441976194715812%8eda0b1da6bda"}
	}
	if n > 1 {
		res.Success = n - 1
		res.Failure = 1
		res.Results[n-1].Error = "InvalidRegistration"
	}
	enc := json.NewEncoder(w)
	enc.Encode(res)
}
