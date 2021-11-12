package repository

import (
	"context"
	"github.com/Kamva/mgm/v2"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
)

// Client is a structure that holds the necessary data of a client
type Client struct {
	mgm.DefaultModel `bson:",inline"`

	Email                 string   `json:"email" bson:"email"`
	Name                  string   `json:"name" bson:"name"`
	Surname               string   `json:"surname" bson:"surname"`
}

// AddClient adds a client to clients
func AddClient(newClient *Client) error {
	repo := NewMongoRepository(newClient)
	return repo.Create()
}

func FindFirstByName(name string) (*Client, error) {
	filter := generateFilter(name)
	repo := NewMongoRepository(&Client{})
	cursor, err := repo.Find(filter)
	if err != nil {
		return nil, err
	}

	client := Client{}
	c := context.Background()
	if cursor.Next(c) {
		err := cursor.Decode(&client)
		if err != nil {
			return nil, err
		}
		return &client, nil
	}

	return nil, errors.Errorf("not found for name %v", name)
}

func generateFilter(name string) interface{} {
	return bson .M{"name": name}
}

