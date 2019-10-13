package docker

import (
	"os"
	"os/exec"
)

type BuildOptions struct {
	Context    string
	Dockerfile string
	Tag        string
	BuildArg   []string
	Secret     []string
}

func Build(options *BuildOptions) (*exec.Cmd, error) {

	cmd := exec.Command("docker", "build")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.Env = append(os.Environ(),
		"DOCKER_BUILDKIT=1",
	)

	if options.Tag != "" {
		cmd.Args = append(cmd.Args, "--tag", options.Tag)
	}

	if options.Dockerfile != "" {
		cmd.Args = append(cmd.Args, "--file", options.Dockerfile)
	}

	if options.BuildArg != nil {
		for _, a := range options.BuildArg {
			cmd.Args = append(cmd.Args, "--build-arg", a)
		}
	}

	if options.Secret != nil {
		for _, o := range options.Secret {
			cmd.Args = append(cmd.Args, "--secret", o)
		}
	}

	if options.Context != "" {
		cmd.Args = append(cmd.Args, options.Context)
	}

	return cmd, nil
}

func Push(name string) (*exec.Cmd, error) {
	cmd := exec.Command("docker", "push", name)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd, nil
}
