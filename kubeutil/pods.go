package kubeutil

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"sort"
	"strings"
)

func GetPodsBySelector(selector string) ([]string, error) {
	buf := new(bytes.Buffer)
	cmd := getPodsBySelectorCommand(selector, buf, os.Stderr)
	err := cmd.Run()
	if err != nil {
		return nil, err
	}

	b, err := ioutil.ReadAll(buf)
	if err != nil {
		return nil, err
	}

	var names []string
	for _, n := range strings.Split(string(b), "\n") {
		n = strings.TrimSpace(n)
		if n != "" {
			names = append(names, n)
		}
	}
	sort.Strings(names)

	return names, nil
}

func getPodsBySelectorCommand(selector string, stdout, stderr io.Writer) *exec.Cmd {
	cmd := exec.Command("kubectl", "get", "pods")
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	cmd.Env = os.Environ()

	cmd.Args = append(cmd.Args, "--selector", selector)
	cmd.Args = append(cmd.Args, "--no-headers", "-o", "custom-columns=:metadata.name")

	return cmd
}
