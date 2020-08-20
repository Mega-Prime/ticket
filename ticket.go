package ticket

import (
	"fmt"
)

const (
	StatusOpen int = iota
	StatusClosed
)

type Ticket struct {
	Subject     string `json:"subject,omitempty"`
	Description string `json:"description,omitempty"`
	ID          int    `json:"id,omitempty"`
	Status      int    `json:"status"`
}

type ticketMap map[int]*Ticket

// Store stores tickets.
type Store struct {
	highestID int
	tickets   ticketMap
}

// NewStore returns a pointer to a new Store.
func NewStore() *Store {
	return &Store{
		tickets: ticketMap{},
	}
}

func (s *Store) GetByID(ID int) (*Ticket, error) {
	tk, ok := s.tickets[ID]
	if !ok {
		return &Ticket{}, fmt.Errorf("no such ID %d", ID)
	}
	return tk, nil
}

func (s *Store) NewTicket(subject string) *Ticket {
	s.highestID++
	tk := &Ticket{
		Subject: subject,
		ID:      s.highestID,
		Status:  StatusOpen,
	}
	//save ticket here in a map:
	s.tickets[tk.ID] = tk

	//then return ticket:
	return tk
}

func (s *Store) GetByStatus(Status int) (tix []*Ticket, err error) {
	result := []*Ticket{}
	for _, ticket := range s.tickets {
		if ticket.Status == Status {
			result = append(result, ticket)
		}

	}
	return result, err
}
