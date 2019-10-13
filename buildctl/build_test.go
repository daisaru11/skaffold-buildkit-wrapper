package buildctl

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuild(t *testing.T) {
	options := BuildOptions{
		Addr:       "tcp://0.0.0.0:1234",
		Context:    "./foo",
		Dockerfile: "./bar",
		Output:     "type=image,name=foo/bar,push=true",
		Opt:        []string{"build-arg:arg1"},
		Secret:     []string{"id=secret1,src=.secret1", "id=secret2,src=.secret2"},
	}

	got, err := Build(&options)
	if !assert.NoError(t, err) {
		return
	}

	assert.Equal(t, []string{
		"buildctl",
		"--addr",
		"tcp://0.0.0.0:1234",
		"build",
		"--frontend",
		"dockerfile.v0",
		"--local",
		"context=./foo",
		"--local",
		"dockerfile=./foo",
		"--output",
		"type=image,name=foo/bar,push=true",
		"--opt",
		"build-arg:arg1",
		"--secret",
		"id=secret1,src=.secret1",
		"--secret",
		"id=secret2,src=.secret2",
	}, got.Args)
}
