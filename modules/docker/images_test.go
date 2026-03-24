package docker_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/docker"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/stretchr/testify/assert"
)

func TestListImagesAndDeleteImage(t *testing.T) {
	t.Parallel()

	uniqueID := strings.ToLower(random.UniqueID())
	repo := "gruntwork-io/test-image"
	tag := "v1-" + uniqueID
	img := fmt.Sprintf("%s:%s", repo, tag)

	options := &docker.BuildOptions{
		Tags: []string{img},
	}
	docker.Build(t, "../../test/fixtures/docker", options)

	assert.True(t, docker.DoesImageExist(t, img, nil))
	docker.DeleteImage(t, img, nil)
	assert.False(t, docker.DoesImageExist(t, img, nil))
}
