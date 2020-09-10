package ticket_test

import (
	"bytes"
	"testing"
	"ticket"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

var ignoreID = cmpopts.IgnoreFields(ticket.Ticket{}, "ID")

func TestAddTicketMem(t *testing.T) {
	t.Parallel()
	s := ticket.NewMemoryStore()
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

func TestGetByIDMem(t *testing.T) {
	// create ticket in system
	t.Parallel()
	s := ticket.NewMemoryStore()
	want := ticket.Ticket{
		Subject: "My screen Broke",
	}
	// pointer to created ticket
	ID, err := s.AddTicket(want)
	if err != nil {
		t.Fatal(err)
	}

	// look up ticket by ID and get a pointer to it
	got, err := s.GetByID(ID)
	if err != nil {
		t.Fatal(err)
	}

	// both pointers should point to the same ticket
	if !cmp.Equal(&want, got, ignoreID) {
		t.Error(cmp.Diff(&want, got))
	}
}

func TestGetByStatusMem(t *testing.T) {
	t.Parallel()
	s := ticket.NewMemoryStore().(*ticket.MemoryStore)
	wantOpen := []*ticket.Ticket{
		{
			Subject: "This is open",
			Status:  ticket.StatusOpen,
		},
	}
	wantClosed := []*ticket.Ticket{
		{
			Subject: "This is closed",
			Status:  ticket.StatusClosed,
		},
	}
	s.AddTicket(*wantOpen[0])
	s.AddTicket(*wantClosed[0])
	gotOpen, err := s.GetByStatus(ticket.StatusOpen)
	if err != nil {
		t.Error(err)
	}
	if !cmp.Equal(wantOpen, gotOpen, ignoreID) {
		t.Errorf("GetByStatus(StatusOpen): %v", cmp.Diff(wantOpen, gotOpen))
	}
	gotClosed, err := s.GetByStatus(ticket.StatusClosed)
	if err != nil {
		t.Error(err)
	}
	if !cmp.Equal(wantClosed, gotClosed, ignoreID) {
		t.Errorf("GetByStatus(StatusClosed): %v", cmp.Diff(wantClosed, gotClosed))
	}
}

func TestWriteJSONTo(t *testing.T) {
	t.Parallel()
	want := "[{\"subject\":\"This is a test ticket\",\"id\":1,\"status\":1}]\n"
	buf := &bytes.Buffer{}
	s := ticket.NewMemoryStore().(*ticket.MemoryStore)
	s.AddTicket(ticket.Ticket{
		Subject: "This is a test ticket",
		Status:  ticket.StatusOpen,
	})
	err := s.WriteJSONTo(buf)
	if err != nil {
		t.Fatal(err)
	}

	got := buf.String()

	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestReadJSONFrom(t *testing.T) {
	t.Parallel()
	want := &ticket.Ticket{
		Subject: "This is a test ticket",
		ID:      1,
		Status:  ticket.StatusOpen,
	}
	buf := bytes.NewBufferString("[{\"subject\":\"This is a test ticket\",\"id\":1,\"status\":1}]\n")
	s, err := ticket.ReadJSONFrom(buf)
	if err != nil {
		t.Fatal(err)
	}
	got, err := s.GetByID(want.ID)
	if err != nil {
		t.Fatal(err)
	}

	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}
