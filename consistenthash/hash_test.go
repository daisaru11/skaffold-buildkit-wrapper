package consistenthash

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChoice(t *testing.T) {
	nodes := []string{
		"buildkitd-0",
		"buildkitd-1",
		"buildkitd-2",
		"buildkitd-3",
	}

	keys := []string{
		"github.com/daisaru11/repo_1",
		"github.com/daisaru11/repo_2",
		"github.com/daisaru11/repo_3",
		"github.com/daisaru11/repo_4",
	}

	got := []string{}
	for _, k := range keys {
		c, err := Choice(nodes, k)
		if !assert.NoError(t, err) {
			return
		}

		got = append(got, c)
	}

	want := []string{
		"buildkitd-1",
		"buildkitd-1",
		"buildkitd-0",
		"buildkitd-3",
	}

	assert.Equal(t, want, got)

	// Test Re-balancing
	nodes = append(nodes, "buildkitd-4")

	got = []string{}
	for _, k := range keys {
		c, err := Choice(nodes, k)
		if !assert.NoError(t, err) {
			return
		}

		got = append(got, c)
	}

	want = []string{
		"buildkitd-1",
		"buildkitd-1",
		"buildkitd-4",
		"buildkitd-3",
	}

	assert.Equal(t, want, got)
}
