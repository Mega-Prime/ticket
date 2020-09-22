package ticket

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoStore struct {
}

func NewMongoStore(addr string, collectionName string) Store {

	return &MongoStore{}
}

func (s *MongoStore) AddTicket(tk Ticket) (ID, error) {
	// given existing Mongo collection
	dbURI := os.Getenv("MONGO_TICKET_STORE_URL")
	if dbURI == "" {

		return "", errors.New("Database URI not set. Set envVar: MONGO_TICKET_STORE_URL = mongodb://localhost:27017")
	}
	client, err := mongo.NewClient(options.Client().ApplyURI(dbURI))

	if err != nil {
		return "", fmt.Errorf("tried to access mongo URI %q, got %q", dbURI, err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		return "", err
	}
	defer client.Disconnect(ctx)

	ticketCollection := client.Database("dbStore").Collection("ticketTest")
	// turn ticket into BSON data
	// call InsertOne on the collection
	res, err := ticketCollection.InsertOne(ctx, tk)
	if err != nil {
		return "", err
	}

	// convert res.InsertedID to 'ID' type
	dbID, ok := res.InsertedID.(primitive.ObjectID)

	if !ok {
		return "", fmt.Errorf("INTERNAL ERROR: Expected 'string' got %T", res.InsertedID)
	}

	// and return it
	return ID(dbID.Hex()), nil

}

func (s *MongoStore) GetByID(ID) (*Ticket, error) {
	return nil, nil
}
