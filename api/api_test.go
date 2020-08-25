package api_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
	"testing"
	"ticket"
	"ticket/api"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/phayes/freeport"
)

func startTestServer() (addr string, s *ticket.Store) {
	s = ticket.NewStore()

	port, err := freeport.GetFreePort()
	log.Println(err)

	addr = net.JoinHostPort("localhost", strconv.Itoa(port))

	go api.ListenAndServe(addr, s)
	url := "http://" + addr + "/get/"
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	for resp.StatusCode == http.StatusNotFound {
		log.Println("RETRYING...")
		log.Println(url)
		time.Sleep(10 * time.Millisecond)
		resp, _ = http.Get(url)

	}
	log.Println("SERVER READY")
	return addr, s
}
func TestGet(t *testing.T) {
	t.Parallel()
	addr, s := startTestServer()
	_, _ = s.AddTicket(ticket.Ticket{
		Subject: "This is ticket 1",
	})
	_, _ = s.AddTicket(ticket.Ticket{
		Subject: "This is ticket 2",
	})
	url := "http://" + addr + "/get/2"
	log.Println(url)
	resp, err := http.Get(url)
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
	t.Parallel()
	addr, _ := startTestServer()

	want := ticket.Ticket{
		Subject:     "I hope this gets created",
		Description: "My screen broke!",
	}
	data, err := json.Marshal(want)
	if err != nil {
		t.Error(err)
	}
	resp, err := http.Post("http://"+addr+"/create", "application/json", bytes.NewBuffer(data))
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
