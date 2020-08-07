package ticket

// const (
// 	Open int = iota
// 	InProgress
// 	Closed
// )

var StatusOpen = "Open"

/*
Why Doesnt the map work when used in NewTicket function:
i.e.
	Status:      StatusOpen[Open],

var StatusOpen = map[int]string{
	Open:       "open",
	InProgress: "In Progress",
	Closed:     "Complete",
}
*/
type Ticket struct {
	Subject     string
	Description string
	ID          int
	Status      string
}

type Project struct {
	Name            string
	highestID       int
	ProjDescription string
}

func NewProject(name string) *Project {
	return &Project{
		Name: name,
	}
}

func (p *Project) NewTicket(s string) Ticket {
	p.highestID++
	p.ProjDescription = "Pixels missing!"

	return Ticket{
		Subject:     s,
		ID:          p.highestID,
		Description: p.ProjDescription,
		Status:      StatusOpen,
	}
}
