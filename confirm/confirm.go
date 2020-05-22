package confirm

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"path"
)

var Channel = make(chan Payload, 100)
var LogError = log.Print

func init() {
	go func() {
		for payload := range Channel {
			confirmDelivery(&payload)
		}
	}()
}

func confirmDelivery(payload *Payload) {
	u, err := payload.requestURL()
	if err != nil {
		logError(err)
		return
	}
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{
		InsecureSkipVerify: true,
	}
	r, err := http.Post(u.String(), "application/json", payload.requestBody())
	if err != nil {
		logError(err)
		return
	}
	defer r.Body.Close()
	err = responseBody(r)
	if err != nil {
		logError(err)
	}
}

func (payload *Payload) requestURL() (*url.URL, error) {
	u, err := url.Parse(payload.BaseURL)
	if err != nil {
		return nil, err
	}
	u.Path = path.Join(u.Path, "/api/publish/MobileBackend/"+
		"ConfirmNotificationDeliveryV2")
	return u, nil
}

func (payload *Payload) requestBody() io.Reader {
	b := &bytes.Buffer{}
	enc := json.NewEncoder(b)
	err := enc.Encode(struct {
		ApplicationID     uint64
		NotificationToken string
		Platform          string
	}{
		ApplicationID:     payload.ApplicationID,
		Platform:          payload.Platform,
		NotificationToken: payload.Token,
	})
	if err != nil {
		panic(fmt.Errorf("could not encode request body: %v", err))
	}
	return b
}

func responseBody(r *http.Response) error {
	if r.StatusCode != http.StatusOK {
		return fmt.Errorf("confirmDelivery response: "+
			"unexpected status: %v", r.Status)
	}
	dec := json.NewDecoder(r.Body)
	b := struct {
		Error string
	}{}
	if err := dec.Decode(&b); err != nil && err != io.EOF {
		return fmt.Errorf("confirmDelivery response: "+
			"could not decode json: %v", err)
	}
	if b.Error != "" {
		return fmt.Errorf("confirmDelivery response: error: %v", b.Error)
	}
	return nil
}

func logError(err error) {
	if LogError != nil {
		LogError(err)
	}
}
