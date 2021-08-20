package workflows

import (
	"github.com/nullstone-modules/mongo-db-admin/mongodb"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

func EnsureUser(client *mongo.Client, newUser mongodb.User) error {
	log.Printf("ensuring user %q\n", newUser.Name)

	return newUser.Ensure(client)
}
