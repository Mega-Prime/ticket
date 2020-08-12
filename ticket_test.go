package ticket_test

import (
	"testing"
	"ticket"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestNewTicket(t *testing.T) {
	t.Parallel()
	p := ticket.NewProject("test")
	want := ticket.Ticket{
		Status:  ticket.StatusOpen,
		Subject: "My screen broke again",
		//Description: "Pixels missing!",
	}
	got := p.NewTicket(want.Subject)
	if !cmp.Equal(want, got, cmpopts.IgnoreFields(want, "ID")) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestFields(t *testing.T) {
	t.Parallel()
	p := ticket.NewProject("test field")
	newTic := p.NewTicket("testing something")
	newTic.Description = "Pixels missing"
}

func TestID(t *testing.T) {
	t.Parallel()
	p := ticket.NewProject("test")
	t1 := p.NewTicket("test ticket")
	if t1.ID == 0 {
		t.Errorf("invalid id: %v", t1.ID)
	}
	t2 := p.NewTicket("another test ticket")
	if t1.ID == t2.ID {
		t.Errorf("want different IDs, got both == %v", t1.ID)
	}
}

func TestGetTicket(t *testing.T) {
	//create ticket in system
	t.Parallel()
	p := ticket.NewProject("test")
	want := "My screen broke again"

	//created ticket
	createTicket := p.NewTicket("testing ticket search")

	//got the id from created ticket

	got := ticket.Get(createTicket.ID)

	//check if it is correct subject
	if want != got.Subject {
		t.Errorf("want %s, got: %s", want, got.Subject)
	}

}

func BenchmarkNewTicket(b *testing.B) {
	p := ticket.NewProject("speedtest")
	for i := 0; i < b.N; i++ {
		p.NewTicket("high speed")
	}
}
