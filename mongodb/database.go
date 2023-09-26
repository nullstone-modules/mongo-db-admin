package mongodb

import (
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

type Database struct {
	Name string
}

func (d Database) Ensure(client *mongo.Client) error {
	if exists, err := u.Exists(client); exists {
		log.Printf("User %q already exists\n", u.Name)
		return nil
	} else if err != nil {
		return fmt.Errorf("error checking for user %q: %w", u.Name, err)
	}
	if err := u.Create(client); err != nil {
		return fmt.Errorf("error creating user %q: %w", u.Name, err)
	}
	return nil
}
