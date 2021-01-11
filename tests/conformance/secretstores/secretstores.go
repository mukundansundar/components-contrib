// +build conf

package secretstores

import (
	"fmt"
	"testing"

	"github.com/dapr/components-contrib/secretstores"
	"github.com/dapr/components-contrib/tests/conformance"
	"github.com/stretchr/testify/assert"
)

type TestConfig struct {
	req                  secretstores.GetSecretRequest
	bulkReq              secretstores.BulkGetSecretRequest
	expectedResponse     map[string]string
	expectedBulkResponse map[string]string
}

func runWithSecretStore(t *testing.T, name string, componentFactory func() secretstores.SecretStore, config TestConfig) {
	store := componentFactory()
	comps, err := conformance.LoadComponents(fmt.Sprintf("../../config/secretstores/%s", name))
	assert.Nil(t, err)
	assert.Equal(t, len(comps), 1) // We only expect a single component per state store

	c := comps[0]
	props := conformance.ConvertMetadataToProperties(c.Spec.Metadata)
	// Run the state store conformance tests
	secretStoreConformanceTests(t, props, store, config)
}

func secretStoreConformanceTests(t *testing.T, props map[string]string, store secretstores.SecretStore, config TestConfig) {
	t.Run("Init secret store", func(t *testing.T) {
		err := store.Init(secretstores.Metadata{
			Properties: props,
		})
		assert.NoError(t, err, "expected no error on getting secret %v", config.req)
	})

	t.Run("Get secret", func(t *testing.T) {
		resp, err := store.GetSecret(config.req)
		assert.NoError(t, err, "expected no error on getting secret %v", config.req)
		assert.NotNil(t, resp, "expected value to be returned")
		assert.NotNil(t, resp.Data, "expected value to be returned")
		assert.Equal(t, resp.Data, config.expectedResponse, "expected values to be equal")
	})

	t.Run("Get bulk secret", func(t *testing.T) {
		resp, err := store.BulkGetSecret(config.bulkReq)
		assert.NoError(t, err, "expected no error on getting secret %v", config.req)
		assert.NotNil(t, resp, "expected value to be returned")
		assert.NotNil(t, resp.Data, "expected value to be returned")
		assert.Equal(t, resp.Data, config.expectedBulkResponse, "expected values to be equal")
	})
}
