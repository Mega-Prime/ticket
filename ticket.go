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

func (t *Ticket) FromJSON(r io.Reader) error {

	return json.NewDecoder(r).Decode(t)
}

func (t *Ticket) ToJSON(w io.Writer) error {
	return json.NewEncoder(w).Encode(t)
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

// AddTicket creates a ticket
func (s *Store) AddTicket(t Ticket) (int, error) {
	s.highestID++

	//Store id in t
	t.ID = s.highestID
	//save ticket here in a map:
	s.tickets[t.ID] = &t
	return t.ID, nil

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

// create OpenStore func. Takes io.reader, decodes a [] tickets
// tickets need to be added to a new store in sequence. for loop?
// close reader as to not create invalid writes
// How do we write updates to the store? flush? Store.Writeto?
