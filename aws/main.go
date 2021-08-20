package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/nullstone-modules/mongo-db-admin/mongodb"
	"github.com/nullstone-modules/mongo-db-admin/workflows"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
)

const (
	dbConnUrlSecretIdEnvVar = "DB_CONN_URL_SECRET_ID"

	eventTypeCreateUser = "create-user"
)

type AdminEvent struct {
	Type     string            `json:"type"`
	Metadata map[string]string `json:"metadata"`
}

func main() {
	lambda.Start(HandleRequest)
}

func HandleRequest(ctx context.Context, event AdminEvent) error {
	connUrl, err := getConnectionUrl(ctx)
	if err != nil {
		return err
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connUrl))
	if err != nil {
		return fmt.Errorf("error connecting to db: %w", err)
	}
	defer client.Disconnect(ctx)

	switch event.Type {
	case eventTypeCreateUser:
		newUser := mongodb.User{
			RoleName: "dbOwner", // readWrite|dbAdmin|userAdmin
		}
		newUser.Name, _ = event.Metadata["username"]
		if newUser.Name == "" {
			return fmt.Errorf("cannot create user: username is required")
		}
		newUser.Password, _ = event.Metadata["password"]
		if newUser.Password == "" {
			return fmt.Errorf("cannot create user: password is required")
		}
		newUser.DatabaseName, _ = event.Metadata["databaseName"]
		if newUser.DatabaseName == "" {
			return fmt.Errorf("cannot create user: databaseName is required")
		}
		return workflows.EnsureUser(client, newUser)
	default:
		return fmt.Errorf("unknown event %q", event.Type)
	}
}

func getConnectionUrl(ctx context.Context) (string, error) {
	awsConfig, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return "", fmt.Errorf("error accessing aws: %w", err)
	}
	sm := secretsmanager.NewFromConfig(awsConfig)
	secretId := os.Getenv(dbConnUrlSecretIdEnvVar)
	out, err := sm.GetSecretValue(ctx, &secretsmanager.GetSecretValueInput{SecretId: aws.String(secretId)})
	if err != nil {
		return "", fmt.Errorf("error accessing secret: %w", err)
	}
	if out.SecretString == nil {
		return "", nil
	}
	return *out.SecretString, nil
}
