package fakegcm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/martindrlik/org/confirm"
)

func TestResponseCodeAndError(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/fcm/send", bytes.NewBuffer([]byte(`{
		"message": {
			"data": {
				"ResponseCode": "302",
				"ResponseError": "Error"
			},
			"registration_ids": ["x"]
		}
}`)))
	config := Configuration{
		confirmAdd: func(confirm.Payload) {},
		queryAdd:   func(string, []byte) {},
		Println:    func(...interface{}) {},
	}
	config.handle(w, r)
	rr := w.Result()
	if rr.StatusCode != http.StatusFound {
		t.Errorf("expected response code to be %v, got %v", http.StatusFound, rr.StatusCode)
	}
	var br Response
	dec := json.NewDecoder(rr.Body)
	if err := dec.Decode(&br); err != nil {
		t.Fatal(err)
	}
	if br.Results[0].Error != "Error" {
		t.Errorf("expected results[0].Error to be \"Error\", got %q", br.Results[0].Error)
	}
}

func TestQueryAdd(t *testing.T) {
	w := httptest.NewRecorder()
	data := []byte(`{"message":{"data": {},"registration_ids": ["x", "y"]}}`)
	r := httptest.NewRequest(http.MethodPost, "/fcm/send", bytes.NewBuffer(data))
	nqueryadd := 0
	config := Configuration{
		queryAdd: func(name string, actualData []byte) {
			if name != "x,y" {
				t.Errorf("expected name for query to be \"x,y\", got %q", name)
			}
			if !bytes.Equal(data, actualData) {
				t.Error("actual data for query does not match expected data:")
				t.Error("expected", data)
				t.Error("actual  ", actualData)
			}
			nqueryadd++
		},
		Println: func(...interface{}) {},
	}
	config.handle(w, r)
	if nqueryadd != 1 {
		t.Errorf("expected queryAdd to be called once, but called %v times", nqueryadd)
	}
}

func TestConfirmAdd(t *testing.T) {
	w := httptest.NewRecorder()
	data := []byte(`{"message":{"data": {},"registration_ids": ["x", "y"]}}`)
	r := httptest.NewRequest(http.MethodPost, "/fcm/send", bytes.NewBuffer(data))
	nconfirmadd := 0
	config := Configuration{
		ConfirmDelivery: true,
		confirmAdd: func(confirm.Payload) {
			nconfirmadd++
		},
		queryAdd: func(string, []byte) {},
		Println:  func(...interface{}) {},
	}
	config.handle(w, r)
	if nconfirmadd != 1 {
		t.Errorf("expected confirmAdd to be called once, but called %v times", nconfirmadd)
	}
}

func TestMessageOnly(t *testing.T) {
	w := httptest.NewRecorder()
	data := []byte(`{"message":{"data": {"message": "Hello"},"registration_ids": ["x", "y"]}}`)
	r := httptest.NewRequest(http.MethodPost, "/fcm/send", bytes.NewBuffer(data))
	nprintln := 0
	config := Configuration{
		MessageOnly: true,
		confirmAdd:  func(confirm.Payload) {},
		queryAdd:    func(string, []byte) {},
		Println: func(v ...interface{}) {
			if actual := fmt.Sprint(v...); actual != "Hello" {
				t.Errorf("expected message to be \"Hello\", got %q", actual)
			}
			nprintln++
		},
	}
	config.handle(w, r)
	if nprintln != 1 {
		t.Errorf("expected Println to be called once, but called %v times", nprintln)
	}
}

func TestInvalidJson(t *testing.T) {
	expectedMessages := []string{
		"POST /fcm/send HTTP/1.1\r\nHost: example.com\r\nUser-Agent: Go-http-client/1.1\r\nContent-Length: 8\r\n\r\n{\"data\":",
		"unexpected EOF",
	}
	w := httptest.NewRecorder()
	data := []byte(`{"data":`)
	r := httptest.NewRequest(http.MethodPost, "/fcm/send", bytes.NewBuffer(data))
	nprintln := 0
	config := Configuration{
		confirmAdd: func(confirm.Payload) {},
		queryAdd:   func(string, []byte) {},
		Println: func(v ...interface{}) {
			if actual := fmt.Sprint(v...); actual != expectedMessages[nprintln] {
				t.Errorf("expected message to be\n%q\ngot\n%q", expectedMessages[nprintln], actual)
			}
			nprintln++
		},
	}
	config.handle(w, r)
	if nprintln != 2 {
		t.Errorf("expected Println to be called twice, but called %v times", nprintln)
	}
}
