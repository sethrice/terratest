package terragrunt

import (
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

// RunAll runs terragrunt run --all <command> with the given options and returns stdout/stderr.
// This is a generic wrapper that allows running any terraform command with --all flag.
func RunAll(t testing.TestingT, options *Options, command string, additionalArgs ...string) string {
	out, err := RunAllE(t, options, command, additionalArgs...)
	require.NoError(t, err)
	return out
}

// RunAllE runs terragrunt run --all <command> with the given options and returns stdout/stderr.
// This is a generic wrapper that allows running any terraform command with --all flag.
func RunAllE(t testing.TestingT, options *Options, command string, additionalArgs ...string) (string, error) {
	args := append([]string{"--all", command}, additionalArgs...)
	return runTerragruntCommandE(t, options, "run", args...)
}
