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

func respond(w io.Writer, n int) error {
	res := Response{
		MulticastId: 2371663165171299815,
		Results:     make([]Result, n),
	}
	for i := 0; i < n; i++ {
		res.Results[i] = Result{MessageId: "0:5219441976194715812%8eda0b1da6bda"}
	}
	if n > 1 {
		res.Success = n - 1
		res.Failure = 1
		res.Results[n-1].Error = "InvalidRegistration"
	}
	enc := json.NewEncoder(w)
	return enc.Encode(res)
}
