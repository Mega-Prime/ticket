//+build integration

package ticket_test

import (
	"context"
	"os"
	"testing"
	"ticket"
	"time"

	"github.com/google/go-cmp/cmp"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func cleanUpTestCollection(t *testing.T, dbURI string) {

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
	err = client.Database("dbStore").Collection("ticketTest").Drop(ctx)
	if err != nil {
		t.Fatalf("tried to access mongo URI %q, got %q", dbURI, err)
	}

}

func TestAddTicketMongo(t *testing.T) {
	t.Parallel()
	dbURI := os.Getenv("MONGO_TICKET_STORE_URL")
	if dbURI == "" {

		t.Fatal("Database URI not set. Set envVar: MONGO_TICKET_STORE_URL = mongodb://localhost:27017")
	}

	t.Cleanup(func() {
		cleanUpTestCollection(t, dbURI)
	})

	s := ticket.NewMongoStore(dbURI, "ticketTest")
	tk1 := ticket.Ticket{
		Subject: "Test Ticket",
	}
	ID1, err := s.AddTicket(tk1)
	if err != nil {
		t.Fatal(err)
	}
	tk2 := ticket.Ticket{
		Subject: "Ticket 2",
	}
	ID2, err := s.AddTicket(tk2)
	if err != nil {
		t.Fatal(err)
	}
	// IDs are sequential - potential security issue in the future?
	if ID1 == ID2 {
		t.Errorf("want different IDs, got both == %v", ID1)
	}

}

func TestGetTicketByID(t *testing.T) {
	t.Parallel()

	dbURI := os.Getenv("MONGO_TICKET_STORE_URL")
	if dbURI == "" {

		t.Fatal("Database URI not set. Set envVar: MONGO_TICKET_STORE_URL = mongodb://localhost:27017")
	}

	t.Cleanup(func() {
		cleanUpTestCollection(t, dbURI)
	})

	s := ticket.NewMongoStore(dbURI, "ticketTest")
	want := ticket.Ticket{
		Subject: "Test Ticket",
	}
	ID, err := s.AddTicket(want)
	if err != nil {
		t.Fatal(err)
	}

	got, err := s.GetByID(ID)
	if err != nil {
		t.Fatal(err)
	}

	if !cmp.Equal(&want, got, ignoreID) {
		t.Error(cmp.Diff(&want, got))
	}

}
