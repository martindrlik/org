package fakegcm

import (
	"encoding/json"
	"io"
)

type Response struct {
	Failure     int
	MulticastId int `json:"multicast_id"`
	Results     []Result
	Success     int
}

type Result struct {
	Error          string
	MessageId      string `json:"message_id"`
	RegistrationId string `json:"registration_id"`
}

type TopicResponse struct {
	Error     string
	MessageId string `json:"message_id"`
}

func respond(w io.Writer, sr *SendRequest) error {
	n := len(sr.Message.RegistrationIds)
	res := Response{
		MulticastId: 2371663165171299815,
		Results:     make([]Result, n),
	}
	for i := 0; i < n; i++ {
		res.Results[i] = Result{
			MessageId: "0:5219441976194715812%8eda0b1da6bda",
			Error:     sr.Message.Data.ResponseError,
		}
	}
	if sr.Message.Data.ResponseError == "" {
		res.Success = n
	} else {
		res.Failure = n
	}
	enc := json.NewEncoder(w)
	return enc.Encode(res)
}
