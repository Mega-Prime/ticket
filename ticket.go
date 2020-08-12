package ticket

const (
	StatusOpen int = iota
)

type Ticket struct {
	Subject     string
	Description string
	ID          int
	Status      int
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


func Get(id int) Ticket {
	return Ticket{}
}

func (p *Project) NewTicket(s string) Ticket {
	p.highestID++
	p.ProjDescription = "Pixels missing!"
	//save ticket here in a map:

	//then return ticket:
	return Ticket{
		Subject:     s,
		ID:          p.highestID,
		Description: p.ProjDescription,
		Status:      StatusOpen,
	}
}
