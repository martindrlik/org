package confirm

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestConfirmDelivery(t *testing.T) {
	mux := http.NewServeMux()
	confirmEndpoint := "/api/publish/MobileBackend/ConfirmNotificationDeliveryV2"
	called := 0
	actualPayload := struct {
		ApplicationID uint64
		Platfrom      string
		Token         string
	}{}
	mux.HandleFunc(confirmEndpoint, func(w http.ResponseWriter, r *http.Request) {
		called++
		defer r.Body.Close()
		dec := json.NewDecoder(r.Body)
		err := dec.Decode(&actualPayload)
		if err != nil {
			t.Fatal(err)
		}
		fmt.Fprintln(w, `{"error": "x"}`)
	})
	ts := httptest.NewServer(mux)
	defer ts.Close()
	confirmDelivery(&Payload{
		ApplicationID: "1",
		BaseURL:       ts.URL,
		Platform:      "iOS",
		Token:         "token",
	})
	if called != 1 {
		t.Errorf("expected confirmDelivery to be called once; called %v times", called)
	}
}
