// SPDX-License-Identifier: AGPL-3.0-only

package azure

import (
	"testing"

	"github.com/go-kit/log"
	"github.com/grafana/dskit/flagext"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thanos-io/objstore/providers/azure"
	"gopkg.in/yaml.v2"
)

func TestNewBucketClient(t *testing.T) {
	t.Run("empty endpoint should be replaced with default", func(t *testing.T) {
		cfg := Config{
			StorageAccountName: "test",
			StorageAccountKey:  flagext.SecretWithValue("test"),
			ContainerName:      "test",
			MaxRetries:         3,
		}
		bkt, err := NewBucketClient(cfg, "test", log.NewNopLogger(), fakeFactory(t, cfg))
		require.NoError(t, err)
		require.NotNil(t, bkt)
	})

	t.Run("non-empty endpoint should be kept", func(t *testing.T) {
		cfg := Config{
			StorageAccountName: "test",
			StorageAccountKey:  flagext.SecretWithValue("test"),
			ContainerName:      "test",
			MaxRetries:         3,
			Endpoint:           "test-endpoint",
		}
		bkt, err := NewBucketClient(cfg, "test", log.NewNopLogger(), fakeFactory(t, cfg))
		require.NoError(t, err)
		require.NotNil(t, bkt)
	})
}

// fakeFactory is a test utility to act as an azure.Bucket factory, but in reality verify the input config.
func fakeFactory(t *testing.T, cfg Config) func(log.Logger, []byte, string) (*azure.Bucket, error) {
	expCfg := azure.DefaultConfig
	expCfg.StorageAccountName = cfg.StorageAccountName
	expCfg.StorageAccountKey = cfg.StorageAccountKey.String()
	expCfg.ContainerName = cfg.ContainerName
	expCfg.MaxRetries = cfg.MaxRetries
	expCfg.UserAssignedID = cfg.UserAssignedID
	if cfg.Endpoint != "" {
		expCfg.Endpoint = cfg.Endpoint
	}

	return func(_ log.Logger, c []byte, _ string) (*azure.Bucket, error) {
		t.Helper()

		var azCfg azure.Config
		require.NoError(t, yaml.Unmarshal(c, &azCfg))
		assert.Equal(t, expCfg, azCfg)

		return &azure.Bucket{}, nil
	}
}