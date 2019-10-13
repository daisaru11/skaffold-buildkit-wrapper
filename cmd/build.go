package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/daisaru11/skaffold-buildkit-wrapper/buildctl"
	"github.com/daisaru11/skaffold-buildkit-wrapper/consistenthash"
	"github.com/daisaru11/skaffold-buildkit-wrapper/docker"
	"github.com/daisaru11/skaffold-buildkit-wrapper/kubeutil"
	"github.com/spf13/cobra"
)

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "build images with buildctl",
	Long:  `build images with buildctl`,
	RunE:  runBuild,
}

func init() {
	rootCmd.AddCommand(buildCmd)

	buildCmd.Flags().String("cli", "docker", "CLI used to build image")
	buildCmd.Flags().String("file", "Dockerfile", "Name of Dockerfile")
	buildCmd.Flags().StringArray("build-arg", []string{}, "Build argument")
	buildCmd.Flags().StringArray("secret", []string{}, "Secret")

	buildCmd.Flags().String("addr", "", "Address of buildkitd")
	buildCmd.Flags().String("kube-pod-selector", "", "Selector of kubernetes pod buildkitd running")
	buildCmd.Flags().String("kube-pod-balancing-hash-key", "", "Key of consistent hashing to choice pods")
}

func runBuild(cmd *cobra.Command, args []string) error {
	var err error
	cli, err := cmd.Flags().GetString("cli")
	if err != nil {
		return err
	}

	switch cli {
	case "docker":
		return runBuildWithDocker(cmd, args)
	case "buildctl":
		return runBuildWithBuildctl(cmd, args)
	default:
		return fmt.Errorf("unknown cli")
	}
}

func runBuildWithDocker(cmd *cobra.Command, args []string) error {
	var err error

	options := docker.BuildOptions{}
	err = setDockerBuildOptionsFromCustomBuildEnv(&options)
	if err != nil {
		return fmt.Errorf("failed to set build options from environment variables: %w", err)
	}

	err = setDockerBuildOptionsFromFlags(&options, cmd)
	if err != nil {
		return fmt.Errorf("failed to set build options from flags: %w", err)
	}

	c, err := docker.Build(&options)
	if err != nil {
		return fmt.Errorf("failed to build command: %w", err)
	}
	log.Printf("run: %s", c.String())
	err = c.Run()
	if err != nil {
		return fmt.Errorf("failed to run docker cli: %w", err)
	}
	log.Printf("build finished successfully with buildctl")

	if os.Getenv("PUSH_IMAGE") == "true" {
		c, err := docker.Push(options.Tag)
		if err != nil {
			return fmt.Errorf("failed to build command: %w", err)
		}
		err = c.Run()
		if err != nil {
			return fmt.Errorf("failed to run docker cli: %w", err)
		}
	}

	return nil
}

func setDockerBuildOptionsFromCustomBuildEnv(options *docker.BuildOptions) error {
	context := os.Getenv("BUILD_CONTEXT")
	if context != "" {
		options.Context = context
	}

	images := os.Getenv("IMAGES")
	if images != "" {
		options.Tag = images
	}

	return nil
}

func setDockerBuildOptionsFromFlags(options *docker.BuildOptions, cmd *cobra.Command) error {
	file, err := cmd.Flags().GetString("file")
	if err != nil {
		return err
	}
	options.Dockerfile = file

	args, err := cmd.Flags().GetStringArray("build-arg")
	if err != nil {
		return err
	}
	options.BuildArg = args

	secret, err := cmd.Flags().GetStringArray("secret")
	if err != nil {
		return err
	}
	options.Secret = secret

	return nil
}

func runBuildWithBuildctl(cmd *cobra.Command, args []string) error {
	var err error

	options := buildctl.BuildOptions{}
	err = setBuildctlBuildOptionsFromCustomBuildEnv(&options)
	if err != nil {
		return fmt.Errorf("failed to set build options from environment variables: %w", err)
	}

	err = setBuildctlBuildOptionsFromFlags(&options, cmd)
	if err != nil {
		return fmt.Errorf("failed to set build options from flags: %w", err)
	}

	c, err := buildctl.Build(&options)
	if err != nil {
		return fmt.Errorf("failed to build command: %w", err)
	}
	log.Printf("run: %s", c.String())
	err = c.Run()
	if err != nil {
		return fmt.Errorf("failed to run buildctl: %w", err)
	}
	log.Printf("build finished successfully with buildctl")

	return nil
}

func setBuildctlBuildOptionsFromCustomBuildEnv(options *buildctl.BuildOptions) error {
	context := os.Getenv("BUILD_CONTEXT")
	if context != "" {
		options.Context = context
		options.Dockerfile = context
	}

	images := os.Getenv("IMAGES")
	push := os.Getenv("PUSH_IMAGE")
	if images != "" {
		options.Output = fmt.Sprintf("type=image,name=%s,push=%s", images, push)
	}

	return nil
}

func setBuildctlBuildOptionsFromFlags(options *buildctl.BuildOptions, cmd *cobra.Command) error {
	addr, err := cmd.Flags().GetString("addr")
	if err != nil {
		return err
	}
	options.Addr = addr

	kubePodSelector, err := cmd.Flags().GetString("kube-pod-selector")
	if err != nil {
		return err
	}
	kubePodBalancingHashKey, err := cmd.Flags().GetString("kube-pod-balancing-hash-key")
	if err != nil {
		return err
	}
	if kubePodSelector != "" && kubePodBalancingHashKey != "" {
		pods, err := kubeutil.GetPodsBySelector(kubePodSelector)
		if err != nil {
			return err
		}
		pod, err := consistenthash.Choice(pods, kubePodBalancingHashKey)
		if err != nil {
			return err
		}
		options.Addr = fmt.Sprintf("kube-pod://%s", pod)
	}

	options.Opt = []string{}

	file, err := cmd.Flags().GetString("file")
	if err != nil {
		return err
	}
	options.Opt = append(options.Opt, fmt.Sprintf("filename=%s", file))

	args, err := cmd.Flags().GetStringArray("build-arg")
	if err != nil {
		return err
	}
	for _, arg := range args {
		options.Opt = append(options.Opt, fmt.Sprintf("build-arg:%s", arg))
	}

	secret, err := cmd.Flags().GetStringArray("secret")
	if err != nil {
		return err
	}
	options.Secret = secret

	return nil
}
