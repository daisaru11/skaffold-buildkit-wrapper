package consistenthash

import (
	"fmt"

	"github.com/serialx/hashring"
)

func Choice(nodes []string, key string) (string, error) {
	ring := hashring.New(nodes)
	x, ok := ring.GetNode(key)
	if !ok {
		return "", fmt.Errorf("no node found for key %q", key)
	}
	return x, nil
}
