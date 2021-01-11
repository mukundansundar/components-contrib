// +build conf

package pubsub

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/dapr/components-contrib/pubsub"
	"github.com/dapr/components-contrib/tests/conformance"
	"github.com/stretchr/testify/assert"
)

type TestConfig struct {
	pubsubName          string
	testTopicName       string
	publishMetadata     map[string]string
	subscribeMetadata   map[string]string
	messageCount        int
	maxReadDurationInMs time.Duration
}

func runWithPubSub(t *testing.T, name string, componentFactory func() pubsub.PubSub, config TestConfig) {
	store := componentFactory()
	comps, err := conformance.LoadComponents(fmt.Sprintf("../../config/pubsub/%s", name))
	assert.Nil(t, err)
	assert.Equal(t, len(comps), 1) // We only expect a single component per state store

	c := comps[0]
	props := conformance.ConvertMetadataToProperties(c.Spec.Metadata)
	// Run the pubSub conformance tests
	pubSubConformanceTests(t, props, store, config)
}

func pubSubConformanceTests(t *testing.T, props map[string]string, store pubsub.PubSub, config TestConfig) {
	actualReadCount := 0
	t.Run("Init and Subscribe to pubsub", func(t *testing.T) {
		err := store.Init(pubsub.Metadata{
			Properties: props,
		})
		assert.NoError(t, err, "expected no error on setting up pubsub")
		err = store.Subscribe(pubsub.SubscribeRequest{
			Topic:    config.testTopicName,
			Metadata: config.subscribeMetadata,
		}, func(msg *pubsub.NewMessage) error {
			actualReadCount++
			return nil
		})
		assert.NoError(t, err, "expected no error on subscribe")
	})

	t.Run("Publish data to pubsub", func(t *testing.T) {
		for k := 0; k < config.messageCount; k++ {
			data := []byte("message-" + strconv.Itoa(k))
			err := store.Publish(&pubsub.PublishRequest{
				Data:       data,
				PubsubName: config.pubsubName,
				Topic:      config.testTopicName,
				Metadata:   config.publishMetadata,
			})
			assert.NoError(t, err, "expected no error on publishing data %s", data)
		}
	})

	t.Run("Read data from pubsub", func(t *testing.T) {
		t.Logf("waiting for %v ms to complete read", config.maxReadDurationInMs)
		time.Sleep(config.maxReadDurationInMs)
		assert.LessOrEqual(t, config.messageCount, actualReadCount, "expected to read %v messages", config.messageCount)
		// Properly close connection to store
		store.Close()
	})

}
