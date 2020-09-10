package ticket

type Store interface {
	AddTicket(Ticket) (int, error)
	GetByID(int) (*Ticket, error)
}
