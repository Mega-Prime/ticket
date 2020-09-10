package ticket

type MongoStore struct {
}

func NewMongoStore(addr string) Store {
	return &MongoStore{}
}

func (s *MongoStore) AddTicket(Ticket) (ID, error) {
	// given existing Mongo collection
	// turn ticket into BSON data
	// call InsertOne on the collection
	// convert res.InsertedID to 'ID' type
	// and return it
	return "0", nil
}

func (s *MongoStore) GetByID(ID) (*Ticket, error) {
	return nil, nil
}
