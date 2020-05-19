package fakegcm

import (
	"bytes"
	"log"
	"net/http"

	"github.com/martindrlik/org/confirm"
	"github.com/tomchl/logfilter"
)

type Configuration struct {
	Addr     string
	CertFile string
	KeyFile  string

	ConfirmDelivery bool
	MessageOnly     bool
}

var (
	confirmDelivery bool
	messageOnly     bool
)

func ListenAndServeTLS(configuration Configuration) error {
	confirmDelivery = configuration.ConfirmDelivery
	messageOnly = configuration.MessageOnly
	http.HandleFunc("/", handle)

	server := &http.Server{Addr: configuration.Addr, ErrorLog: log.New(&logfilter.IgnoreHTTPWriter{}, "", 0)}
	return server.ListenAndServeTLS(configuration.CertFile, configuration.KeyFile)
}

func handle(w http.ResponseWriter, r *http.Request) {
	buf := new(bytes.Buffer)
	r.Write(buf)

	sp := bytes.SplitN(buf.Bytes(), []byte("\r\n\r\n"), 2)
	notification, err := requestNotification(bytes.NewBuffer(sp[1]))
	if err != nil {
		log.Println(buf)
		log.Println(err)
		return
	}
	if messageOnly {
		log.Println(notification.Data.Message)
	} else {
		log.Println(buf)
	}
	if confirmDelivery &&
		isSuccessCode(notification.Data.ResponseCode) &&
		notification.Data.ResponseError == "" {
		confirm.Channel <- confirm.Payload{
			ApplicationId: notification.Data.Content.ApplicationId,
			BaseURL:       notification.Data.Content.BaseURL,
			Platform:      "Android",
			Token:         notification.Data.NotificationToken,
		}
	}
	if notification.Data.ResponseCode != 0 {
		w.WriteHeader(notification.Data.ResponseCode)
	}
	err = respond(w, notification)
	if err != nil {
		log.Println(err)
	}
}

func isSuccessCode(i int) bool {
	return i == 0 || i >= 200 && i <= 300
}
