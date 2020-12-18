// +build conf

package state

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/dapr/components-contrib/state"
	"github.com/dapr/components-contrib/tests/conformance"
	"github.com/stretchr/testify/assert"
)

type ValueType struct {
	Message string `json:"message"`
}

type TestConfig struct {
	maxInitDurationInMs   time.Duration
	maxSetDurationInMs    time.Duration
	maxGetDurationInMs    time.Duration
	maxDeleteDurationInMs time.Duration
	numBulkRequests       int
}

func runWithStateStore(t *testing.T, name string, componentFactory func() state.Store, config TestConfig) {
	store := componentFactory()
	comps, err := conformance.LoadComponents(fmt.Sprintf("../../config/state/%s", name))
	assert.Nil(t, err)
	assert.Equal(t, len(comps), 1) // We only expect a single component per state store

	c := comps[0]
	props := conformance.ConvertMetadataToProperties(c.Spec.Metadata)
	// Run the state store conformance tests
	stateStoreConformanceTests(t, props, store, config)
}

/*
	State store component tests
*/
func stateStoreConformanceTests(t *testing.T, props map[string]string, statestore state.Store, config TestConfig) {
	// Test vars
	key := conformance.NewRandString(8)
	b, _ := json.Marshal(ValueType{Message: "test"})
	value := b

	// Init
	t.Run("init", func(t *testing.T) {
		start := time.Now()
		err := statestore.Init(state.Metadata{
			Properties: props,
		})
		elapsed := time.Since(start)
		assert.Nil(t, err)
		assert.Lessf(t, elapsed.Microseconds(), config.maxInitDurationInMs.Microseconds(),
			"test took %dμs but must complete in less than %dμs", elapsed.Microseconds(), config.maxDeleteDurationInMs.Microseconds())
	})

	// Set
	t.Run("set", func(t *testing.T) {
		setReq := &state.SetRequest{
			Key:   key,
			Value: value,
		}
		start := time.Now()
		err := statestore.Set(setReq)
		elapsed := time.Since(start)
		assert.Nil(t, err)
		assert.Lessf(t, elapsed.Microseconds(), config.maxSetDurationInMs.Microseconds(),
			"test took %dμs but must complete in less than %dμs", elapsed.Microseconds(), config.maxSetDurationInMs.Microseconds())
	})

	// Get
	t.Run("get", func(t *testing.T) {
		getReq := &state.GetRequest{
			Key: key,
		}
		start := time.Now()
		getRes, err := statestore.Get(getReq) // nolint:govet
		elapsed := time.Since(start)
		assert.Nil(t, err)
		assert.Equal(t, value, getRes.Data)
		assert.Lessf(t, elapsed.Microseconds(), config.maxGetDurationInMs.Microseconds(),
			"test took %dμs but must complete in less than %dμs", elapsed.Microseconds(), config.maxGetDurationInMs.Microseconds())
	})

	// Delete
	t.Run("delete", func(t *testing.T) {
		delReq := &state.DeleteRequest{
			Key: key,
		}
		start := time.Now()
		err := statestore.Delete(delReq)
		elapsed := time.Since(start)
		assert.Nil(t, err)
		assert.Lessf(t, elapsed.Microseconds(), config.maxDeleteDurationInMs.Microseconds(),
			"test took %dμs but must complete in less than %dμs", elapsed.Microseconds(), config.maxDeleteDurationInMs.Microseconds())
	})

	// Bulk test vars
	var bulkSetReqs []state.SetRequest
	var bulkDeleteReqs []state.DeleteRequest
	for k := 0; k < config.numBulkRequests; k++ {
		bkey := fmt.Sprintf("%s-%d", key, k)
		bulkSetReqs = append(bulkSetReqs, state.SetRequest{
			Key:   bkey,
			Value: value,
		})
		bulkDeleteReqs = append(bulkDeleteReqs, state.DeleteRequest{
			Key: bkey,
		})
	}

	// BulkSet
	t.Run("bulkset", func(t *testing.T) {
		start := time.Now()
		err := statestore.BulkSet(bulkSetReqs)
		elapsed := time.Since(start)
		maxElapsed := config.maxSetDurationInMs * time.Duration(config.numBulkRequests) // assumes at least linear scale
		assert.Nil(t, err)
		assert.Lessf(t, elapsed.Microseconds(), maxElapsed.Microseconds(),
			"test took %dμs but must complete in less than %dμs", elapsed.Microseconds(), maxElapsed.Microseconds())
		for k := 0; k < config.numBulkRequests; k++ {
			bkey := fmt.Sprintf("%s-%d", key, k)
			greq := &state.GetRequest{
				Key: bkey,
			}
			_, err = statestore.Get(greq)
			assert.Nil(t, err)
		}
	})

	// BulkDelete
	t.Run("bulkdelete", func(t *testing.T) {
		start := time.Now()
		err := statestore.BulkDelete(bulkDeleteReqs)
		elapsed := time.Since(start)
		maxElapsed := config.maxDeleteDurationInMs * time.Duration(config.numBulkRequests) // assumes at least linear scale
		assert.Nil(t, err)
		assert.Lessf(t, elapsed.Microseconds(), maxElapsed.Microseconds(),
			"test took %dμs but must complete in less than %dμs", elapsed.Microseconds(), maxElapsed.Microseconds())
	})
}
