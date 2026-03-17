package terragrunt

import (
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

// Run runs terragrunt run [tgArgs...] -- <command> [tfArgs...] with the given options and returns stdout/stderr.
// This is a generic wrapper that allows running any tofu/terraform command through terragrunt run.
// The -- separator disambiguates Terragrunt flags from OpenTofu/Terraform flags.
func Run(t testing.TestingT, options *Options, command string, tgArgs []string, tfArgs []string) string {
	out, err := RunE(t, options, command, tgArgs, tfArgs)
	require.NoError(t, err)
	return out
}

// RunE runs terragrunt run [tgArgs...] -- <command> [tfArgs...] with the given options and returns stdout/stderr.
// This is a generic wrapper that allows running any tofu/terraform command through terragrunt run.
// The -- separator disambiguates Terragrunt flags from OpenTofu/Terraform flags.
func RunE(t testing.TestingT, options *Options, command string, tgArgs []string, tfArgs []string) (string, error) {
	args := buildRunArgs(tgArgs, command, tfArgs)
	return runTerragruntCommandE(t, options, "run", args...)
}
