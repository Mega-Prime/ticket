package api_test

import (
	"bytes"
	"context"
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
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	testDB         = "dbStore"
	testCollection = "ticket_test"
	// TODO get this from environment
	dbURI          = "mongodb+srv://megaPrime:H*rdlifease22@prime1.9fk1d.mongodb.net/quickstart?retryWrites=true&w=majority"
)

func startTestServer(t *testing.T) (addr string) {
	t.Helper()
	port, err := freeport.GetFreePort()
	if err != nil {
		t.Fatal(err)
	}

	addr = net.JoinHostPort("localhost", strconv.Itoa(port))

	go func() {
		err := api.ListenAndServe(addr, dbURI, testDB, testCollection)
		fmt.Println("Listening")
		if err != nil {
			t.Fatal(err)
		}

	}()

	url := "http://" + addr + "/healthz"
	// Not too happy with this
	time.Sleep(800 * time.Millisecond)
	_, err = http.Get(url)
	for err != nil {
		log.Println("Waiting for connection...")
		time.Sleep(30 * time.Millisecond)
		_, err = http.Get(url)
	}
	return addr
}

func TestCreateAndGet(t *testing.T) {
	t.Parallel()
	addr := startTestServer(t)
	t.Cleanup(func() {
		cleanUpTestCollection(t, dbURI)
	})
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
	// Can't run in parallel because we need to know the whole contents
	// of the ticket store.
	addr := startTestServer(t)
	t.Cleanup(func() {
		cleanUpTestCollection(t, dbURI)
	})
	want := []ticket.Ticket{
		{
			Subject: "Test Get all tickets #1",
		},
		{
			Subject: "Test Get all tickets #2",
		},
	}
	buf := bytes.Buffer{}
	err := want[0].ToJSON(&buf)
	if err != nil {
		t.Fatal(err)
	}
	resp, err := http.Post("http://"+addr+"/create", "application/json", &buf)
	if err != nil {
		t.Fatal(err)
	}
	err = want[1].ToJSON(&buf)
	if err != nil {
		t.Fatal(err)
	}
	resp, err = http.Post("http://"+addr+"/create", "application/json", &buf)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	var got []ticket.Ticket
	url := "http://" + addr + "/all"
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

func cleanUpTestCollection(t *testing.T, dbURI string) {
	t.Helper()
	client, err := mongo.NewClient(options.Client().ApplyURI(dbURI))

	if err != nil {
		t.Fatalf("tried to access mongo URI %q, got %q", dbURI, err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		t.Fatalf("tried to access mongo URI %q, got %q", dbURI, err)
	}
	defer client.Disconnect(ctx)

	// Don't care if it fails. Omit error.
	err = client.Database(testDB).Collection(testCollection).Drop(ctx)
	if err != nil {
		t.Fatalf("tried to access mongo URI %q, got %q", dbURI, err)
	}
}