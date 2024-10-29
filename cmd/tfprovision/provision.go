package tfplan

import (
	"github.com/clouddrove/smurf/cmd"
	"github.com/clouddrove/smurf/internal/terraform"
	"github.com/spf13/cobra"
)

var provisionCmd = &cobra.Command{
	Use:   "provision",
	Short: "Its the combination of init, drift, plan, apply, output for Terraform",
	RunE: func(cmd *cobra.Command, args []string) error {

		err := terraform.Init()

		if err != nil {
			return err
		}

		err = terraform.DetectDrift()

		if err != nil {
			return err
		}

		err = terraform.Plan()

		if err != nil {
			return err
		}

		err = terraform.Apply()

		if err != nil {
			return err
		}

		return terraform.Output()
	},
}

func init() {
	cmd.RootCmd.AddCommand(provisionCmd)
}
