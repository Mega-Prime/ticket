// +build integration

package ticket_test

import (
	"testing"
	"ticket"
)

func TestAddTicketMongo(t *testing.T) {
	t.Parallel()
	s := ticket.NewMongoStore("mongo://localhost:27017")
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
