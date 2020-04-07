package fakegcm

import (
	"bytes"
	"log"
	"net/http"

	"github.com/martindrlik/org/confirm"
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
	confirm.LogError = log.Println
	confirmDelivery = configuration.ConfirmDelivery
	messageOnly = configuration.MessageOnly
	http.HandleFunc("/", handle)
	return http.ListenAndServeTLS(
		configuration.Addr,
		configuration.CertFile,
		configuration.KeyFile,
		nil)
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
	if confirmDelivery {
		confirm.Channel <- confirm.Payload{
			ApplicationId: notification.Data.Content.ApplicationId,
			BaseURL:       notification.Data.Content.BaseURL,
			Platform:      "Android",
			Token:         notification.Data.NotificationToken,
		}
	}
	err = respond(w, len(notification.Regs))
	if err != nil {
		log.Println(err)
	}
}
