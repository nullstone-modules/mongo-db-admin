package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/nullstone-modules/mongo-db-admin/mongodb"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"time"
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
	switch event.Type {
	case eventTypeCreateUser:
		return ensureUser(ctx, event.Metadata)
	default:
		return fmt.Errorf("unknown event %q", event.Type)
	}
}

func ensureUser(ctx context.Context, metadata map[string]string) error {
	user := mongodb.User{}
	user.Username, _ = metadata["username"]
	if user.Username == "" {
		return fmt.Errorf("cannot create user: username is required")
	}
	user.Password, _ = metadata["password"]
	if user.Password == "" {
		return fmt.Errorf("cannot create user: password is required")
	}
	user.DatabaseName, _ = metadata["databaseName"]
	if user.DatabaseName == "" {
		return fmt.Errorf("cannot create user: databaseName is required")
	}

	client, err := getClient(ctx)
	if err != nil {
		return fmt.Errorf("error connecting to mongo: %w", err)
	}
	defer client.Disconnect(ctx)

	return user.Create(client)
}

func getClient(ctx context.Context) (*mongo.Client, error) {
	connUrl, err := getConnectionUrl(ctx)
	if err != nil {
		return nil, fmt.Errorf("error retrieving mongo connection url: %w", err)
	}
	timeoutCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	return mongo.Connect(timeoutCtx, options.Client().ApplyURI(connUrl))
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
