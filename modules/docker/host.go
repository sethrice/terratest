package docker

import (
	"os"
	"strings"
)

// minDockerURLParts is the minimum number of colon-separated parts in a valid
// DOCKER_HOST URL (e.g. "tcp://host:port" splits into at least 2 parts).
const minDockerURLParts = 2

// GetDockerHost returns the name or address of the host on which the Docker engine is running.
func GetDockerHost() string {
	return GetDockerHostFromEnv(os.Environ())
}

// GetDockerHostFromEnv parses the DOCKER_HOST value from the given environment
// variable list and returns the host portion. If DOCKER_HOST is absent, empty,
// or uses a local transport (unix/npipe), it returns "localhost".
func GetDockerHostFromEnv(env []string) string {
	// Parses the DOCKER_HOST environment variable to find the address
	//
	// For valid formats see:
	// https://github.com/docker/cli/blob/6916b427a0b07e8581d121967633235ced6db9a1/opts/hosts.go#L69
	var dockerURL []string

	for _, item := range env {
		envVar := strings.Split(item, "=")
		if len(envVar) == 2 && envVar[0] == "DOCKER_HOST" {
			dockerURL = strings.Split(envVar[1], ":")
			break
		}
	}

	if len(dockerURL) < minDockerURLParts {
		// DOCKER_HOST was empty, not present or not a valid URL
		return "localhost"
	}

	switch dockerURL[0] {
	case "tcp", "ssh", "fd":
		return strings.TrimPrefix(dockerURL[1], "//")
	default:
		// if DOCKER_HOST is not in one of the formats listed above, return default
		return "localhost"
	}
}
