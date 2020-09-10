package ticket

import (
	"encoding/json"
	"fmt"
	"io"
)

type ticketMap map[int]*Ticket

// MemoryStore stores tickets in memory.
type MemoryStore struct {
	highestID int
	tickets   ticketMap
}

// NewMemoryStore returns a Store value containing a MemoryStore.
// A MemoryStore is not a Store. A pointer to MemoryStore is a Store.
func NewMemoryStore() Store {
	ms := MemoryStore{
		tickets: ticketMap{},
	}
	return &ms
}

// AddTicket creates a ticket
func (s *MemoryStore) AddTicket(tk Ticket) (int, error) {
	s.highestID++

	// MemoryStore id in t
	tk.ID = s.highestID
	// save ticket here in a map:
	s.tickets[tk.ID] = &tk
	return tk.ID, nil
}

func (s *MemoryStore) GetByID(ID int) (*Ticket, error) {
	tk, ok := s.tickets[ID]
	if !ok {
		return &Ticket{}, fmt.Errorf("no such ID %d", ID)
	}
	return tk, nil
}

func (s *MemoryStore) GetByStatus(Status int) (tix []*Ticket, err error) {
	result := []*Ticket{}
	for _, ticket := range s.tickets {
		if ticket.Status == Status {
			result = append(result, ticket)
		}
	}
	return result, err
}

// ReadJSONFrom takes an io.Reader, and tries to read JSON data representing a
// slice of tickets.
func ReadJSONFrom(r io.Reader) (*MemoryStore, error) {
	s := NewMemoryStore().(*MemoryStore)
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

// WriteTo writes encodes data to JSON and writes to file.
func (s *MemoryStore) WriteJSONTo(w io.Writer) error {
	tks := []Ticket{}
	for _, tk := range s.tickets {
		tks = append(tks, *tk)
	}

	return json.NewEncoder(w).Encode(tks)
}
