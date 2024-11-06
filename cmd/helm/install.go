package helm

import (
	"github.com/clouddrove/smurf/internal/helm"
	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use:   "install [RELEASE] [CHART]",
	Short: "Install a Helm chart into a Kubernetes cluster.",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		releaseName := args[0]
		chartPath := args[1]
		return helm.HelmInstall(releaseName, chartPath, "default") // Assuming 'default' namespace or make it a flag
	},
}

func init() {
	selmCmd.AddCommand(installCmd)
}
