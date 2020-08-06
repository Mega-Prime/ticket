package ticket_test

import (
	"testing"
	"ticket"

	"github.com/google/go-cmp/cmp"
)

func TestSubject(t *testing.T) {
	var got ticket.Ticket
	want := "My screen broke"
	got = ticket.New(want)
	if want != got.Subject {
		t.Error(cmp.Diff(want, got))
	}
}

func TestID(t *testing.T) {
	got1 := ticket.New("test ticket")
	if got1.ID == 0 {
		t.Errorf("invalid id: %v", got1.ID)
	}
	got2 := ticket.New("another test ticket")
	if got1.ID == got2.ID {
		t.Errorf("want different IDs, got both == %v", got1.ID)
	}
}
