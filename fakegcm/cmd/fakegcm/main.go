package main

import (
	"flag"

	"github.com/martindrlik/org/fakegcm"
)

var (
	addr  = flag.String("addr", ":8080", "")
	cert  = flag.String("cert", "cert.pem", "")
	key   = flag.String("key", "key.pem", "")
	monly = flag.Bool("message-only", false, "log notification message only")
	confr = flag.Bool("confirm-delivery", false, "confirm notification delivery")
)

func main() {
	flag.Parse()
	fakegcm.ListenAndServeTLS(fakegcm.Configuration{
		Addr:            *addr,
		CertFile:        *cert,
		KeyFile:         *key,
		MessageOnly:     *monly,
		ConfirmDelivery: *confr,
	})
}
