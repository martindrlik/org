package notquery

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	Add("a", []byte(`{"b":"c"}`))
	Add("d", []byte(`{"e":"f"}`))
	os.Exit(m.Run())
}

func TestAll(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/all", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(all)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	dec := json.NewDecoder(rr.Body)
	expected := []struct {
		Body []byte
		Name string
	}{
		{Body: []byte(`{"b":"c"}`), Name: "a"},
		{Body: []byte(`{"e":"f"}`), Name: "d"},
	}
	for _, expected := range expected {
		actual := &responseEntry{}
		if err := dec.Decode(actual); err != nil {
			t.Fatalf("error while decoding response: %v", err)
		}
		if expected := base64.StdEncoding.EncodeToString(expected.Body); expected != actual.Body {
			t.Errorf("expected body to be base64 string %v, got %v", expected, actual.Body)
		}
		if expected.Name != actual.Name {
			t.Errorf("expected name to be a, got %v", actual.Name)
		}
	}
}

func TestQDevice(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/q?device=d", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(q)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	dec := json.NewDecoder(rr.Body)
	expected := []struct {
		Body []byte
		Name string
	}{
		{Body: []byte(`{"e":"f"}`), Name: "d"},
	}
	for _, expected := range expected {
		actual := &responseEntry{}
		if err := dec.Decode(actual); err != nil {
			t.Fatalf("error while decoding response: %v", err)
		}
		if expected := base64.StdEncoding.EncodeToString(expected.Body); expected != actual.Body {
			t.Errorf("expected body to be base64 string %v, got %v", expected, actual.Body)
		}
		if expected.Name != actual.Name {
			t.Errorf("expected name to be %v, got %v", expected.Name, actual.Name)
		}
	}
}
