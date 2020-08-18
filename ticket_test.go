package ticket_test

import (
	"testing"
	"ticket"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

var ignoreID = cmpopts.IgnoreFields(ticket.Ticket{}, "ID")

func TestNewTicket(t *testing.T) {
	t.Parallel()
	s := ticket.NewStore()
	want := &ticket.Ticket{
		Status:  ticket.StatusOpen,
		Subject: "My screen broke again",
	}
	got := s.NewTicket(want.Subject)
	if !cmp.Equal(want, got, ignoreID) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestFields(t *testing.T) {
	t.Parallel()
	s := ticket.NewStore()
	tk := s.NewTicket("testing something")
	tk.Description = "Pixels missing"
}

func TestID(t *testing.T) {
	t.Parallel()
	s := ticket.NewStore()
	t1 := s.NewTicket("test ticket")
	if t1.ID == 0 {
		t.Errorf("invalid id: %v", t1.ID)
	}
	t2 := s.NewTicket("another test ticket")
	// IDs are sequential - potential security issue in the future?
	if t1.ID == t2.ID {
		t.Errorf("want different IDs, got both == %v", t1.ID)
	}
}

func TestGetTicket(t *testing.T) {
	//create ticket in system
	t.Parallel()
	s := ticket.NewStore()

	//pointer to created ticket
	want := s.NewTicket("My screen broke again")

	//look up ticket by ID and get a pointer to it
	got, err := s.Get(want.ID)
	if err != nil {
		t.Fatal(err)
	}

	//both pointers should point to the same ticket
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestGetByStatus(t *testing.T) {
	t.Parallel()
	s := ticket.NewStore()
	tkOpen := s.NewTicket("this is open")
	tkClosed := s.NewTicket("this will be closed")
	tkClosed.Status = ticket.StatusClosed
	want := []*ticket.Ticket{
		tkOpen,
	}
	wantClosed := []*ticket.Ticket{
		tkClosed,
	}
	got, err := s.GetByStatus(ticket.StatusOpen)
	if err != nil {
		t.Error(err)
	}

	gotclosed, err2 := s.GetByStatus(ticket.StatusClosed)
	if err != nil {
		t.Error(err2)
	}
	// if !cmp.Equal(want, got) {
	// 	t.Error(cmp.Diff(want, got))
	// }
	if !cmp.Equal(want, got) {
		t.Errorf("GetByStatus(StatusOpen): %v", cmp.Diff(want, got))
	}
	if !cmp.Equal(wantClosed, gotclosed) {
		t.Errorf("GetByStatus(StatusOpen): %v", cmp.Diff(want, got))
	}

}

func BenchmarkNewTicket(b *testing.B) {
	p := ticket.NewStore()
	for i := 0; i < b.N; i++ {
		p.NewTicket("high speed")
	}
}
