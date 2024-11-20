package docker

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/clouddrove/smurf/internal/docker"
	"github.com/spf13/cobra"
)

var (
	dockerfilePath string
	noCache        bool
	buildArgs      []string
	target         string
	platform       string
	contextDir     string
)

var buildCmd = &cobra.Command{
	Use:   "build [IMAGE_NAME] [TAG]",
	Short: "Build a Docker image with the given name and tag.",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		buildArgsMap := make(map[string]string)
		for _, arg := range buildArgs {
			parts := strings.SplitN(arg, "=", 2)
			if len(parts) == 2 {
				buildArgsMap[parts[0]] = parts[1]
			}
		}

		if contextDir == "" {
			wd, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("failed to get current working directory: %w", err)
			}
			contextDir = wd
		}

		if dockerfilePath == "" {
			dockerfilePath = filepath.Join(contextDir, "Dockerfile")
		} else {
			dockerfilePath = filepath.Join(contextDir, dockerfilePath)
		}

		if _, err := os.Stat(dockerfilePath); os.IsNotExist(err) {
			return fmt.Errorf("dockerfile not found at %s", dockerfilePath)
		}

		opts := docker.BuildOptions{
			ContextDir:     contextDir,
			DockerfilePath: dockerfilePath,
			NoCache:        noCache,
			BuildArgs:      buildArgsMap,
			Target:         target,
			Platform:       platform,
		}

		return docker.Build(args[0], args[1], opts)
	},
	Example: `
	smurf sdkr build my-image my-tag
	smurf sdkr build my-image my-tag --file Dockerfile --context ./build-context --no-cache --build-arg key1=value1 --build-arg key2=value2 --target my-target --platform linux/amd64`,
}

func init() {
	buildCmd.Flags().StringVarP(&dockerfilePath, "file", "f", "", "Path to Dockerfile relative to context directory")
	buildCmd.Flags().StringVar(&contextDir, "context", "", "Build context directory (default: current directory)")
	buildCmd.Flags().BoolVar(&noCache, "no-cache", false, "Do not use cache when building the image")
	buildCmd.Flags().StringArrayVar(&buildArgs, "build-arg", []string{}, "Set build-time variables")
	buildCmd.Flags().StringVar(&target, "target", "", "Set the target build stage to build")
	buildCmd.Flags().StringVar(&platform, "platform", "", "Set the platform for the build (e.g., linux/amd64, linux/arm64)")

	sdkrCmd.AddCommand(buildCmd)
}