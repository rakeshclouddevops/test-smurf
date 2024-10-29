/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)



// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "smurf",
	Short: "Smurf is a tool for automating common commands across Terraform, Docker, and more, streamlining DevOps workflows.",
	Long: `Smurf is a command-line interface built with Cobra, designed to simplify and automate commands for essential tools like Terraform and Docker. It provides intuitive, unified commands to execute Terraform plans, Docker container management, and other DevOps tasks seamlessly from one interface. Whether you need to spin up environments, manage containers, or apply infrastructure as code, this CLI streamlines multi-tool operations, boosting productivity and reducing context-switching.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.smurf.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}


