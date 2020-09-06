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

// OpenStore takes an io.Reader, and tries to read JSON data representing a
// slice of tickets.
func OpenStore(r io.Reader) (*Store, error) {
	s := NewStore()
	var tks []Ticket
	decoder := json.NewDecoder(r)
	err := decoder.Decode(&tks)
	if err != nil {
		return nil, err
	}
	for _, tk := range tks {
		s.tickets[tk.ID] = &tk
	}
	return s, nil
}

// Yada.
func (s *Store) WriteTo(w io.Writer) error {
	w = s.tickets
	for _, tk := range w{
		json.NewEncoder(tk).Encode
	}
	encoder := json.NewEncoder(w)
	return nil
}

// AddTicket creates a ticket
func (s *Store) AddTicket(tk Ticket) (int, error) {
	s.highestID++

	//Store id in t
	tk.ID = s.highestID
	//save ticket here in a map:
	s.tickets[tk.ID] = &tk
	return tk.ID, nil
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
