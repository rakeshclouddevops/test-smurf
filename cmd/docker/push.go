package docker

import (
    "github.com/clouddrove/smurf/internal/docker"
    "github.com/pterm/pterm"
    "github.com/spf13/cobra"
)

var imageName string
var deleteAfterPush bool

var pushCmd = &cobra.Command{
    Use:   "push",
    Short: "Push a Docker image to a Docker registry",
    Long: `Push a Docker image to a Docker registry using the provided image name and tag.
Ensure that DOCKER_USERNAME and DOCKER_PASSWORD environment variables are set.`,
    RunE: func(cmd *cobra.Command, args []string) error {
        opts := docker.PushOptions{
            ImageName: imageName,
        }
        if err := docker.PushImage(opts); err != nil {
            return err
        }
        if deleteAfterPush {
            if err := docker.RemoveImage(imageName); err != nil {
                return err
            }
            pterm.Success.Println("Successfully deleted local image:", imageName)
        }
        return nil
    },
}

func init() {
    pushCmd.Flags().StringVarP(&imageName, "image", "i", "", "Full image name including the tag (e.g., username/repository:tag)")
    pushCmd.Flags().BoolVarP(&deleteAfterPush, "delete", "d", false, "Delete the local image after pushing")
    pushCmd.MarkFlagRequired("image")

    sdkrCmd.AddCommand(pushCmd)
}
