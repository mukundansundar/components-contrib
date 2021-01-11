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
		// inital init time and set time are higher than follow-up requests
		maxInitDurationInMs:   time.Duration(20) * time.Millisecond,
		maxSetDurationInMs:    time.Duration(20) * time.Millisecond,
		maxDeleteDurationInMs: time.Duration(10) * time.Millisecond,
		maxGetDurationInMs:    time.Duration(10) * time.Millisecond,
		numBulkRequests:       10,
	})
}
