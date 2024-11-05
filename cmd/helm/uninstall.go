package helm

import (
	"github.com/clouddrove/smurf/internal/helm"
	"github.com/spf13/cobra"
)

var uninstallCmd = &cobra.Command{
	Use:   "uninstall [NAME]",
	Short: "Uninstall a Helm release.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		releaseName := args[0]
		return helm.HelmUninstall(releaseName, "default") // Assuming 'default' namespace or make it a flag
	},
}

func init() {
	selmCmd.AddCommand(uninstallCmd)
}
