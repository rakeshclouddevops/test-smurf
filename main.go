/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"github.com/clouddrove/smurf/cmd"
	_ "github.com/clouddrove/smurf/cmd/apply"
	_ "github.com/clouddrove/smurf/cmd/destroy"
	_ "github.com/clouddrove/smurf/cmd/drift"
	_ "github.com/clouddrove/smurf/cmd/init"
	_ "github.com/clouddrove/smurf/cmd/output"
	_ "github.com/clouddrove/smurf/cmd/plan"
	_ "github.com/clouddrove/smurf/cmd/provision"
)

func main() {
	cmd.Execute()
}
