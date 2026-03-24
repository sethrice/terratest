package docker_test

import (
	"context"
	"crypto/tls"
	"fmt"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/docker"
	http_helper "github.com/gruntwork-io/terratest/modules/http-helper"
	"github.com/gruntwork-io/terratest/modules/shell"
	"github.com/stretchr/testify/require"
)

func TestStop(t *testing.T) {
	t.Parallel()

	// appending timestamp to container name to run tests in parallel
	name := "test-nginx" + strconv.FormatInt(time.Now().UnixNano(), 10)

	// choosing a unique port since 80 may not fly well on test machines
	port := "13030"
	host := docker.GetDockerHost()

	testURL := fmt.Sprintf("http://%s:%s", host, port)

	// for testing the stopping of a docker container
	// we got to run a container first and then stop it
	runOpts := &docker.RunOptions{
		Detach:       true,
		Name:         name,
		Remove:       true,
		OtherOptions: []string{"-p", port + ":80"},
	}
	docker.Run(t, "nginx:1.17-alpine", runOpts)

	// verify nginx is running
	http_helper.HttpGetWithRetryWithCustomValidation(t, testURL, &tls.Config{}, 60, 2*time.Second, verifyNginxIsUp)

	// try to stop it now
	out := docker.Stop(t, []string{name}, &docker.StopOptions{})
	require.Contains(t, out, name)

	// verify nginx is down
	// run a docker ps with name filter
	command := &shell.Command{
		Command: "docker",
		Args:    []string{"ps", "-q", "--filter", "name=" + name},
	}
	output := shell.RunCommandContextAndGetStdOut(t, context.Background(), command)
	require.Empty(t, output)
}

func verifyNginxIsUp(statusCode int, body string) bool {
	return statusCode == 200 && strings.Contains(body, "nginx!")
}
