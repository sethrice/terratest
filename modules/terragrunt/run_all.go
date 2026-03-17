package terragrunt

import (
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

// RunAll runs terragrunt run --all [tgArgs...] -- <command> [tfArgs...] with the given options and returns stdout/stderr.
// This is a generic wrapper that allows running any OpenTofu/Terraform command with --all flag.
// The -- separator disambiguates Terragrunt flags from OpenTofu/Terraform flags.
func RunAll(t testing.TestingT, options *Options, command string, tgArgs []string, tfArgs []string) string {
	out, err := RunAllE(t, options, command, tgArgs, tfArgs)
	require.NoError(t, err)
	return out
}

// RunAllE runs terragrunt run --all [tgArgs...] -- <command> [tfArgs...] with the given options and returns stdout/stderr.
// This is a generic wrapper that allows running any OpenTofu/Terraform command with --all flag.
// The -- separator disambiguates Terragrunt flags from OpenTofu/Terraform flags.
func RunAllE(t testing.TestingT, options *Options, command string, tgArgs []string, tfArgs []string) (string, error) {
	args := buildRunArgs(append([]string{"--all"}, tgArgs...), command, tfArgs)
	return runTerragruntCommandE(t, options, "run", args...)
}
