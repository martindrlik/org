package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
)

var (
	addr = flag.String("addr", ":8085", "listen and serve on addr")
	cert = flag.String("cert", "cert.pem", "")
	key  = flag.String("key", "key.pem", "")
)

func main() {
	flag.Parse()
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServeTLS(
		*addr,
		*cert,
		*key,
		nil))
}

type Reply struct {
	Message    string
	RetryAfter string
	Status     int
}

func handler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	dec := json.NewDecoder(r.Body)
	b := struct {
		Reply Reply
	}{}
	err := dec.Decode(&b)
	if err != nil {
		err = writeError(w, err)
		return
	}
	if err == nil {
		err = writeReply(w, b.Reply)
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}
}

func writeError(w http.ResponseWriter, err error) error {
	w.WriteHeader(http.StatusBadRequest)
	enc := json.NewEncoder(w)
	return enc.Encode(struct {
		Error string
	}{
		Error: err.Error(),
	})
}

func writeReply(w http.ResponseWriter, r Reply) error {
	if r.Status != 0 {
		w.WriteHeader(r.Status)
	}
	if r.RetryAfter != "" {
		w.Header().Add("Retry-After", r.RetryAfter)
	}
	enc := json.NewEncoder(w)
	return enc.Encode(struct {
		Message string
	}{
		Message: r.Message,
	})
}
