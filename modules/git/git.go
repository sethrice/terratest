// Package git allows to interact with Git.
package git

import (
	"context"
	"os"
	"os/exec"
	"strings"

	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

// GetCurrentBranchName retrieves the current branch name or
// empty string in case of detached state.
//
// Deprecated: Use GetCurrentBranchNameContext instead.
func GetCurrentBranchName(t testing.TestingT) string {
	return GetCurrentBranchNameContext(t, context.Background())
}

// GetCurrentBranchNameContext is like GetCurrentBranchName but includes a context.
func GetCurrentBranchNameContext(t testing.TestingT, ctx context.Context) string {
	out, err := GetCurrentBranchNameContextE(t, ctx)
	if err != nil {
		t.Fatal(err)
	}

	return out
}

// GetCurrentBranchNameE retrieves the current branch name or
// empty string in case of detached state.
// Uses branch --show-current, which was introduced in git v2.22.
// Falls back to rev-parse for users of the older version, like Ubuntu 18.04.
//
// Deprecated: Use GetCurrentBranchNameContextE instead.
func GetCurrentBranchNameE(t testing.TestingT) (string, error) {
	return GetCurrentBranchNameContextE(t, context.Background())
}

// GetCurrentBranchNameContextE is like GetCurrentBranchNameE but includes a context.
func GetCurrentBranchNameContextE(t testing.TestingT, ctx context.Context) (string, error) {
	cmd := exec.CommandContext(ctx, "git", "branch", "--show-current")

	bytes, err := cmd.Output()
	if err != nil {
		return GetCurrentBranchNameOldContextE(t, ctx)
	}

	name := strings.TrimSpace(string(bytes))
	if name == "HEAD" {
		return "", nil
	}

	return name, nil
}

// GetCurrentBranchNameOldE retrieves the current branch name or
// empty string in case of detached state. This uses the older pattern
// of `git rev-parse` rather than `git branch --show-current`.
//
// Deprecated: Use GetCurrentBranchNameOldContextE instead.
func GetCurrentBranchNameOldE(t testing.TestingT) (string, error) {
	return GetCurrentBranchNameOldContextE(t, context.Background())
}

// GetCurrentBranchNameOldContextE is like GetCurrentBranchNameOldE but includes a context.
func GetCurrentBranchNameOldContextE(t testing.TestingT, ctx context.Context) (string, error) {
	cmd := exec.CommandContext(ctx, "git", "rev-parse", "--abbrev-ref", "HEAD")

	bytes, err := cmd.Output()
	if err != nil {
		return "", err
	}

	name := strings.TrimSpace(string(bytes))
	if name == "HEAD" {
		return "", nil
	}

	return name, nil
}

// GetCurrentGitRef retrieves current branch name, lightweight (non-annotated) tag or
// if tag points to the commit exact tag value.
//
// Deprecated: Use GetCurrentGitRefContext instead.
func GetCurrentGitRef(t testing.TestingT) string {
	return GetCurrentGitRefContext(t, context.Background())
}

// GetCurrentGitRefContext is like GetCurrentGitRef but includes a context.
func GetCurrentGitRefContext(t testing.TestingT, ctx context.Context) string {
	out, err := GetCurrentGitRefContextE(t, ctx)
	if err != nil {
		t.Fatal(err)
	}

	return out
}

// GetCurrentGitRefE retrieves current branch name, lightweight (non-annotated) tag or
// if tag points to the commit exact tag value.
//
// Deprecated: Use GetCurrentGitRefContextE instead.
func GetCurrentGitRefE(t testing.TestingT) (string, error) {
	return GetCurrentGitRefContextE(t, context.Background())
}

// GetCurrentGitRefContextE is like GetCurrentGitRefE but includes a context.
func GetCurrentGitRefContextE(t testing.TestingT, ctx context.Context) (string, error) {
	out, err := GetCurrentBranchNameContextE(t, ctx)
	if err != nil {
		return "", err
	}

	if out != "" {
		return out, nil
	}

	out, err = GetTagContextE(t, ctx)
	if err != nil {
		return "", err
	}

	return out, nil
}

// GetTagE retrieves lightweight (non-annotated) tag or if tag points
// to the commit exact tag value.
//
// Deprecated: Use GetTagContextE instead.
func GetTagE(t testing.TestingT) (string, error) {
	return GetTagContextE(t, context.Background())
}

// GetTagContextE is like GetTagE but includes a context.
func GetTagContextE(t testing.TestingT, ctx context.Context) (string, error) {
	cmd := exec.CommandContext(ctx, "git", "describe", "--tags")

	bytes, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(bytes)), nil
}

// GetRepoRoot retrieves the path to the root directory of the repo. This fails the test if there is an error.
//
// Deprecated: Use GetRepoRootContext instead.
func GetRepoRoot(t testing.TestingT) string {
	return GetRepoRootContext(t, context.Background())
}

// GetRepoRootContext is like GetRepoRoot but includes a context.
func GetRepoRootContext(t testing.TestingT, ctx context.Context) string {
	out, err := GetRepoRootContextE(t, ctx)
	require.NoError(t, err)

	return out
}

// GetRepoRootE retrieves the path to the root directory of the repo.
//
// Deprecated: Use GetRepoRootContextE instead.
func GetRepoRootE(t testing.TestingT) (string, error) {
	return GetRepoRootContextE(t, context.Background())
}

// GetRepoRootContextE is like GetRepoRootE but includes a context.
func GetRepoRootContextE(t testing.TestingT, ctx context.Context) (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	return GetRepoRootForDirContextE(t, ctx, dir)
}

// GetRepoRootForDir retrieves the path to the root directory of the repo in which dir resides
//
// Deprecated: Use GetRepoRootForDirContext instead.
func GetRepoRootForDir(t testing.TestingT, dir string) string {
	return GetRepoRootForDirContext(t, context.Background(), dir)
}

// GetRepoRootForDirContext is like GetRepoRootForDir but includes a context.
func GetRepoRootForDirContext(t testing.TestingT, ctx context.Context, dir string) string {
	out, err := GetRepoRootForDirContextE(t, ctx, dir)
	require.NoError(t, err)

	return out
}

// GetRepoRootForDirE retrieves the path to the root directory of the repo in which dir resides
//
// Deprecated: Use GetRepoRootForDirContextE instead.
func GetRepoRootForDirE(t testing.TestingT, dir string) (string, error) {
	return GetRepoRootForDirContextE(t, context.Background(), dir)
}

// GetRepoRootForDirContextE is like GetRepoRootForDirE but includes a context.
func GetRepoRootForDirContextE(t testing.TestingT, ctx context.Context, dir string) (string, error) {
	cmd := exec.CommandContext(ctx, "git", "rev-parse", "--show-toplevel")
	cmd.Dir = dir

	bytes, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(bytes)), nil
}
