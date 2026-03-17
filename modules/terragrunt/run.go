package terragrunt

import (
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

// Run runs terragrunt run <command> with the given options and returns stdout/stderr.
// This is a generic wrapper that allows running any tofu/terraform command through terragrunt run.
func Run(t testing.TestingT, options *Options, command string, additionalArgs ...string) string {
	out, err := RunE(t, options, command, additionalArgs...)
	require.NoError(t, err)
	return out
}

// RunE runs terragrunt run <command> with the given options and returns stdout/stderr.
// This is a generic wrapper that allows running any tofu/terraform command through terragrunt run.
func RunE(t testing.TestingT, options *Options, command string, additionalArgs ...string) (string, error) {
	args := append([]string{command}, additionalArgs...)
	return runTerragruntCommandE(t, options, "run", args...)
}
