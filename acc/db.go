package acc

import (
	"context"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"testing"
)

func createClient(t *testing.T) *mongo.Client {
	connUrl := "mongodb://mda:mda@localhost:27017/admin"
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(connUrl))
	require.NoError(t, err, "error connecting to mongo")
	return client
}
