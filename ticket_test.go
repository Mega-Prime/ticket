package ticket_test

import (
	"bytes"
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
	got := s.AddTicket(want.Subject)
	if !cmp.Equal(want, got, ignoreID) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestFields(t *testing.T) {
	t.Parallel()
	s := ticket.NewStore()
	tk := s.AddTicket("testing something")
	tk.Description = "Pixels missing"
}

func TestID(t *testing.T) {
	t.Parallel()
	s := ticket.NewStore()
	t1 := s.AddTicket("test ticket")
	if t1.ID == 0 {
		t.Errorf("invalid id: %v", t1.ID)
	}
	t2 := s.AddTicket("another test ticket")
	// IDs are sequential - potential security issue in the future?
	if t1.ID == t2.ID {
		t.Errorf("want different IDs, got both == %v", t1.ID)
	}
}

func TestGetByID(t *testing.T) {
	//create ticket in system
	t.Parallel()
	s := ticket.NewStore()

	//pointer to created ticket
	want := s.AddTicket("My screen broke again")

	//look up ticket by ID and get a pointer to it
	got, err := s.GetByID(want.ID)
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
	tkOpen := s.AddTicket("this is open")
	tkClosed := s.AddTicket("this will be closed")
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

	if !cmp.Equal(want, got) {
		t.Errorf("GetByStatus(StatusOpen): %v", cmp.Diff(want, got))
	}
	if !cmp.Equal(wantClosed, gotclosed) {
		t.Errorf("GetByStatus(StatusOpen): %v", cmp.Diff(want, got))
	}

}

func TestOpenStore(t *testing.T) {
	t.Parallel()
	want := "This is ticket 1"

	data := `[{"subject": "This is ticket 1"}]`
	reader := bytes.NewBufferString(data)
	s, err := ticket.OpenStore(reader)
	if err != nil {
		t.Fatal(err)
	}
	tk, err := s.GetByID(1)
	if err != nil {
		t.Fatal(err)
	}
	got := tk.Subject

	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func BenchmarkNewTicket(b *testing.B) {
	p := ticket.NewStore()
	for i := 0; i < b.N; i++ {
		p.AddTicket("high speed")
	}
}
