package ticket

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoStore struct {
	ctx        context.Context
	collection *mongo.Collection
}

func NewMongoStore(ctx context.Context, dbURI string, collectionName string) (Store, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(dbURI))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Mongo at %q: %v", dbURI, err)
	}
	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}
	coll := client.Database("dbStore").Collection(collectionName)
	return &MongoStore{
		ctx:        ctx,
		collection: coll,
	}, nil
}

func (s *MongoStore) AddTicket(tk Ticket) (ID, error) {
	res, err := s.collection.InsertOne(s.ctx, tk)
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

func (s *MongoStore) GetByID(id ID) (*Ticket, error) {
	var tk Ticket
	oid, err := primitive.ObjectIDFromHex(string(id))
	if err != nil {
		return nil, err
	}
	res := s.collection.FindOne(s.ctx, bson.M{"_id": oid})
	err = res.Decode(&tk)
	if err != nil {
		return nil, err
	}
	tk.ID = id
	return &tk, nil
}

func (s *MongoStore) GetAll() ([]*Ticket, error) {
	cur, err := s.collection.Find(s.ctx, bson.D{})

	if err != nil {
		return nil, err
	}
	defer cur.Close(s.ctx)

	var tks []*Ticket

	err = cur.All(s.ctx, &tks)
	if err != nil {
		return nil, err
	}
	return tks, nil
}

// func (s *MongoStore) UpdateTicket(id ID, field string, update string) {
// 	coll := s.collection
// 	oid, _ := primitive.ObjectIDFromHex(string(id))
// 	result, err := coll.UpdateOne(s.ctx, bson.M{"_id": oid},
// 		bson.D{
// 			{"$set", bson.D{{field, update}}}})

// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	fmt.Printf("Updated %v Documents!\n", result.ModifiedCount)
// }
