package terragrunt

import (
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

// PlanAllExitCode runs terragrunt run --all plan with the given options and returns the detailed exit code.
// This will fail the test if there is an error in the command.
func PlanAllExitCode(t testing.TestingT, options *Options) int {
	exitCode, err := PlanAllExitCodeE(t, options)
	require.NoError(t, err)
	return exitCode
}

// PlanAllExitCodeE runs terragrunt run --all -- plan with the given options and returns the detailed exit code.
func PlanAllExitCodeE(t testing.TestingT, options *Options) (int, error) {
	args := buildRunArgs([]string{"--all"}, "plan", []string{"-input=false", "-lock=true", "-detailed-exitcode"})
	return getExitCodeForTerragruntCommandE(t, options, append([]string{"run"}, args...)...)
}

// Plan runs terragrunt run plan for a single unit and returns stdout/stderr.
func Plan(t testing.TestingT, options *Options) string {
	out, err := PlanE(t, options)
	require.NoError(t, err)
	return out
}

// PlanE runs terragrunt run -- plan for a single unit and returns stdout/stderr.
// Uses -lock=false since plan is a read-only operation that does not need state locking.
func PlanE(t testing.TestingT, options *Options) (string, error) {
	args := buildRunArgs([]string{}, "plan", []string{"-input=false", "-lock=false"})
	return runTerragruntCommandE(t, options, "run", args...)
}

// PlanExitCode runs terragrunt run plan for a single unit and returns the detailed exit code.
// This will fail the test if there is an error in the command.
func PlanExitCode(t testing.TestingT, options *Options) int {
	exitCode, err := PlanExitCodeE(t, options)
	require.NoError(t, err)
	return exitCode
}

// PlanExitCodeE runs terragrunt run -- plan for a single unit and returns the detailed exit code.
func PlanExitCodeE(t testing.TestingT, options *Options) (int, error) {
	args := buildRunArgs([]string{}, "plan", []string{"-input=false", "-lock=true", "-detailed-exitcode"})
	return getExitCodeForTerragruntCommandE(t, options, append([]string{"run"}, args...)...)
}

// InitAndPlan runs terragrunt init followed by plan for a single unit and returns the plan stdout/stderr.
func InitAndPlan(t testing.TestingT, options *Options) string {
	out, err := InitAndPlanE(t, options)
	require.NoError(t, err)
	return out
}

// InitAndPlanE runs terragrunt init followed by plan for a single unit and returns the plan stdout/stderr.
func InitAndPlanE(t testing.TestingT, options *Options) (string, error) {
	if _, err := InitE(t, options); err != nil {
		return "", err
	}
	return PlanE(t, options)
}
