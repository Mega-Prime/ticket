package ticket

type Store interface {
	AddTicket(Ticket) (ID, error)
	GetByID(ID) (*Ticket, error)
	GetAll() ([]*Ticket, error)
	//UpdateTicket(id ID, update string)
}
