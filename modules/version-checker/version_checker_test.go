package version_checker_test //nolint:staticcheck // package name matches directory convention

import (
	"errors"
	"testing"

	checker "github.com/gruntwork-io/terratest/modules/version-checker"
	"github.com/stretchr/testify/require"
)

func TestParamValidation(t *testing.T) {
	t.Parallel()

	t.Run("Empty Params", func(t *testing.T) {
		t.Parallel()
		err := checker.CheckVersionContextE(t, t.Context(), checker.CheckVersionParams{})

		var missingErr *checker.MissingParamErr

		require.ErrorAs(t, err, &missingErr)
	})

	t.Run("Missing VersionConstraint", func(t *testing.T) {
		t.Parallel()
		err := checker.CheckVersionContextE(t, t.Context(), checker.CheckVersionParams{
			Binary:            checker.Docker,
			VersionConstraint: "",
			WorkingDir:        ".",
		})

		var missingErr *checker.MissingParamErr

		require.ErrorAs(t, err, &missingErr)
	})

	t.Run("Invalid Version Constraint Format", func(t *testing.T) {
		t.Parallel()
		err := checker.CheckVersionContextE(t, t.Context(), checker.CheckVersionParams{
			Binary:            checker.Docker,
			VersionConstraint: "abc",
			WorkingDir:        ".",
		})

		var constraintErr *checker.InvalidVersionConstraintErr

		require.ErrorAs(t, err, &constraintErr)
	})

	t.Run("Valid Params", func(t *testing.T) {
		t.Parallel()
		err := checker.CheckVersionContextE(t, t.Context(), checker.CheckVersionParams{
			Binary:            checker.Docker,
			VersionConstraint: ">= 0.0.1",
			WorkingDir:        ".",
		})
		require.NoError(t, err)
	})
}

func TestCheckVersionConstraintMismatch(t *testing.T) {
	t.Parallel()

	err := checker.CheckVersionContextE(t, t.Context(), checker.CheckVersionParams{
		Binary:            checker.Docker,
		VersionConstraint: ">= 999.999.999",
		WorkingDir:        ".",
	})

	var mismatchErr *checker.VersionMismatchErr

	require.ErrorAs(t, err, &mismatchErr)
	require.Equal(t, ">= 999.999.999", mismatchErr.Constraint)
}

func TestMissingParamErrFields(t *testing.T) {
	t.Parallel()

	err := checker.CheckVersionContextE(t, t.Context(), checker.CheckVersionParams{})

	var missingErr *checker.MissingParamErr

	require.ErrorAs(t, err, &missingErr)
	require.Equal(t, "WorkingDir", missingErr.Param)
}

func TestInvalidVersionConstraintErrUnwrap(t *testing.T) {
	t.Parallel()

	err := checker.CheckVersionContextE(t, t.Context(), checker.CheckVersionParams{
		Binary:            checker.Docker,
		VersionConstraint: "abc",
		WorkingDir:        ".",
	})

	var constraintErr *checker.InvalidVersionConstraintErr

	require.ErrorAs(t, err, &constraintErr)
	require.Error(t, errors.Unwrap(constraintErr), "underlying parse error should be wrapped")
}

// TestCheckVersionEndToEnd assumes Docker, Terraform, and Packer are installed
// and their versions are greater than 0.0.1.
func TestCheckVersionEndToEnd(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		param checker.CheckVersionParams
	}{
		{
			name: "Docker",
			param: checker.CheckVersionParams{
				Binary:            checker.Docker,
				VersionConstraint: ">= 0.0.1",
				WorkingDir:        ".",
			},
		},
		{
			name: "Terraform",
			param: checker.CheckVersionParams{
				Binary:            checker.Terraform,
				VersionConstraint: ">= 0.0.1",
				WorkingDir:        ".",
			},
		},
		{
			name: "Packer",
			param: checker.CheckVersionParams{
				BinaryPath:        "/usr/local/bin/packer",
				Binary:            checker.Packer,
				VersionConstraint: ">= 0.0.1",
				WorkingDir:        ".",
			},
		},
	}

	for _, tc := range tests {
		err := checker.CheckVersionContextE(t, t.Context(), tc.param)
		require.NoError(t, err, tc.name)
	}
}
