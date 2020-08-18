package api_test

import (
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
	"ticket/api"
)

func TestHandler(t *testing.T) {
	go api.ListenAndServe()

	resp, err := http.Get("http://localhost:9090")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("wanted http response status %d, got: %d", http.StatusOK, resp.StatusCode)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), "this is open") {
		t.Errorf("string not found. Got: %q", data)
	}

}
