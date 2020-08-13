package ticket

import "fmt"

const (
	StatusOpen int = iota
)

type Ticket struct {
	Subject     string
	Description string
	ID          int
	Status      int
}

type Store struct {
	highestID int
	tickets   map[int]*Ticket
}

func NewStore() *Store {
	return &Store{
		tickets: map[int]*Ticket{},
	}
}

func (s *Store) Get(ID int) (*Ticket, error) {
	tk, ok := s.tickets[ID]
	if !ok {
		return &Ticket{}, fmt.Errorf("no such ID %d", ID)
	}
	return tk, nil
}

func (s *Store) NewTicket(subject string) Ticket {
	s.highestID++
	tk := Ticket{
		Subject: subject,
		ID:      s.highestID,
		Status:  StatusOpen,
	}
	//save ticket here in a map:
	s.tickets[tk.ID] = &tk

	//then return ticket:
	return tk
}
