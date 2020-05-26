package notquery

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/martindrlik/org/store"
)

var s = store.NewStore(1000)

func Add(name string, data []byte) { s.Add(name, data) }

func init() {
	http.HandleFunc("/all", all)
	http.HandleFunc("/q", q)
}

func all(w http.ResponseWriter, r *http.Request) {
	enc := json.NewEncoder(w)
	for _, u := range s.All() {
		b := base64.StdEncoding.EncodeToString(u.Data)
		if err := enc.Encode(responseEntry{b, u.Name, u.Time}); err != nil {
			log.Print(err)
			continue
		}
	}
}

func q(w http.ResponseWriter, r *http.Request) {
	vals, ok := s.ByName(r.URL.Query().Get("device"))
	if !ok {
		return
	}
	enc := json.NewEncoder(w)
	for _, u := range vals {
		b := base64.StdEncoding.EncodeToString(u.Data)
		if err := enc.Encode(responseEntry{b, u.Name, u.Time}); err != nil {
			log.Print(err)
			continue
		}
	}
}

type responseEntry struct {
	Body string
	Name string
	Time time.Time
}
