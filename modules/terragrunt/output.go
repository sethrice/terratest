package terragrunt

import (
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

// TODO: Add OutputAll/OutputAllE when terragrunt supports combined JSON output format.
// Currently, `output --all -json` returns separate JSON objects per module without module prefixes,
// making it impossible to reliably map outputs to their source modules.

// OutputAllJson runs terragrunt run --all output -json and returns the raw JSON string.
// Note: Current terragrunt versions return separate JSON objects per module, not a combined object.
func OutputAllJson(t testing.TestingT, options *Options) string {
	out, err := OutputAllJsonE(t, options)
	require.NoError(t, err)
	return out
}

// OutputAllJsonE runs terragrunt run --all output -json and returns the raw JSON string.
// Note: Current terragrunt versions return separate JSON objects per module, not a combined object.
func OutputAllJsonE(t testing.TestingT, options *Options) (string, error) {
	optsCopy := *options
	optsCopy.TerragruntArgs = append([]string{"--no-color"}, options.TerragruntArgs...)

	args := buildRunArgs([]string{"--all"}, []string{"output", "-json"})
	rawOutput, err := runTerragruntCommandE(t, &optsCopy, "run", args...)
	if err != nil {
		return "", err
	}

	// Extract only JSON content from output, filtering log lines and other terragrunt messages
	return extractJsonContent(rawOutput)
}

// OutputJson runs terragrunt run output -json for a single unit and returns clean JSON.
// If key is non-empty, returns the JSON value for that specific output.
// If key is empty, returns all outputs as JSON.
func OutputJson(t testing.TestingT, options *Options, key string) string {
	out, err := OutputJsonE(t, options, key)
	require.NoError(t, err)
	return out
}

// OutputJsonE runs terragrunt run output -json for a single unit and returns clean JSON.
// If key is non-empty, returns the JSON value for that specific output.
// If key is empty, returns all outputs as JSON.
func OutputJsonE(t testing.TestingT, options *Options, key string) (string, error) {
	optsCopy := *options
	optsCopy.TerragruntArgs = append([]string{"--no-color"}, options.TerragruntArgs...)

	tfArgs := []string{"-json"}
	if key != "" {
		tfArgs = append(tfArgs, key)
	}

	args := buildRunArgs([]string{}, append([]string{"output"}, tfArgs...))
	rawOutput, err := runTerragruntCommandE(t, &optsCopy, "run", args...)
	if err != nil {
		return "", err
	}

	return cleanTerragruntJson(rawOutput)
}
