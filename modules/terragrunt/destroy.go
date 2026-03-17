package terragrunt

import (
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

// DestroyAll runs terragrunt run --all destroy with the given options and returns stdout.
func DestroyAll(t testing.TestingT, options *Options) string {
	out, err := DestroyAllE(t, options)
	require.NoError(t, err)
	return out
}

// DestroyAllE runs terragrunt run --all -- destroy with the given options and returns stdout.
func DestroyAllE(t testing.TestingT, options *Options) (string, error) {
	args := buildRunArgs([]string{"--all"}, "destroy", []string{"-auto-approve", "-input=false"})
	return runTerragruntCommandE(t, options, "run", args...)
}

// Destroy runs terragrunt run destroy for a single unit and returns stdout/stderr.
func Destroy(t testing.TestingT, options *Options) string {
	out, err := DestroyE(t, options)
	require.NoError(t, err)
	return out
}

// DestroyE runs terragrunt run -- destroy for a single unit and returns stdout/stderr.
func DestroyE(t testing.TestingT, options *Options) (string, error) {
	args := buildRunArgs([]string{}, "destroy", []string{"-auto-approve", "-input=false"})
	return runTerragruntCommandE(t, options, "run", args...)
}
