package ticket

import (
	"encoding/json"
	"io"
)

const (
	StatusInvalid int = iota
	StatusOpen
	StatusClosed
)

type ID string

type Ticket struct {
	Subject     string `json:"subject,omitempty"`
	Description string `json:"description,omitempty"`
	ID          ID     `json:"id,omitempty"`
	Status      int    `json:"status"`
}

func (t *Ticket) FromJSON(r io.Reader) error {
	return json.NewDecoder(r).Decode(t)
}

func (t *Ticket) ToJSON(w io.Writer) error {
	return json.NewEncoder(w).Encode(t)
}
