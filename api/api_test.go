package api_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"testing"
	"ticket"
	"ticket/api"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestGet(t *testing.T) {
	s := ticket.NewStore()
	_ = s.NewTicket("This is ticket 1")
	_ = s.NewTicket("This is ticket 2")
	go func() {
		err := api.ListenAndServe(s)
		if err != nil {
			t.Error(err)
		}
	}()
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

func TestCreate(t *testing.T) {
	s := ticket.NewStore()
	go func() {
		err := api.ListenAndServe(s)
		if err != nil {
			t.Error(err)
		}
	}()
	want := ticket.Ticket{
		Subject:     "I hope this gets created",
		Description: "My screen broke!",
	}
	data, err := json.Marshal(want)
	if err != nil {
		t.Error(err)
	}
	resp, err := http.Post("http://localhost:9090/create", "application/json", bytes.NewBuffer(data))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		log.Println(resp.StatusCode)
		t.Fatalf("wanted http response status %d, got: %d", http.StatusOK, resp.StatusCode)
	}
	got := ticket.Ticket{}
	err = json.NewDecoder(resp.Body).Decode(&got)
	if err != nil {
		t.Fatal(err)
	}
	if !cmp.Equal(want, got, cmpopts.IgnoreFields(ticket.Ticket{}, "ID", "Status")) {
		t.Error(cmp.Diff(want, got))
	}
}
