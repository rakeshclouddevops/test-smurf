/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"github.com/clouddrove/smurf/cmd"
	_ "github.com/clouddrove/smurf/cmd/tfapply"
	_ "github.com/clouddrove/smurf/cmd/tfdestroy"
	_ "github.com/clouddrove/smurf/cmd/tfdrift"
	_ "github.com/clouddrove/smurf/cmd/tfinit"
	_ "github.com/clouddrove/smurf/cmd/tfoutput"
	_ "github.com/clouddrove/smurf/cmd/tfplan"
	_ "github.com/clouddrove/smurf/cmd/tfprovision"
)

func main() {
	cmd.Execute()
}
