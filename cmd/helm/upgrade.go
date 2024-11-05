package helm

import (
	"github.com/clouddrove/smurf/internal/helm"
	"github.com/spf13/cobra"
)

var upgradeCmd = &cobra.Command{
	Use:   "upgrade [NAME] [CHART]",
	Short: "Upgrade a deployed Helm chart.",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		releaseName := args[0]
		chartPath := args[1]
		return helm.HelmUpgrade(releaseName, chartPath, "default") // Assuming 'default' namespace or make it a flag
	},
}

func init() {
	selmCmd.AddCommand(upgradeCmd)
}
