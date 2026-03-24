package gcp

import (
	"context"
	"fmt"

	"github.com/google/go-containerregistry/pkg/authn"
	gcrname "github.com/google/go-containerregistry/pkg/name"
	gcrgoogle "github.com/google/go-containerregistry/pkg/v1/google"
	gcrremote "github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

// DeleteGCRRepo deletes a GCR repository including all tagged images.
func DeleteGCRRepo(t testing.TestingT, repo string) {
	err := DeleteGCRRepoE(t, repo)
	require.NoError(t, err)
}

// DeleteGCRRepoE deletes a GCR repository including all tagged images.
func DeleteGCRRepoE(t testing.TestingT, repo string) error {
	// create a new authenticator for the API calls
	authenticator, err := newGCRAuthenticator()
	if err != nil {
		return fmt.Errorf("failed to create authenticator: %w", err)
	}

	gcrrepo, err := gcrname.NewRepository(repo)
	if err != nil {
		return fmt.Errorf("failed to get repo: %w", err)
	}

	logger.Default.Logf(t, "Retrieving Image Digests %s", gcrrepo)

	tags, err := gcrgoogle.List(gcrrepo, gcrgoogle.WithAuth(authenticator))
	if err != nil {
		return fmt.Errorf("failed to list tags for repo %s: %w", repo, err)
	}

	// attempt to delete the latest image tag
	latestRef := repo + ":latest"
	logger.Default.Logf(t, "Deleting Image Ref %s", latestRef)

	if err := DeleteGCRImageRefE(t, latestRef); err != nil {
		return fmt.Errorf("failed to delete GCR image reference %s: %w", latestRef, err)
	}

	// delete image references sequentially
	for k := range tags.Manifests {
		ref := repo + "@" + k
		logger.Default.Logf(t, "Deleting Image Ref %s", ref)

		if err := DeleteGCRImageRefE(t, ref); err != nil {
			return fmt.Errorf("failed to delete GCR image reference %s: %w", ref, err)
		}
	}

	return nil
}

// DeleteGCRImageRef deletes a single repo image ref/digest.
func DeleteGCRImageRef(t testing.TestingT, ref string) {
	err := DeleteGCRImageRefE(t, ref)
	require.NoError(t, err)
}

// DeleteGCRImageRefE deletes a single repo image ref/digest.
func DeleteGCRImageRefE(t testing.TestingT, ref string) error {
	name, err := gcrname.ParseReference(ref)
	if err != nil {
		return fmt.Errorf("failed to parse reference %s: %w", ref, err)
	}

	// create a new authenticator for the API calls
	authenticator, err := newGCRAuthenticator()
	if err != nil {
		return fmt.Errorf("failed to create authenticator: %w", err)
	}

	opts := gcrremote.WithAuth(authenticator)

	if err := gcrremote.Delete(name, opts); err != nil {
		return fmt.Errorf("failed to delete %s: %w", name, err)
	}

	return nil
}

func newGCRAuthenticator() (authn.Authenticator, error) {
	if ts, ok := getStaticTokenSource(); ok {
		return gcrgoogle.NewTokenSourceAuthenticator(ts), nil
	}

	return gcrgoogle.NewEnvAuthenticator(context.Background())
}
