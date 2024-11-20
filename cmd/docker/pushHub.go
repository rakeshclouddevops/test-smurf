package docker

import (
	"github.com/clouddrove/smurf/internal/docker"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var (
	hubImageName       string
	hubImageTag        string
	hubDeleteAfterPush bool
)

var pushHubCmd = &cobra.Command{
	Use:   "hub",
	Short: "push Docker images to Docker Hub",
	Long: `
	Push Docker images to Docker Hub
	export DOCKER_USERNAME=<username>
	export DOCKER_PASSWORD=<password>
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		opts := docker.PushOptions{
			ImageName: hubImageName,
		}
		if err := docker.PushImage(opts); err != nil {
			return err
		}
		if hubDeleteAfterPush {
			if err := docker.RemoveImage(hubImageName); err != nil {
				return err
			}
			pterm.Success.Println("Successfully deleted local image:", hubImageName)
		}
		return nil
	},
	Example: `
	smurf sdkr push hub --image <image-name> --tag <image-tag>
	smurf sdkr push hub --image <image-name> --tag <image-tag> --delete
	`,
}

func init() {
	pushHubCmd.Flags().StringVarP(&hubImageName, "image", "i", "", "Image name (e.g., myapp)")
	pushHubCmd.Flags().StringVarP(&hubImageTag, "tag", "t", "latest", "Image tag (default: latest)")
	pushHubCmd.Flags().BoolVarP(&hubDeleteAfterPush, "delete", "d", false, "Delete the local image after pushing")

	pushHubCmd.MarkFlagRequired("image")

	pushCmd.AddCommand(pushHubCmd)
}
