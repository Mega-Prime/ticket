package ticket_test

import (
	"testing"
	"ticket"
)

func TestFields(t *testing.T) {
	t.Parallel()
	_ = ticket.Ticket{
		ID:          "577f9cecd71d71fa1fb6f43a",
		Subject:     "My screen broke",
		Description: "Testing Something",
		Status:      ticket.StatusOpen,
	}
}
