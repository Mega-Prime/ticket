package api_test

import (
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
	"ticket"
	"ticket/api"
)

func TestHandler(t *testing.T) {
	s := ticket.NewStore()
	_ = s.NewTicket("This is ticket 1")
	_ = s.NewTicket("This is ticket 2")
	go api.ListenAndServe(s)

	resp, err := http.Get("http://localhost:9090/get/2")
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
	if !strings.Contains(string(data), "This is ticket 2") {
		t.Errorf("string not found. Got: %q", data)
	}

}
