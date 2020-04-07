package fakegcm

import (
	"encoding/json"
	"io"
)

type Notification struct {
	Data Payload
	Regs []string `json:"registration_ids"`
}

type Payload struct {
	Content
	Message string

	NotificationToken string
}

type Content struct {
	ApplicationId uint64
	BaseURL       string
}

func requestNotification(r io.Reader) (*Notification, error) {
	n := &Notification{}
	dec := json.NewDecoder(r)
	err := dec.Decode(n)
	if err != nil && err != io.EOF {
		return nil, err
	}
	return n, nil
}
