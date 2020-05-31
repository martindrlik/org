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
	Add("g", []byte(`{"h":"i"}`))
	Add("g", []byte(`{"j":"k"}`))
	os.Exit(m.Run())
}

func TestQuery(t *testing.T) {
	type Result struct {
		Body []byte
		Name string
	}
	tt := []struct {
		Query    string
		Handdler http.HandlerFunc
		Expected []Result
	}{
		{
			"all", all, []Result{
				{Body: []byte(`{"b":"c"}`), Name: "a"},
				{Body: []byte(`{"e":"f"}`), Name: "d"},
				{Body: []byte(`{"h":"i"}`), Name: "g"},
				{Body: []byte(`{"j":"k"}`), Name: "g"},
			},
		},
		{
			"q?device=d", q, []Result{
				{Body: []byte(`{"e":"f"}`), Name: "d"},
			},
		},
		{
			"q?device=g", q, []Result{
				{Body: []byte(`{"h":"i"}`), Name: "g"},
				{Body: []byte(`{"j":"k"}`), Name: "g"},
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.Query, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, "/"+tc.Query, nil)
			if err != nil {
				t.Fatal(err)
			}
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(tc.Handdler)
			handler.ServeHTTP(rr, req)
			if status := rr.Code; status != http.StatusOK {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, http.StatusOK)
			}
			dec := json.NewDecoder(rr.Body)
			for _, result := range tc.Expected {
				actualResult := &responseEntry{}
				if err := dec.Decode(actualResult); err != nil {
					t.Fatalf("error while decoding response: %v", err)
				}
				expect := string(result.Body)
				actualByte, err := base64.StdEncoding.DecodeString(actualResult.Body)
				if err != nil {
					t.Fatalf("error while decoding base64 body: %v", err)
				}
				actual := string(actualByte)
				if expect != actual {
					t.Errorf("expected body to be %q, got %q", expect, actual)
				}
				if result.Name != actualResult.Name {
					t.Errorf("expected name to be %q, got %q", result.Name, actualResult.Name)
				}
			}
		})
	}
}
