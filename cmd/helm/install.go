package helm

import (
    "github.com/clouddrove/smurf/internal/helm"
    "github.com/spf13/cobra"
)

var n string // This will hold the namespace value from the flag

var installCmd = &cobra.Command{
    Use:   "install [RELEASE] [CHART]",
    Short: "Install a Helm chart into a Kubernetes cluster.",
    Args:  cobra.ExactArgs(2),
    RunE: func(cmd *cobra.Command, args []string) error {
        releaseName := args[0]
        chartPath := args[1]
        if n == "" { // If no namespace is provided, use default
            n = "default"
        }
        return helm.HelmInstall(releaseName, chartPath, n)
    },
}

func init() {
    // Add the namespace flag
    installCmd.Flags().StringVarP(&n, "namespace", "n", "", "Specify the namespace to install the Helm chart into")
    selmCmd.AddCommand(installCmd)
}
