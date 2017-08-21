package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"io"
	"log"
	"net/http"
	"os"
)

var (
	logger = log.New(os.Stdout, "fakegcm: ", log.LstdFlags)
	addr   = flag.String("addr", ":8080", "")
	cert   = flag.String("cert", "cert.pem", "")
	key    = flag.String("key", "key.pem", "")
	monly  = flag.Bool("message-only", false, "print message property only")
)

func main() {
	flag.Parse()
	http.HandleFunc("/", handler)
	go func() {
		for buf := range ch {
			logger.Println(buf)
		}
	}()
	err := http.ListenAndServeTLS(*addr, *cert, *key, nil)
	if err != nil {
		logger.Fatalln(err)
		os.Exit(1)
	}
}

type RequestData struct {
	Message string
}

type Request struct {
	Data            RequestData
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

var ch = make(chan *bytes.Buffer, 1000)

func handler(w http.ResponseWriter, r *http.Request) {
	buf := new(bytes.Buffer)
	r.Write(buf)

	sp := bytes.SplitN(buf.Bytes(), []byte("\r\n\r\n"), 2)
	body := bytes.NewBuffer(sp[1])
	req, err := readRequest(body)
	if err != nil {
		logger.Println(err)
		ch <- buf
		return
	}
	if *monly {
		ch <- bytes.NewBufferString(req.Data.Message)
	} else {
		ch <- buf
	}
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
