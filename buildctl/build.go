package buildctl

import (
	"fmt"
	"os"
	"os/exec"
)

type BuildOptions struct {
	Addr       string
	Context    string
	Dockerfile string
	Output     string
	Opt        []string
	Secret     []string
}

func Build(options *BuildOptions) (*exec.Cmd, error) {

	cmd := exec.Command("buildctl")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()

	if options.Addr != "" {
		cmd.Args = append(cmd.Args, "--addr", options.Addr)
	}

	cmd.Args = append(cmd.Args, "build")
	cmd.Args = append(cmd.Args, "--frontend", "dockerfile.v0")

	if options.Context != "" {
		cmd.Args = append(cmd.Args, "--local", fmt.Sprintf("context=%s", options.Context))
	}
	if options.Dockerfile != "" {
		cmd.Args = append(cmd.Args, "--local", fmt.Sprintf("dockerfile=%s", options.Context))
	}

	if options.Output != "" {
		cmd.Args = append(cmd.Args, "--output", options.Output)
	}

	if options.Opt != nil {
		for _, o := range options.Opt {
			cmd.Args = append(cmd.Args, "--opt", o)
		}
	}

	if options.Secret != nil {
		for _, o := range options.Secret {
			cmd.Args = append(cmd.Args, "--secret", o)
		}
	}

	return cmd, nil
}
