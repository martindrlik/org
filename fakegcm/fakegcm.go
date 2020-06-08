package fakegcm

import (
	"bytes"
	"log"
	"net/http"
	"strings"

	"github.com/martindrlik/org/confirm"
	"github.com/martindrlik/org/notquery"
	"github.com/tomchl/logfilter"
)

type Configuration struct {
	Addr     string
	CertFile string
	KeyFile  string

	ConfirmDelivery bool
	MessageOnly     bool

	Println func(...interface{})

	confirmAdd func(confirm.Payload)
	queryAdd   func(string, []byte)
}

func ListenAndServeTLS(config Configuration) error {
	cp := make(chan confirm.Payload, 500_000)
	confirm.Init(50, cp)
	config.confirmAdd = func(p confirm.Payload) { cp <- p }
	config.queryAdd = notquery.Add
	if config.Println == nil {
		config.Println = log.Println
	}
	http.HandleFunc("/", config.handle)
	srv := &http.Server{Addr: config.Addr, ErrorLog: log.New(&logfilter.IgnoreHTTPWriter{}, "", 0)}
	return srv.ListenAndServeTLS(config.CertFile, config.KeyFile)
}

func (config Configuration) handle(w http.ResponseWriter, r *http.Request) {
	buf := new(bytes.Buffer)
	r.Write(buf)
	sp := bytes.SplitN(buf.Bytes(), []byte("\r\n\r\n"), 2)
	sr, err := decodeSendRequest(bytes.NewBuffer(sp[1]))
	if err != nil {
		config.queryAdd("", sp[1])
		config.Println(buf)
		config.Println(err)
		return
	}
	config.queryAdd(strings.Join(sr.RegistrationIds, ","), sp[1])
	if config.MessageOnly {
		config.Println(sr.Data.Message)
	} else {
		config.Println(buf)
	}
	if config.ConfirmDelivery &&
		isSuccessCode(sr.Data.ResponseCode) &&
		sr.Data.ResponseError == "" {
		config.confirmAdd(confirm.Payload{
			ApplicationID: sr.Data.Content.ApplicationID,
			BaseURL:       sr.Data.Content.BaseURL,
			Platform:      "Android",
			Token:         sr.Data.NotificationToken,
		})
	}
	if sr.Data.ResponseCode != 0 {
		w.WriteHeader(sr.Data.ResponseCode)
	}
	err = respond(w, sr)
	if err != nil {
		config.Println(err)
	}
}

func isSuccessCode(i int) bool {
	return i == 0 || i >= 200 && i <= 300
}
