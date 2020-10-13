package api_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"
	"testing"
	"ticket"
	"ticket/api"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/phayes/freeport"
)

const (
	testDB         = "dbStore"
	testCollection = "ticket_test"
	dbURI          = "mongodb+srv://megaPrime:H*rdlifease22@prime1.9fk1d.mongodb.net/quickstart?retryWrites=true&w=majority"
)

func startTestServer(t *testing.T) (addr string) {

	port, err := freeport.GetFreePort()
	if err != nil {
		t.Fatal(err)
	}

	addr = net.JoinHostPort("localhost", strconv.Itoa(port))

	go func() {
		err := api.ListenAndServe(addr, dbURI, testDB, testCollection)
		if err != nil {
			t.Fatal(err)
		}

	}()

	url := "http://" + addr + "/healthz"

	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	for resp.StatusCode != http.StatusOK {
		log.Println("RETRYING...")
		time.Sleep(10 * time.Millisecond)
		resp, _ = http.Get(url)

	}

	return addr
}

func TestCreateAndGet(t *testing.T) {
	t.Parallel()
	addr := startTestServer(t)

	want := ticket.Ticket{
		Subject:     "I hope this gets created",
		Description: "My screen broke!",
	}
	data, err := json.Marshal(want)
	if err != nil {
		t.Error(err)
	}

	fmt.Printf("This is the data: %q", data)

	fmt.Println("Got here")
	resp, err := http.Post("http://"+addr+"/create", "application/json", bytes.NewBuffer(data))

	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {

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

	url := "http://" + addr + "/get/" + string(got.ID)

	resp, err = http.Get(url)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("wanted http response status %d, got: %d", http.StatusOK, resp.StatusCode)
	}

	err = json.NewDecoder(resp.Body).Decode(&got)
	if err != nil {
		t.Fatal(err)
	}
	if !cmp.Equal(want, got, cmpopts.IgnoreFields(ticket.Ticket{}, "ID", "Status")) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestGetAllTickets(t *testing.T) {
	t.Parallel()
	addr := startTestServer(t)

	want := []ticket.Ticket{
		{
			Subject: "Test Get all tickets #1",
		},
		{
			Subject: "Test Get all tickets #2",
		},
	}

	data, err := json.Marshal(want)
	if err != nil {
		t.Error(err)
	}

	fmt.Printf("this is  the data: %q\n", data)

	fmt.Println("Are we after ticket creation?")
	resp, err := http.Post("http://"+addr+"/create", "application/json", bytes.NewBuffer(data))
	if err != nil {
		t.Fatal(err)
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("Needed http response status %d, got: %d", http.StatusOK, resp.StatusCode)
	}

	var got []ticket.Ticket
	err = json.NewDecoder(resp.Body).Decode(&got)
	if err != nil {
		t.Fatal(err)
	}

	if !cmp.Equal(want, got, cmpopts.IgnoreFields(ticket.Ticket{}, "ID", "Status")) {
		t.Error(cmp.Diff(want, got))
	}

	url := "http://" + addr + "/all"
	fmt.Println("this is the url:", url)
	resp, err = http.Get(url)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("wanted http response status %d, got: %d", http.StatusOK, resp.StatusCode)
	}

	err = json.NewDecoder(resp.Body).Decode(&got)
	if err != nil {
		t.Fatal(err)
	}

	if !cmp.Equal(want, got, cmpopts.IgnoreFields(ticket.Ticket{}, "ID", "Status")) {
		t.Error(cmp.Diff(want, got))
	}

}
