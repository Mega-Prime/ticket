package ticket_test

import (
	"testing"
	"ticket"
)

func TestTicket(t *testing.T) {
	var got ticket.Ticket
	got = ticket.New("My screen broke")
	var want ticket.Ticket
	want = ticket.Ticket{Subject: "My screen broke"}
	if want != got {
		t.Errorf("want %v, got %v", want, got)
	}

}
