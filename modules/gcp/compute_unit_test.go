package gcp_test

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/gcp"
	"github.com/stretchr/testify/assert"
	"google.golang.org/api/compute/v1"
)

// TestNewMetadataPreservesExisting is a regression test for issue #1655.
// Verifies that existing metadata is preserved when adding new key-value pairs.
func TestNewMetadataPreservesExisting(t *testing.T) {
	t.Parallel()

	existingVal := "existing-value"
	oldMetadata := &compute.Metadata{
		Fingerprint: "test-fingerprint",
		Items:       []*compute.MetadataItems{{Key: "existing-key", Value: &existingVal}},
	}

	result := gcp.NewMetadata(oldMetadata, map[string]string{"new-key": "new-value"})

	// Convert to map for easier assertion
	got := make(map[string]string)
	for _, item := range result.Items {
		got[item.Key] = *item.Value
	}

	assert.Equal(t, "test-fingerprint", result.Fingerprint)
	assert.Equal(t, "existing-value", got["existing-key"], "existing metadata should be preserved")
	assert.Equal(t, "new-value", got["new-key"], "new metadata should be added")
}
