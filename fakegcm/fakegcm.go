package main

import (
	"bytes"
	"flag"
	"log"
	"net/http"

	"github.com/martindrlik/org/confirm"
)

var (
	addr = flag.String("addr", ":8080", "")
	cert = flag.String("cert", "cert.pem", "")
	key  = flag.String("key", "key.pem", "")

	monly = flag.Bool("message-only", false, "log notification message only")
	confr = flag.Bool("confirm-delivery", false, "confirm notification delivery")
)

func main() {
	confirm.LogError = log.Println
	flag.Parse()
	http.HandleFunc("/", handle)
	log.Fatal(http.ListenAndServeTLS(
		*addr,
		*cert,
		*key,
		nil))
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
	if *monly {
		log.Println(notification.Data.Message)
	} else {
		log.Println(buf)
	}
	if *confr {
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
