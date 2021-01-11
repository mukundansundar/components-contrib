// +build conf

package pubsub

import (
	"testing"

	"github.com/dapr/components-contrib/pubsub"
	"github.com/dapr/components-contrib/pubsub/redis"
	"github.com/dapr/dapr/pkg/logger"
)

func TestRedis(t *testing.T) {
	runWithPubSub(t, "redis", func() pubsub.PubSub {
		return redis.NewRedisStreams(logger.NewLogger("testLogger"))
	}, TestConfig{
		maxReadDurationInMs: 2000, // wait for 5 seconds
		pubsubName:          "pubsub",
		testTopicName:       "testtopic",
		messageCount:        10,
	})
}
