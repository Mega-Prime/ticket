package ticket_test

import (
	"testing"
	"ticket"
)

func TestFields(t *testing.T) {
	t.Parallel()
	_ = ticket.Ticket{
		ID:          1,
		Subject:     "My screen broke",
		Description: "Testing Something",
		Status:      ticket.StatusOpen,
	}
}
