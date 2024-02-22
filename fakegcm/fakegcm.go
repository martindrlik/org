package fakegcm

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/martindrlik/org/confirm"
	"github.com/martindrlik/org/notquery"
	"github.com/tomchl/logfilter"
)

type Configuration struct {
	Addr     string
	CertFile string
	KeyFile  string

	ConfirmBuffer      int
	ConfirmWorkerCount int

	ConfirmDelivery bool
	MessageOnly     bool

	Println func(...interface{})

	confirmAdd func(confirm.Payload)
	queryAdd   func(string, []byte)
}

func ListenAndServeTLS(config Configuration) error {
	if config.ConfirmBuffer == 0 {
		config.ConfirmBuffer = 2_000_000
	}
	if config.ConfirmWorkerCount == 0 {
		config.ConfirmWorkerCount = 50
	}
	cp := make(chan confirm.Payload, config.ConfirmBuffer)
	confirm.Init(config.ConfirmWorkerCount, cp)
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
	config.queryAdd(strings.Join(sr.Message.RegistrationIds, ","), sp[1])
	if config.MessageOnly {
		config.Println(sr.Message.Data.Message)
	} else {
		config.Println(buf)
	}
	statusCode, err := strconv.Atoi(sr.Message.Data.ResponseCode)
	if sr.Message.Data.ResponseCode != "" && err != nil {
		fmt.Printf("unable to parse response code string: %s\n", sr.Message.Data.ResponseCode)
	}
	if config.ConfirmDelivery &&
		isSuccessCode(statusCode) &&
		sr.Message.Data.ResponseError == "" {
		config.confirmAdd(confirm.Payload{
			ApplicationID: sr.Message.Data.Content.ApplicationID,
			BaseURL:       sr.Message.Data.Content.BaseURL,
			Platform:      "Android",
			Token:         sr.Message.Data.NotificationToken,
		})
	}
	if statusCode != 0 {
		w.WriteHeader(statusCode)
	}
	err = respond(w, sr)
	if err != nil {
		config.Println(err)
	}
}

func isSuccessCode(i int) bool {
	return i == 0 || i >= 200 && i <= 300
}
