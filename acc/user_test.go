package acc

import (
	"context"
	"os"
	"testing"
)

func TestDatabase(t *testing.T) {
	if os.Getenv("ACC") != "1" {
		t.Skip("Set ACC=1 to run e2e tests")
	}

	client := createClient(t)
	defer client.Disconnect(context.Background())
	//
	//database := postgresql.Database{Name: "test-database"}
	//
	//ownerRole := postgresql.Role{Name: database.Name}
	//require.NoError(t, ownerRole.Ensure(client), "error creating owner role")
	//database.Owner = ownerRole.Name
	//
	//require.NoError(t, database.Create(client, *dbInfo), "unexpected error")
	//
	//find := &postgresql.Database{Name: "test-database"}
	//require.NoError(t, find.Read(client), "read database")
	//assert.Equal(t, ownerRole.Name, find.Owner, "mismatched owner")
}
