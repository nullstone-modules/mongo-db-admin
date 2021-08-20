package acc

import (
	"github.com/nullstone-modules/mongo-db-admin/mongodb"
	"github.com/nullstone-modules/mongo-db-admin/workflows"
	"github.com/stretchr/testify/require"
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

	u, _ := url.Parse(connUrl)
	u.Path = "/test-database"
	u.User = url.UserPassword(newUser.Name, newUser.Password)

	appClient, err := mongo.Connect(nil, options.Client().ApplyURI(u.String()))
	require.NoError(t, err, "error connecting to app mongo")
	defer appClient.Disconnect(nil)

	
}
