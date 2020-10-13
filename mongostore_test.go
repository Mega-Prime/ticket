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

const (
	testDBName     = "dbStore"
	testCollection = "ticket_test"
)

func TestAddTicketMongo(t *testing.T) {
	t.Parallel()
	s := newTestMongoStore(t)
	tk1 := ticket.Ticket{
		Subject: "TestAddTicketMongo",
	}
	ID1, err := s.AddTicket(tk1)
	if err != nil {
		t.Fatal(err)
	}
	tk2 := ticket.Ticket{
		Subject: "TestAddTicketMongo",
	}
	ID2, err := s.AddTicket(tk2)
	if err != nil {
		t.Fatal(err)
	}
	if ID1 == ID2 {
		t.Errorf("want different IDs, got both == %v", ID1)
	}
}

func TestGetTicketByID(t *testing.T) {
	t.Parallel()
	s := newTestMongoStore(t)
	// don't care about ID
	_, err := s.AddTicket(ticket.Ticket{})
	if err != nil {
		t.Fatal(err)
	}
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
	if got.ID == "" {
		t.Error("ticket has no ID")
	}
}

func TestGetAll(t *testing.T) {
	// can't run in parallel because we need to know
	// what all the tickets are
	s := newTestMongoStore(t)
	want := []*ticket.Ticket{
		{
			Subject: "This is ticket A",
		},
		{
			Subject: "This is ticket B",
		},
	}
	s.AddTicket(*want[0])
	s.AddTicket(*want[1])
	got, err := s.GetAll()
	if err != nil {
		t.Fatal(err)
	}

	if !cmp.Equal(want, got, ignoreID) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestUpdateTicket(t *testing.T) {
	t.Parallel()
	s := newTestMongoStore(t)
	ID, err := s.AddTicket(ticket.Ticket{
		Subject: "Test UpdateTicket",
	})
	if err != nil {
		t.Fatal(err)
	}

	original, err := s.GetByID(ID)
	if err != nil {
		t.Fatal(err)
	}

	original.Subject = "This has been updated"
	modified := original
	err = s.UpdateTicket(modified)
	if err != nil {
		t.Fatal(err)
	}
	final, err := s.GetByID(modified.ID)
	if err != nil {
		t.Fatal(err)
	}
	if !cmp.Equal(modified, final) {
		cmp.Diff(modified, final)
	}
}

func TestUpdateBogusTicket(t *testing.T) {
	t.Parallel()
	s := newTestMongoStore(t)
	err := s.UpdateTicket(&ticket.Ticket{ID: "0f6a257cf0b6af5783e44658"})
	if err == nil {
		t.Fatal("want error on updating non-existent ticket, got nil")
	}
}

func newTestMongoStore(t *testing.T) ticket.Store {
	t.Helper()
	dbURI := os.Getenv("MONGO_TICKET_STORE_URL")
	if dbURI == "" {
		t.Fatal("Database URI not set:\nexport MONGO_TICKET_STORE_URL=mongodb://localhost:27017")
	}
	t.Cleanup(func() {
		cleanUpTestCollection(t, dbURI)
	})
	s, err := ticket.NewMongoStore(context.Background(), dbURI, testCollection)
	if err != nil {
		t.Fatal(err)
	}
	return s
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
	err = client.Database(testDBName).Collection(testCollection).Drop(ctx)
	if err != nil {
		t.Fatalf("tried to access mongo URI %q, got %q", dbURI, err)
	}
}
