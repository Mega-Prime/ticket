package ticket

import (
	"encoding/json"
	"fmt"
	"io"
)

const (
	StatusInvalid int = iota
	StatusOpen
	StatusClosed
)

type Ticket struct {
	Subject     string `json:"subject,omitempty"`
	Description string `json:"description,omitempty"`
	ID          int    `json:"id,omitempty"`
	Status      int    `json:"status"`
}

func (t *Ticket) FromJson(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(t)
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

func (s *Store) GetByID(ID int) (*Ticket, error) {
	tk, ok := s.tickets[ID]
	if !ok {
		return &Ticket{}, fmt.Errorf("no such ID %d", ID)
	}
	return tk, nil
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

func CreateTicket(t Ticket) (tic *Ticket, err error) {
	tk := &Ticket{
		Subject:     t.Subject,
		Description: t.Description,
		ID:          t.ID,
		Status:      t.Status,
	}
	return tk, err

}
