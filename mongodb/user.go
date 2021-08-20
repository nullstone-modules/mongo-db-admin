package mongodb

import (
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
)

type User struct {
	Username string
	Password string
	DatabaseName string
}

func (u User) Create(client *mongo.Client) error {
	fmt.Printf("creating user %q\n", u.Username)
	//client.Database(u.DatabaseName).RunCommand(context.Background(), )
	return nil
}
