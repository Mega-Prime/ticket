package ticket_test

import (
	"testing"
	"ticket"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestNewTicket(t *testing.T) {
	t.Parallel()
	ts := ticket.NewStore()
	want := ticket.Ticket{
		Status:  ticket.StatusOpen,
		Subject: "My screen broke again",
	}
	got := ts.NewTicket(want.Subject)
	if !cmp.Equal(want, got, cmpopts.IgnoreFields(want, "ID")) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestFields(t *testing.T) {
	t.Parallel()
	ts := ticket.NewStore()
	newTic := ts.NewTicket("testing something")
	newTic.Description = "Pixels missing"
}

func TestID(t *testing.T) {
	t.Parallel()
	ts := ticket.NewStore()
	t1 := ts.NewTicket("test ticket")
	if t1.ID == 0 {
		t.Errorf("invalid id: %v", t1.ID)
	}
	t2 := ts.NewTicket("another test ticket")
	if t1.ID == t2.ID {
		t.Errorf("want different IDs, got both == %v", t1.ID)
	}
}

func TestGetTicket(t *testing.T) {
	//create ticket in system
	t.Parallel()
	s := ticket.NewStore()
	wantSubject := "My screen broke again"

	//created ticket
	createdTicket := s.NewTicket(wantSubject)

	//got the id from created ticket

	got, err := s.Get(createdTicket.ID)
	if err != nil {
		t.Fatal(err)
	}

	//check if it is correct subject
	if wantSubject != got.Subject {
		t.Errorf("want %q, got: %q", wantSubject, got.Subject)
	}

}

func BenchmarkNewTicket(b *testing.B) {
	p := ticket.NewStore()
	for i := 0; i < b.N; i++ {
		p.NewTicket("high speed")
	}
}
