package ticket_test

import (
	"testing"
	"ticket"

	"github.com/google/go-cmp/cmp"
)

func TestSubject(t *testing.T) {
	t.Parallel()
	p := ticket.NewProject("test")
	var t1 ticket.Ticket
	want := "My screen broke"
	t1 = p.NewTicket(want)
	if want != t1.Subject {
		t.Error(cmp.Diff(want, t1))
	}
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

func BenchmarkNewTicket(b *testing.B) {
	p := ticket.NewProject("speedtest")
	for i := 0; i < b.N; i++ {
		p.NewTicket("high speed")
	}
}
