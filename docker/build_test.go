package docker

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuild(t *testing.T) {
	options := BuildOptions{
		Context:    ".",
		Tag:        "foo/bar",
		Dockerfile: "foo.Dockerfile",
		BuildArg:   []string{"arg1", "arg2"},
		Secret:     []string{"secret1", "secret2"},
	}

	got, err := Build(&options)
	if !assert.NoError(t, err) {
		return
	}

	assert.Equal(t, []string{
		"docker",
		"build",
		"--tag",
		"foo/bar",
		"--file",
		"foo.Dockerfile",
		"--build-arg",
		"arg1",
		"--build-arg",
		"arg2",
		"--secret",
		"secret1",
		"--secret",
		"secret2",
		".",
	}, got.Args)
	assert.Contains(t, got.Env, "DOCKER_BUILDKIT=1")
}
