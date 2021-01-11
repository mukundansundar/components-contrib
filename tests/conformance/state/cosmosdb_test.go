// +build conf-dev

package state

import (
	"context"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/dapr/components-contrib/state"
	"github.com/dapr/components-contrib/state/azure/cosmosdb"
	testcosmosdb "github.com/dapr/components-contrib/tests/conformance/state/cosmosdb"
	"github.com/dapr/dapr/pkg/logger"
)

func TestCosmosDB(t *testing.T) {
	ctx := context.Background()
	c, err := testcosmosdb.NewTestCollection(ctx)
	if err != nil {
		t.Fatal("Failed to create cosmo db collection for testing", err)
	}
	defer c.Drop(ctx)
	// TODO: Use a temp location
	base := "../../config/state/cosmodb"
	componentFilePath := base + "/statestore.yaml"
	if err = os.MkdirAll(base, 0755); err != nil {
		t.Fatal("Can't create directory", err)
	}
	yaml, err := c.ComponentsYAML()
	if err != nil {
		t.Fatal("Can't render YAML for cosmosdb", err)
	}
	err = ioutil.WriteFile(componentFilePath, []byte(yaml), 0644)
	if err != nil {
		t.Fatal("Failed to create components YAML file", err)
	}
	defer func() {
		os.Remove(componentFilePath)
	}()
	runWithStateStore(t, "cosmodb", func() state.Store {
		return cosmosdb.NewCosmosDBStateStore(logger.NewLogger("test-cosmos"))
	}, TestConfig{
		maxInitDurationInMs:   time.Duration(1000) * time.Millisecond,
		maxSetDurationInMs:    time.Duration(1000) * time.Millisecond,
		maxDeleteDurationInMs: time.Duration(1000) * time.Millisecond,
		maxGetDurationInMs:    time.Duration(1000) * time.Millisecond,
		numBulkRequests:       10,
	})
}
