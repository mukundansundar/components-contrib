// +build conf

package secretstores

import (
	"os"
	"strings"
	"testing"

	"github.com/dapr/components-contrib/secretstores"
	local_env "github.com/dapr/components-contrib/secretstores/local/env"
	local_file "github.com/dapr/components-contrib/secretstores/local/file"
)

func TestLocalFileSecretStore(t *testing.T) {
	runWithSecretStore(t, "file", func() secretstores.SecretStore {
		return local_file.NewLocalSecretStore(nil)
	}, TestConfig{
		req: secretstores.GetSecretRequest{
			Name: "mysecret",
		},
		bulkReq: secretstores.BulkGetSecretRequest{
			Metadata: map[string]string{},
		},
		expectedResponse: map[string]string{
			"mysecret": "abcd",
		},
		expectedBulkResponse: map[string]string{
			"mysecret":     "abcd",
			"secondsecret": "efgh",
		},
	})
}

func TestLocalEnvSecretStore(t *testing.T) {
	os.Setenv("mysecret", "abcd")
	os.Setenv("secondsecret", "efgh")
	r := map[string]string{}

	for _, element := range os.Environ() {
		envVariable := strings.Split(element, "=")
		r[envVariable[0]] = envVariable[1]
	}
	runWithSecretStore(t, "env", func() secretstores.SecretStore {
		return local_env.NewEnvSecretStore(nil)
	}, TestConfig{
		req: secretstores.GetSecretRequest{
			Name: "mysecret",
		},
		bulkReq: secretstores.BulkGetSecretRequest{
			Metadata: map[string]string{},
		},
		expectedResponse: map[string]string{
			"mysecret": "abcd",
		},
		expectedBulkResponse: r,
	})
}
