package store

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/finest08/jwt-auth-demo/model"
)

func (s *Store) AddPerson(g model.PersonCreate) {
	insertResult, err := s.localColl.InsertOne(context.Background(), g)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("\nInserted a Single Document: %v\n", insertResult)
}

func (s *Store) GetPerson(id string) (model.PersonCreate, error) {
	var p model.PersonCreate
	if err := s.localColl.FindOne(
		context.Background(),
		bson.M{"id": id},
	).Decode(&p); err != nil {
		return model.PersonCreate{}, err
	}

	return p, nil
}

func (s *Store) VerifyPerson(email string) (model.PersonCreate, error) {
	var g model.PersonCreate
	if err := s.localColl.FindOne(
		context.Background(),
		bson.M{"email": email},
	).Decode(&g); err != nil {
		return model.PersonCreate{Password: g.Password}, err
	}

	return g, nil
}

func (s *Store) PersonDetail(email string) (model.Person, error) {
	var n model.Person
	if err := s.localColl.FindOne(
		context.Background(),
		bson.M{"email": email},
	).Decode(&n); err != nil {
		return model.Person{GivenName: n.GivenName, FamilyName: n.FamilyName, Email: n.Email}, err
	}

	return n, nil
}

func (s *Store) UpdatePerson(id string, p model.PersonCreate) {
	insertResult, err := s.localColl.ReplaceOne(context.Background(), bson.M{"id": id}, p)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("\nInserted a Single Document: %v\n", insertResult)
}

func (s *Store) DeletePerson(id string) error {
	removeResult, err := s.localColl.DeleteOne(context.Background(), bson.M{"id": id})
	if err != nil {
		return err
	}
	fmt.Printf("\nRemoved a Single Document: %v\n", removeResult)
	return nil
}

func (s *Store) GetPeople(fn, ln, searchText string, limit, skip *int64) (model.Page, error) {
	filter := bson.M{}

	if searchText != "" {
		filter = bson.M{"$and": bson.A{filter,
			bson.M{"$text": bson.M{"$search": searchText}},
		}}
	}

	opt := options.FindOptions{
		Skip:  skip,
		Limit: limit,
		Sort:  bson.M{"date": -1},
	}

	mctx := context.Background()
	cursor, err := s.localColl.Find(mctx, filter, &opt)
	if err != nil {
		return model.Page{}, err
	}

	// unpack results
	var pg model.Page
	if err := cursor.All(mctx, &pg.Data); err != nil {
		return model.Page{}, err
	}
	if pg.Matches, err = s.localColl.CountDocuments(mctx, filter); err != nil {
		return model.Page{}, err
	}
	return pg, nil
}
