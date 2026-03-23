package version_checker_test //nolint:staticcheck // package name matches directory convention

import (
	"testing"

	checker "github.com/gruntwork-io/terratest/modules/version-checker"
	"github.com/stretchr/testify/require"
)

func TestParamValidation(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name                 string
		expectedErrorMessage string
		param                checker.CheckVersionParams
		containError         bool
	}{
		{
			name:                 "Empty Params",
			param:                checker.CheckVersionParams{},
			containError:         true,
			expectedErrorMessage: "set WorkingDir in params",
		},
		{
			name: "Missing VersionConstraint",
			param: checker.CheckVersionParams{
				Binary:            checker.Docker,
				VersionConstraint: "",
				WorkingDir:        ".",
			},
			containError:         true,
			expectedErrorMessage: "set VersionConstraint in params",
		},
		{
			name: "Invalid Version Constraint Format",
			param: checker.CheckVersionParams{
				Binary:            checker.Docker,
				VersionConstraint: "abc",
				WorkingDir:        ".",
			},
			containError:         true,
			expectedErrorMessage: "invalid version constraint format found {abc}",
		},
		{
			name: "Valid Params",
			param: checker.CheckVersionParams{
				Binary:            checker.Docker,
				VersionConstraint: ">= 0.0.1",
				WorkingDir:        ".",
			},
			containError: false,
		},
	}

	for _, tc := range tests {
		err := checker.CheckVersionE(t, tc.param)
		if tc.containError {
			require.EqualError(t, err, tc.expectedErrorMessage, tc.name)
		} else {
			require.NoError(t, err, tc.name)
		}
	}
}

func TestCheckVersionConstraintMismatch(t *testing.T) {
	t.Parallel()

	err := checker.CheckVersionE(t, checker.CheckVersionParams{
		Binary:            checker.Docker,
		VersionConstraint: ">= 999.999.999",
		WorkingDir:        ".",
	})
	require.Error(t, err, "expected version mismatch error")
	require.Contains(t, err.Error(), "failed the version constraint")
}

// TestCheckVersionEndToEnd assumes Docker, Terraform, and Packer are installed
// and their versions are greater than 0.0.1.
func TestCheckVersionEndToEnd(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		param checker.CheckVersionParams
	}{
		{name: "Docker", param: checker.CheckVersionParams{
			Binary:            checker.Docker,
			VersionConstraint: ">= 0.0.1",
			WorkingDir:        ".",
		}},
		{name: "Terraform", param: checker.CheckVersionParams{
			Binary:            checker.Terraform,
			VersionConstraint: ">= 0.0.1",
			WorkingDir:        ".",
		}},
		{name: "Packer", param: checker.CheckVersionParams{
			BinaryPath:        "/usr/local/bin/packer",
			Binary:            checker.Packer,
			VersionConstraint: ">= 0.0.1",
			WorkingDir:        ".",
		}},
	}

	for _, tc := range tests {
		err := checker.CheckVersionE(t, tc.param)
		require.NoError(t, err, tc.name)
	}
}
