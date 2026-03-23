package opa_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/gruntwork-io/terratest/modules/files"
	"github.com/gruntwork-io/terratest/modules/git"
	"github.com/gruntwork-io/terratest/modules/opa"
)

// TestDownloadPolicyReturnsLocalPath makes sure the DownloadPolicyE function returns a local path without processing it.
func TestDownloadPolicyReturnsLocalPath(t *testing.T) {
	t.Parallel()

	localPath := "../../examples/terraform-opa-example/policy/enforce_source.rego"
	path, err := opa.DownloadPolicyE(t, localPath)
	require.NoError(t, err)
	assert.Equal(t, localPath, path)
}

// TestDownloadPolicyDownloadsRemote makes sure the DownloadPolicyE function returns a remote path to a temporary
// directory.
func TestDownloadPolicyDownloadsRemote(t *testing.T) {
	t.Parallel()

	curRef := git.GetCurrentGitRefContext(t, t.Context(), "")
	baseDir := "git::https://github.com/gruntwork-io/terratest.git?ref=" + curRef
	localPath := "../../examples/terraform-opa-example/policy/enforce_source.rego"
	remotePath := "git::https://github.com/gruntwork-io/terratest.git//examples/terraform-opa-example/policy/enforce_source.rego?ref=" + curRef

	// Make sure we clean up the downloaded file, while simultaneously asserting that the download dir was stored in the
	// cache.
	defer func() {
		downloadPathRaw, inCache := opa.PolicyDirCache.Load(baseDir)
		require.True(t, inCache)

		downloadPath := downloadPathRaw.(string)

		if strings.HasSuffix(downloadPath, "/getter") {
			downloadPath = filepath.Dir(downloadPath)
		}

		assert.NoError(t, os.RemoveAll(downloadPath))
	}()

	path, err := opa.DownloadPolicyE(t, remotePath)
	require.NoError(t, err)

	absPath, err := filepath.Abs(localPath)
	require.NoError(t, err)
	assert.NotEqual(t, absPath, path)

	localContents, err := os.ReadFile(localPath)
	require.NoError(t, err)

	remoteContents, err := os.ReadFile(path)
	require.NoError(t, err)
	assert.Equal(t, localContents, remoteContents)
}

// TestDownloadPolicyReusesCachedDir makes sure the DownloadPolicyE function uses the cache if it has already downloaded
// an existing base path.
func TestDownloadPolicyReusesCachedDir(t *testing.T) {
	t.Parallel()

	baseDir := "git::https://github.com/gruntwork-io/terratest.git?ref=main"
	remotePath := "git::https://github.com/gruntwork-io/terratest.git//examples/terraform-opa-example/policy/enforce_source.rego?ref=main"
	remotePathAltSubPath := "git::https://github.com/gruntwork-io/terratest.git//modules/opa/eval.go?ref=main"

	// Make sure we clean up the downloaded file, while simultaneously asserting that the download dir was stored in the
	// cache.
	defer func() {
		downloadPathRaw, inCache := opa.PolicyDirCache.Load(baseDir)
		require.True(t, inCache)

		downloadPath := downloadPathRaw.(string)

		if strings.HasSuffix(downloadPath, "/getter") {
			downloadPath = filepath.Dir(downloadPath)
		}

		assert.NoError(t, os.RemoveAll(downloadPath))
	}()

	path, err := opa.DownloadPolicyE(t, remotePath)
	require.NoError(t, err)
	files.FileExists(path)

	downloadPathRaw, inCache := opa.PolicyDirCache.Load(baseDir)
	require.True(t, inCache)

	downloadPath := downloadPathRaw.(string)

	// make sure the second call is exactly equal to the first call
	newPath, err := opa.DownloadPolicyE(t, remotePath)
	require.NoError(t, err)
	assert.Equal(t, path, newPath)

	// Also make sure the cache is reused for alternative sub dirs.
	newAltPath, err := opa.DownloadPolicyE(t, remotePathAltSubPath)
	require.NoError(t, err)
	assert.True(t, strings.HasPrefix(path, downloadPath))
	assert.True(t, strings.HasPrefix(newAltPath, downloadPath))
}
