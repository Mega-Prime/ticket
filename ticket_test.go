package ticket_test

import (
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
	tk := ticket.Ticket{
		Subject: "My screen Broke",
	}
	//pointer to created ticket
	want, _ := s.AddTicket(tk)

	//look up ticket by ID and get a pointer to it
	got, err := s.GetByID(tk.ID)
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
	open := ticket.Ticket{
		Subject: "This is open",
	}
	closed := ticket.Ticket{
		Subject: "This will be closed",
	}
	tkOpen, _ := s.AddTicket(open)
	tkClosed, _ := s.AddTicket(closed)
	closed.Status = ticket.StatusClosed
	want := []ticket.Ticket{
		open, //tkOpen,
	}
	wantClosed := []ticket.Ticket{
		closed, //tkClosed,
	}
	got, err := s.GetByStatus(ticket.StatusOpen)
	if err != nil {
		log.Println(tkOpen)
		t.Error(err)
	}

	gotclosed, err2 := s.GetByStatus(ticket.StatusClosed)
	if err != nil {
		log.Println(tkClosed)
		t.Error(err2)
	}

	if !cmp.Equal(want, got) {
		t.Errorf("GetByStatus(StatusOpen): %v", cmp.Diff(want, got))
	}
	if !cmp.Equal(wantClosed, gotclosed) {
		t.Errorf("GetByStatus(StatusOpen): %v", cmp.Diff(want, got))
	}

}

// func TestOpenStore(t *testing.T) {
// 	t.Parallel()
// 	want := "This is ticket 1"

// 	data := `[{"subject": "This is ticket 1"}]`
// 	reader := bytes.NewBufferString(data)
// 	s, err := ticket.OpenStore(reader)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	tk, err := s.GetByID(1)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	got := tk.Subject

// 	if !cmp.Equal(want, got) {
// 		t.Error(cmp.Diff(want, got))
// 	}
// }

func BenchmarkNewTicket(b *testing.B) {
	p := ticket.NewStore()
	for i := 0; i < b.N; i++ {
		p.AddTicket(ticket.Ticket{})
	}
}
