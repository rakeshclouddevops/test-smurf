package helm

import (
    "github.com/clouddrove/smurf/internal/helm"
    "github.com/spf13/cobra"
)

var installNamespace string // This will hold the namespace value from the flag

var installCmd = &cobra.Command{
    Use:   "install [RELEASE] [CHART]",
    Short: "Install a Helm chart into a Kubernetes cluster.",
    Args:  cobra.ExactArgs(2),
    RunE: func(cmd *cobra.Command, args []string) error {
        releaseName := args[0]
        chartPath := args[1]
        if installNamespace == "" { // If no namespace is provided, use default
            installNamespace = "default"
        }
        return helm.HelmInstall(releaseName, chartPath, installNamespace)
    },
}

func init() {
    // Add the namespace flag
    installCmd.Flags().StringVarP(&installNamespace, "namespace", "n", "", "Specify the namespace to install the Helm chart")
    selmCmd.AddCommand(installCmd)
}
