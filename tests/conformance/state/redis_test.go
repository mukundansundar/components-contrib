// +build conf

package state

import (
	"testing"
	"time"

	"github.com/dapr/components-contrib/state"
	"github.com/dapr/components-contrib/state/redis"
)

func TestRedis(t *testing.T) {
	runWithStateStore(t, "redis", func() state.Store {
		return redis.NewRedisStateStore(nil)
	}, TestConfig{
		maxInitDurationInMs:   time.Duration(10) * time.Millisecond,
		maxSetDurationInMs:    time.Duration(10) * time.Millisecond,
		maxDeleteDurationInMs: time.Duration(10) * time.Millisecond,
		maxGetDurationInMs:    time.Duration(10) * time.Millisecond,
		numBulkRequests:       10,
	})
}
