package docker

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"slices"
	"strings"

	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/shell"
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

// Image represents a docker image, and exports all the fields that the docker images command returns for the
// image.
type Image struct {
	// ID is the image ID in docker, and can be used to identify the image in place of the repo and tag.
	ID string `json:"ID"`

	// Repository is the image repository.
	Repository string `json:"Repository"`

	// Tag is the image tag wichin the repository.
	Tag string `json:"Tag"`

	// CreatedAt represents a timestamp for when the image was created.
	CreatedAt string `json:"CreatedAt"`

	// CreatedSince is a diff between when the image was created to now.
	CreatedSince string `json:"CreatedSince"`

	// SharedSize is the amount of space that an image shares with another one (i.e. their common data).
	SharedSize string `json:"SharedSize"`

	// UniqueSize is the amount of space that is only used by a given image.
	UniqueSize string `json:"UniqueSize"`

	// VirtualSize is the total size of the image, combining SharedSize and UniqueSize.
	VirtualSize string `json:"VirtualSize"`

	// Containers represents the list of containers that are using the image.
	Containers string `json:"Containers"`

	// Digest is the hash digest of the image, if requested.
	Digest string `json:"Digest"`
}

// String returns the image reference in "repository:tag" format.
func (image Image) String() string {
	return fmt.Sprintf("%s:%s", image.Repository, image.Tag)
}

// DeleteImage removes a docker image using the Docker CLI. This will fail the test if there is an error.
func DeleteImage(t testing.TestingT, img string, logger *logger.Logger) {
	require.NoError(t, DeleteImageE(t, img, logger))
}

// DeleteImageE removes a docker image using the Docker CLI.
func DeleteImageE(t testing.TestingT, img string, logger *logger.Logger) error {
	cmd := &shell.Command{
		Command: "docker",
		Args:    []string{"rmi", img},
		Logger:  logger,
	}

	return shell.RunCommandContextE(t, context.Background(), cmd)
}

// ListImages calls docker images using the Docker CLI to list the available images on the local docker daemon.
func ListImages(t testing.TestingT, logger *logger.Logger) []Image {
	out, err := ListImagesE(t, logger)
	require.NoError(t, err)

	return out
}

// ListImagesE calls docker images using the Docker CLI to list the available images on the local docker daemon.
func ListImagesE(t testing.TestingT, logger *logger.Logger) ([]Image, error) {
	cmd := &shell.Command{
		Command: "docker",
		Args:    []string{"images", "--format", "{{ json . }}"},
		Logger:  logger,
	}

	out, err := shell.RunCommandContextAndGetOutputE(t, context.Background(), cmd)
	if err != nil {
		return nil, err
	}

	// Parse and return the list of image objects.
	images := []Image{}

	scanner := bufio.NewScanner(strings.NewReader(out))
	for scanner.Scan() {
		line := scanner.Text()

		var image Image

		err := json.Unmarshal([]byte(line), &image)
		if err != nil {
			return nil, err
		}

		images = append(images, image)
	}

	return images, nil
}

// DoesImageExist lists the images in the docker daemon and returns true if the given image label (repo:tag) exists.
// This will fail the test if there is an error.
func DoesImageExist(t testing.TestingT, imgLabel string, logger *logger.Logger) bool {
	images := ListImages(t, logger)

	imageTags := make([]string, 0, len(images))
	for i := range images {
		imageTags = append(imageTags, images[i].String())
	}

	return slices.Contains(imageTags, imgLabel)
}
