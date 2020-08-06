package ticket

type Ticket struct {
	Subject string
	ID      int
}

type Project struct {
	Name      string
	highestID int
}

func NewProject(name string) *Project {
	return &Project{
		Name: name,
	}
}

func (p *Project) NewTicket(s string) Ticket {
	p.highestID++
	return Ticket{
		Subject: s,
		ID:      p.highestID,
	}
}
