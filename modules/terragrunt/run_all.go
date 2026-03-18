package terragrunt

import (
	"github.com/gruntwork-io/terratest/modules/testing"
)

// Deprecated: Use Run with the --all flag in tgArgs instead.
// RunAll runs terragrunt run --all -- <command> with the given options and returns stdout/stderr.
func RunAll(t testing.TestingT, options *Options, command string) string {
	return Run(t, options, []string{"--all"}, []string{command})
}

// Deprecated: Use RunE with the --all flag in tgArgs instead.
// RunAllE runs terragrunt run --all -- <command> with the given options and returns stdout/stderr.
func RunAllE(t testing.TestingT, options *Options, command string) (string, error) {
	return RunE(t, options, []string{"--all"}, []string{command})
}
