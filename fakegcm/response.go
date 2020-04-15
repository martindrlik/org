package fakegcm

import (
	"encoding/json"
	"io"
)

type Response struct {
	MulticastId  uint64   `json:"multicast_id"`
	Success      int      `json:"success"`
	Failure      int      `json:"failure"`
	CanonicalIds int      `json:"canonical_ids"`
	Results      []Result `json:"results"`
}

type Result struct {
	MessageId string `json:"message_id"`
	Error     string `json:"error"`
}

func respond(w io.Writer, notification *Notification) error {
	n := len(notification.Regs)
	res := Response{
		MulticastId: 2371663165171299815,
		Results:     make([]Result, n),
	}
	for i := 0; i < n; i++ {
		res.Results[i] = Result{
			MessageId: "0:5219441976194715812%8eda0b1da6bda",
			Error:     notification.Data.ResponseError,
		}
	}
	if notification.Data.ResponseError == "" {
		res.Success = n
	} else {
		res.Failure = n
	}
	enc := json.NewEncoder(w)
	return enc.Encode(res)
}
