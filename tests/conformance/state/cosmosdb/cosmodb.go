package cosmosdb

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	defaultTimeout       = 60 * time.Second
	collectionNamePrefix = "test-coll-"
	accountEnv           = "COSMODB_ACCOUNT_NAME"
	masterKeyEnv         = "COSMODB_MASTER_KEY"
	dbnameEnv            = "COSMODB_NAME"
)

var (
	connectionStringTmpl = template.Must(template.New("connectionString").Parse("mongodb://{{.Account}}:{{.MasterKey}}@{{.Account}}.documents.azure.com:10255/?ssl=true"))
	componentYamlTmpl    = template.Must(template.New("componentYaml").Parse(`apiVersion: dapr.io/v1alpha1
kind: Component
metadata:
  name: statestore
spec:
  type: state.azure.cosmosdb
  version: v1
  metadata:
  - name: url
    value: https://{{.Account}}.documents.azure.com:443
  - name: masterKey
    value: {{.MasterKey}}
  - name: database
    value: {{.DatabaseName}}
  - name: collection
    value: {{.CollectionName}}
`))
)

type Collection struct {
	Account        string
	MasterKey      string
	DatabaseName   string
	CollectionName string
	c              *mongo.Collection
}

func (c *Collection) ComponentsYAML() (string, error) {
	var buffer bytes.Buffer
	if err := componentYamlTmpl.Execute(&buffer, c); err != nil {
		return "", err
	}
	return buffer.String(), nil
}

func (c *Collection) ConnectionString() (string, error) {
	var buffer bytes.Buffer
	if err := connectionStringTmpl.Execute(&buffer, c); err != nil {
		return "", err
	}
	return buffer.String(), nil
}

func (c *Collection) Drop(ctx context.Context) error {
	return c.c.Drop(ctx)
}

func NewTestCollection(ctx context.Context) (collection *Collection, err error) {
	var account, masterKey, databaseName string
	if account, err = MustGetenv(accountEnv); err != nil {
		return nil, err
	}
	if masterKey, err = MustGetenv(masterKeyEnv); err != nil {
		return nil, err
	}
	if databaseName, err = MustGetenv(dbnameEnv); err != nil {
		return nil, err
	}
	collectionName := collectionNamePrefix + uuid.New().String()
	collection = &Collection{
		Account:        account,
		MasterKey:      masterKey,
		DatabaseName:   databaseName,
		CollectionName: collectionName,
	}
	connectionString, err := collection.ConnectionString()
	if err != nil {
		return nil, err
	}
	clientOptions := options.Client().ApplyURI(connectionString).SetDirect(true)

	// The following calls are blocking, so we want to establish a timeout to avoid
	// having the tests hanging forever.
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	// First, initialize a connection.
	c, err := mongo.NewClient(clientOptions)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to create mongodb client")
	}
	err = c.Connect(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to initialize connection")
	}
	err = c.Ping(ctx, nil)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to connect")
	}
	if err := c.Database(databaseName).RunCommand(ctx, bson.D{
		{"shardCollection", databaseName + "." + collectionName},
		{"key", bson.D{{"partitionKey", "hashed"}}},
	}).Err(); err != nil {
		return nil, errors.Wrapf(err, "Unable to create a new collection %q", collectionName)
	}
	collection.c = c.Database(databaseName).Collection(collectionName)
	return collection, nil
}

func MustGetenv(env string) (string, error) {
	if val := os.Getenv(env); val != "" {
		return val, nil
	}
	return "", fmt.Errorf("The environment variable %q must be set", env)
}
