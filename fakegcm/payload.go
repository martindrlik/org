package fakegcm

import (
	"encoding/json"
	"io"
)

type Targets struct {
	Condition       string
	NotificationKey string   `json:"notification_key"`
	RegistrationIds []string `json:"registration_ids"`
	To              string
}

type Options struct {
	CollapseKey           string `json:"collapse_key"`
	ContentAvailable      string `json:"content_available"`
	DryRun                bool   `json:"dry_run"`
	MutableContent        bool   `json:"mutable_content"`
	Priority              string
	RestrictedPackageName string `json:"restricted_package_name"`
	TimeToLive            int    `json:"time_to_live"`
}

type Payload struct {
	Data         Data
	Notification AndroidNotification
}

type Data struct {
	Content
	Message           string
	NotificationToken string
}

type Content struct {
	ApplicationID string
	BaseURL       string
	ResponseCode  string
	ResponseError string
}

type AndroidNotification struct {
	Body         string
	BodyLocArgs  []string `json:"body_loc_args"`
	BodyLocKey   string   `json:"body_loc_key"`
	ChannelId    string   `json:"android_channel_id"`
	ClickAction  string   `json:"click_action"`
	Color        string
	Icon         string
	Sound        string
	Tag          string
	Title        string
	TitleLocArgs []string `json:"title_loc_args"`
	TitleLocKey  string   `json:"title_loc_key"`
}

type Message struct {
	Targets
	Options
	Payload
}

type SendRequest struct {
	Message Message `json:"message"`
}

func decodeSendRequest(r io.Reader) (*SendRequest, error) {
	sr := new(SendRequest)
	dec := json.NewDecoder(r)
	err := dec.Decode(sr)
	if err != nil && err != io.EOF {
		return nil, err
	}
	return sr, nil
}
