package acc

import (
	"context"
	"github.com/nullstone-modules/mongo-db-admin/mongodb"
	"github.com/nullstone-modules/mongo-db-admin/workflows"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/url"
	"os"
	"testing"
)

func TestFull(t *testing.T) {
	if os.Getenv("ACC") != "1" {
		t.Skip("Set ACC=1 to run e2e tests")
	}

	connUrl := "mongodb://mda:mda@localhost:27017/admin"
	client, err := mongo.Connect(nil, options.Client().ApplyURI(connUrl))
	require.NoError(t, err, "error connecting to mongo")
	defer client.Disconnect(nil)

	newUser := mongodb.User{
		Name:         "test-user",
		Password:     "test-password",
		RoleName:     "dbOwner",
		DatabaseName: "test-database",
	}
	require.NoError(t, workflows.EnsureUser(client, newUser))

	ctx := context.Background()

	u, _ := url.Parse(connUrl)
	u.Path = "/test-database"
	u.User = url.UserPassword(newUser.Name, newUser.Password)

	appClient, err := mongo.Connect(ctx, options.Client().ApplyURI(u.String()))
	require.NoError(t, err, "error connecting to app mongo")
	defer appClient.Disconnect(ctx)
	appDb := appClient.Database("test-database")

	// Attempt to create collections
	todosSchema := bson.M{
		"bsonType": "object",
		"required": []string{"name"},
		"properties": bson.M{
			"name": bson.M{
				"bsonType":    "string",
				"description": "Name of the todo",
			},
		},
	}
	validator := bson.M{"$jsonSchema": todosSchema}
	opts := options.CreateCollection().SetValidator(validator)
	err = appDb.CreateCollection(ctx, "todos", opts)
	require.NoError(t, err, "create collection")

	// Attempt to insert rows into collection
	docs := []interface{}{
		bson.D{{"name", "item1"}},
		bson.D{{"name", "item2"}},
		bson.D{{"name", "item3"}},
	}
	_, err = appDb.Collection("todos").InsertMany(ctx, docs)
	require.NoError(t, err, "insert todos")

	// Attempt to retrieve them
	cursor, err := appDb.Collection("todos").Find(ctx, bson.D{})
	require.NoError(t, err, "find todos")
	var got []bson.M
	require.NoError(t, cursor.All(ctx, &got), "scan todos")
	results := make([]string, 0)
	for _, result := range got {
		results= append(results, result["name"].(string))
	}
	assert.Equal(t, []string{"item1", "item2", "item3"}, results)
}
