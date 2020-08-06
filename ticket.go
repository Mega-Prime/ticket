package ticket

type Ticket struct {
	Subject string
	ID      int
}

// New returns struct values
func New(s string) Ticket {
	return Ticket{Subject: s, ID: 1} // Should there be an auto-generated number using the Math/random package?
}
