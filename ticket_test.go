package ticket_test

import (
	"bytes"
	"log"
	"testing"
	"ticket"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

var ignoreID = cmpopts.IgnoreFields(ticket.Ticket{}, "ID")

func TestNewTicket(t *testing.T) {
	t.Parallel()
	s := ticket.NewStore()
	want := ticket.Ticket{
		Status:  ticket.StatusOpen,
		Subject: "My screen broke again",
	}
	got, err := s.AddTicket(want)

	if err != nil {
		t.Fatal(err)
	}
	if cmp.Equal(want, got, ignoreID) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestFields(t *testing.T) {
	t.Parallel()
	s := ticket.NewStore()
	want := ticket.Ticket{
		Description: "Testing Something",
	}
	_, err := s.AddTicket(want)
	if err != nil {
		t.Fatal(err)
	}
	want.Description = "Pixels missing"
}

func TestID(t *testing.T) {
	t.Parallel()
	s := ticket.NewStore()
	tk1 := ticket.Ticket{
		Subject: "Test Ticket",
		ID:      1,
	}
	tk2 := ticket.Ticket{
		Subject: "Ticket 2",
		ID:      2,
	}
	t1, _ := s.AddTicket(tk1)
	if tk1.ID == 0 {
		log.Println("This is id of tk: ", t1)
		t.Errorf("invalid id: %v", t1)
	}
	t2, _ := s.AddTicket(tk2)
	// IDs are sequential - potential security issue in the future?
	if tk1.ID == tk2.ID {
		log.Println("This is id of tk: ", t2)
		t.Errorf("want different IDs, got both == %v", tk1.ID)
	}
}

func TestGetByID(t *testing.T) {
	//create ticket in system
	t.Parallel()
	s := ticket.NewStore()
	want := ticket.Ticket{
		Subject: "My screen Broke",
	}
	//pointer to created ticket
	ID, err := s.AddTicket(want)
	if err != nil {
		t.Fatal(err)
	}

	//look up ticket by ID and get a pointer to it
	got, err := s.GetByID(ID)
	if err != nil {
		t.Fatal(err)
	}

	//both pointers should point to the same ticket
	if !cmp.Equal(&want, got, ignoreID) {
		t.Error(cmp.Diff(&want, got))
	}
}

func TestGetByStatus(t *testing.T) {
	t.Parallel()
	s := ticket.NewStore()
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

func TestOpenStore(t *testing.T) {
	t.Parallel()
	want := "This is a test ticket"

	data := `[{"ID": 99, "subject": "This is a test ticket"}]`
	r := bytes.NewBufferString(data)
	s, err := ticket.OpenStore(r)
	if err != nil {
		t.Fatal(err)
	}
	if s == nil {
		t.Fatal("store is nil")
	}
	tk, err := s.GetByID(99)
	if err != nil {
		t.Fatal(err)
	}
	got := tk.Subject

	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestWriteJSONTo(t *testing.T) {
	t.Parallel()
	want := "[{\"subject\":\"This is a test ticket\",\"id\":1,\"status\":1}]\n"
	var buf = &bytes.Buffer{}
	s := ticket.NewStore()
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
	want := ticket.Ticket{
		Subject: "This is a test ticket",
		ID:      1,
		Status:  ticket.StatusOpen,
	}

	var buf = bytes.NewBufferString("[{\"subject\":\"This is a test ticket\",\"id\":1}]\n")
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

func BenchmarkNewTicket(b *testing.B) {
	p := ticket.NewStore()
	for i := 0; i < b.N; i++ {
		p.AddTicket(ticket.Ticket{})
	}
}
