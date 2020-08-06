package ticket

type Ticket struct {
	Subject string
	ID      int
}

var highestID int

// New returns struct values
func New(s string) Ticket {
	highestID++
	return Ticket{
		Subject: s,
		ID:      highestID,
	}
}
