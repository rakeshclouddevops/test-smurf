// File: cmd/push.go

package docker

import (
	"fmt"

	"github.com/clouddrove/smurf/internal/docker"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var (
	imageNamePush         string
	imageTagPush          string
	deleteAfterPushPush   bool
	useAWSPush            bool
	useAzurePush          bool
	useGCPPush            bool
	regionPush            string
	repositoryNamePush    string
	projectIDPush         string
	subscriptionIDPush    string
	resourceGroupNamePush string
	registryNamePush      string
)

var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "Push a Docker image to a Docker registry",
	Long: `Push a Docker image to a Docker registry using the provided image name and tag.
Ensure that you have authenticated with the appropriate cloud provider.
Set the GOOGLE_APPLICATION_CREDENTIALS environment variable to the path of your service account JSON key file.
export GOOGLE_APPLICATION_CREDENTIALS="/path/to/your/service-account-key.json"`,
	RunE: func(cmd *cobra.Command, args []string) error {
		providersUsed := 0
		if useAWSPush {
			providersUsed++
		}
		if useAzurePush {
			providersUsed++
		}
		if useGCPPush {
			providersUsed++
		}
		if providersUsed != 1 {
			return fmt.Errorf("you must specify exactly one of --aws, --azure, or --gcp flags")
		}

		var fullImageNamePush string
		if useAWSPush {
			if regionPush == "" || repositoryNamePush == "" {
				return fmt.Errorf("--aws requires both --region and --repository flags")
			}

			fullImageNamePush = fmt.Sprintf("%s.dkr.ecr.%s.amazonaws.com/%s:%s", imageNamePush, regionPush, repositoryNamePush, imageTagPush)
			pterm.Info.Println("Pushing image to AWS ECR...")
			if err := docker.PushImageToECR(imageNamePush, regionPush, repositoryNamePush); err != nil {
				return err
			}
			pterm.Success.Println("Successfully pushed image to ECR:", fullImageNamePush)
		} else if useAzurePush {
			if subscriptionIDPush == "" || resourceGroupNamePush == "" || registryNamePush == "" {
				return fmt.Errorf("--azure requires --subscription-id, --resource-group, and --registry-name flags")
			}

			fullImageNamePush = fmt.Sprintf("%s.azurecr.io/%s:%s", registryNamePush, imageNamePush, imageTagPush)

			pterm.Info.Println("Pushing image to Azure Container Registry...")
			if err := docker.PushImageToACR(subscriptionIDPush, resourceGroupNamePush, registryNamePush, imageNamePush); err != nil {
				return err
			}
			pterm.Success.Println("Successfully pushed image to ACR:", fullImageNamePush)
		} else if useGCPPush {
			if projectIDPush == "" {
				return fmt.Errorf("--gcp requires --project-id flag")
			}

			fullImageNamePush = fmt.Sprintf("gcr.io/%s/%s:%s", projectIDPush, imageNamePush, imageTagPush)

			pterm.Info.Println("Pushing image to Google Container Registry...")
			if err := docker.PushImageToGCR(projectIDPush, imageNamePush); err != nil {
				return err
			}
			pterm.Success.Println("Successfully pushed image to GCR:", fullImageNamePush)
		}

		if deleteAfterPushPush {
			if err := docker.RemoveImage(fullImageNamePush); err != nil {
				return err
			}
			pterm.Success.Println("Successfully deleted local image:", fullImageNamePush)
		}
		return nil
	},
}

func init() {
	pushCmd.Flags().StringVarP(&imageNamePush, "image", "i", "", "Image name (e.g., myapp)")
	pushCmd.Flags().StringVarP(&imageTagPush, "tag", "t", "latest", "Image tag (default: latest)")
	pushCmd.Flags().BoolVarP(&deleteAfterPushPush, "delete", "d", false, "Delete the local image after pushing")

	pushCmd.Flags().BoolVar(&useAWSPush, "aws", false, "Push the image to AWS ECR")
	pushCmd.Flags().StringVarP(&regionPush, "region", "r", "", "AWS region (required with --aws)")
	pushCmd.Flags().StringVarP(&repositoryNamePush, "repository", "R", "", "AWS ECR repository name (required with --aws)")

	pushCmd.Flags().BoolVar(&useAzurePush, "azure", false, "Push the image to Azure Container Registry")
	pushCmd.Flags().StringVar(&subscriptionIDPush, "subscription-id", "", "Azure subscription ID (required with --azure)")
	pushCmd.Flags().StringVar(&resourceGroupNamePush, "resource-group", "", "Azure resource group name (required with --azure)")
	pushCmd.Flags().StringVar(&registryNamePush, "registry-name", "", "Azure Container Registry name (required with --azure)")

	pushCmd.Flags().BoolVar(&useGCPPush, "gcp", false, "Push the image to Google Container Registry")
	pushCmd.Flags().StringVar(&projectIDPush, "project-id", "", "GCP project ID (required with --gcp)")

	pushCmd.MarkFlagRequired("image")
	sdkrCmd.AddCommand(pushCmd)
}
