package terragrunt

import (
	"fmt"

	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

// Run runs terragrunt run [tgArgs...] -- [tfArgs...] with the given options and returns stdout/stderr.
// This is a generic wrapper that allows running any OpenTofu/Terraform command through terragrunt run.
// The -- separator disambiguates Terragrunt flags from OpenTofu/Terraform flags.
// The OpenTofu/Terraform command (e.g. "apply") should be the first element of tfArgs.
func Run(t testing.TestingT, options *Options, tgArgs []string, tfArgs []string) string {
	out, err := RunE(t, options, tgArgs, tfArgs)
	require.NoError(t, err)
	return out
}

// RunE runs terragrunt run [tgArgs...] -- [tfArgs...] with the given options and returns stdout/stderr.
// This is a generic wrapper that allows running any OpenTofu/Terraform command through terragrunt run.
// The -- separator disambiguates Terragrunt flags from OpenTofu/Terraform flags.
// The OpenTofu/Terraform command (e.g. "apply") should be the first element of tfArgs.
func RunE(t testing.TestingT, options *Options, tgArgs []string, tfArgs []string) (string, error) {
	if len(tfArgs) == 0 {
		return "", fmt.Errorf("tfArgs cannot be empty; at minimum, an OpenTofu/Terraform command (e.g. \"apply\") is required")
	}
	args := buildRunArgs(tgArgs, tfArgs)
	return runTerragruntCommandE(t, options, "run", args...)
}
