package helm

import (
	"fmt"

	"github.com/clouddrove/smurf/cmd"
	"github.com/spf13/cobra"
)

// selmCmd represents the 'selm' command
var selmCmd = &cobra.Command{
	Use:   "selm",
	Short: "Subcommand for Helm-related actions",
	Long:  `selm is a subcommand that groups various Helm-related actions under a single command.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Use 'smurf selm [command]' to run Helm-related actions")
	},
}

func init() {
	// Add the 'selm' command as a subcommand of RootCmd
	cmd.RootCmd.AddCommand(selmCmd)
}
