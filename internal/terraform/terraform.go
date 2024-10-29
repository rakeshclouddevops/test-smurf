package terraform

import (
	"context"
	"fmt"
	"os"
	"os/exec"

	"github.com/hashicorp/terraform-exec/tfexec"
	log "github.com/sirupsen/logrus"
)

func getTerraform() (*tfexec.Terraform, error) {
	// Attempt to find the Terraform binary in the system's PATH
	terraformBinary, err := exec.LookPath("terraform")
	if err != nil {
		log.Errorf("Terraform binary not found in PATH. Please install Terraform.")
		return nil, err
	}

	// Create a new Terraform instance using the found binary
	tf, err := tfexec.NewTerraform(".", terraformBinary)
	if err != nil {
		log.Errorf("Error creating Terraform instance: %v", err)
		return nil, err
	}

	log.Infof("Using Terraform binary at: %s", terraformBinary)
	return tf, nil
}

// terraform init cmd
func Init() error {

	tf, err := getTerraform()
	if err != nil {
		return err
	}

	log.Info("Initializing Terraform...")
	err = tf.Init(context.Background(), tfexec.Upgrade(true))
	if err != nil {
		log.Errorf("Terraform init failed: %v", err)
		return err
	}

	log.Info("Terraform initialized successfully.")
	return nil
}

// terraform plan cmd
func Plan() error {

	tf, err := getTerraform()
	if err != nil {
		return err
	}

	// Set the stdout and stderr to os.Stdout and os.Stderr to capture Terraform's output
	tf.SetStdout(os.Stdout)
	tf.SetStderr(os.Stderr)

	log.Info("Running Terraform plan...")

	// Run the plan and output to console
	_, err = tf.Plan(context.Background())
	if err != nil {
		log.Errorf("Terraform plan failed: %v", err)
		return err
	}

	return nil
}

// terraform apply cmd
func Apply() error {

	tf, err := getTerraform()
	if err != nil {
		return err
	}

	log.Info("Applying Terraform changes...")
	err = tf.Apply(context.Background())
	if err != nil {
		log.Errorf("Terraform apply failed: %v", err)
		return err
	}

	log.Info("Terraform applied successfully.")
	return nil
}

// terraform destroy cmd
func Destroy() error {

	tf, err := getTerraform()
	if err != nil {
		return err
	}

	log.Info("Destroying Terraform resources...")
	err = tf.Destroy(context.Background())
	if err != nil {
		log.Errorf("Terraform destroy failed: %v", err)
		return err
	}

	log.Info("Terraform resources destroyed successfully.")
	return nil
}

// terraform drift cmd
func DetectDrift() error {

	tf, err := getTerraform()
	if err != nil {
		return err
	}

	planFile := "drift.plan"
	log.Info("Checking for drift...")
	_, err = tf.Plan(context.Background(), tfexec.Out(planFile), tfexec.Refresh(true))

	if err != nil {
		log.Errorf("Terraform plan for drift detection failed: %v", err)
		return err
	}

	plan, err := tf.ShowPlanFile(context.Background(), planFile)
	if err != nil {
		log.Errorf("Error showing plan file: %v", err)
		return err
	}

	if len(plan.ResourceChanges) > 0 {
		log.Info("Drift detected:")
		for _, change := range plan.ResourceChanges {
			fmt.Printf("- %s: %s\n", change.Address, change.Change.Actions)
		}
	} else {
		log.Info("No drift detected.")
	}

	return nil
}

// terraform output cmd
func Output() error {
	tf, err := getTerraform()
	if err != nil {
		return err
	}

	// Set stdout and stderr to capture any messages
	tf.SetStdout(os.Stdout)
	tf.SetStderr(os.Stderr)

	// Refresh the state to ensure outputs are up to date
	err = tf.Refresh(context.Background())
	if err != nil {
		log.Errorf("Error refreshing Terraform state: %v", err)
		return err
	}

	outputs, err := tf.Output(context.Background())
	if err != nil {
		log.Errorf("Error getting Terraform outputs: %v", err)
		return err
	}

	if len(outputs) == 0 {
		log.Info("No outputs found.")
		return nil
	}

	log.Info("Terraform outputs:")
	for key, value := range outputs {
		// Check if the output is sensitive
		if value.Sensitive {
			fmt.Printf("%s: [sensitive value hidden]\n", key)
		} else {
			fmt.Printf("%s: %v\n", key, value.Value)
		}
	}

	return nil
}
