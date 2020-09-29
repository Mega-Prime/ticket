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
	testDB         = "dbStore"
	testCollection = "ticket_test"
)

func TestAddTicketMongo(t *testing.T) {
	//t.Parallel()
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
	// IDs are sequential - potential security issue in the future?
	if ID1 == ID2 {
		t.Errorf("want different IDs, got both == %v", ID1)
	}

}

func TestGetTicketByID(t *testing.T) {
	//t.Parallel()
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
	//t.Parallel()
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

// func TestUpdateTicket(t *testing.T) {
// 	s := newTestMongoStore(t)
// 	tk, err := s.AddTicket(ticket.Ticket{
// 		Subject:     "This is a Test",
// 		Description: "This is a start",
// 	})
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	want := ticket.Ticket{
// 		Subject:     "This is a test",
// 		Description: "This has been updated",
// 	}

// 	got := s.UpdateTicket()
// 	if !cmp.Equal(&want, got, ignoreID) {
// 		t.Error(cmp.Diff(&want, got))
// 	}
// }

func newTestMongoStore(t *testing.T) ticket.Store {
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
