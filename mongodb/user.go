package mongodb

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

type User struct {
	Name         string
	Password     string
	RoleName     string
	DatabaseName string
}

func (u User) Ensure(client *mongo.Client) error {
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

func (u User) Create(client *mongo.Client) error {
	fmt.Printf("creating user %q\n", u.Name)
	command := bson.D{
		{"createUser", u.Name},
		{"pwd", u.Password},
		{"roles", []bson.M{{"role": u.RoleName, "db": u.DatabaseName}}},
	}
	result := client.Database("admin").RunCommand(nil, command)
	if err := result.Err(); err != nil {
		return fmt.Errorf("error creating user: %w", err)
	}
	return nil
}

func (u User) Exists(client *mongo.Client) (bool, error) {
	check := User{Name: u.Name}
	if err := check.Read(client); err != nil {
		if err == mongo.ErrNoDocuments {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (u User) Read(client *mongo.Client) error {
	var result bson.M
	err := client.Database("admin").Collection("system.users").
		FindOne(nil, bson.D{{"user", u.Name}}).
		Decode(&result)
	if err != nil {
		return err
	}
	return nil
}
