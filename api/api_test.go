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
	_ = s.NewTicket("New test ticket")

	go api.ListenAndServe(s)

	resp, err := http.Get("http://localhost:9090/1") //make this work
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
	if !strings.Contains(string(data), "New test ticket") {
		t.Errorf("string not found. Got: %q", data)
	}

}


