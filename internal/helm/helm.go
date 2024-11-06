package helm

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/fatih/color"
	"github.com/pterm/pterm"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/chartutil"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/release"
	"helm.sh/helm/v3/pkg/strvals"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var settings = cli.New()

func init() {
	if os.Getenv("KUBECONFIG") != "" {
		settings.KubeConfig = os.Getenv("KUBECONFIG")
	} else {
		home := homedir.HomeDir()
		settings.KubeConfig = filepath.Join(home, ".kube", "config")
	}
}

func getKubeClient() (*kubernetes.Clientset, error) {
	var kubeconfig string
	if os.Getenv("KUBECONFIG") != "" {
		kubeconfig = os.Getenv("KUBECONFIG")
	} else {
		home := homedir.HomeDir()
		kubeconfig = filepath.Join(home, ".kube", "config")
	}

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return clientset, nil
}

func CreateChart(chartName, saveDir string) error {
	spinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Installing release '%s'", chartName))
	// Ensure the directory exists
	if _, err := os.Stat(saveDir); os.IsNotExist(err) {
		os.MkdirAll(saveDir, os.ModePerm)
	}

	// Create the chart in the specified directory
	_, err := chartutil.Create(chartName, saveDir)
	if err != nil {
		spinner.Fail()
		return err
	}

	spinner.Success(fmt.Sprintf("Chart created successfully: %s", chartName))
	return nil
}

func HelmInstall(releaseName, chartPath, namespace string) error {
	spinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Installing release '%s'", releaseName))

	// Initialize action configuration
	actionConfig := new(action.Configuration)
	debugLog := func(format string, v ...interface{}) {
		fmt.Printf(format, v...)
	}
	if err := actionConfig.Init(settings.RESTClientGetter(), namespace, os.Getenv("HELM_DRIVER"), debugLog); err != nil {
		spinner.Fail("Failed to initialize Helm action configuration: " + err.Error())
		return err
	}

	kubeClient, err := getKubeClient()

	_ = kubeClient
	if err != nil {
		spinner.Fail("Failed to get Kubernetes client: " + err.Error())
		return err
	}

	client := action.NewInstall(actionConfig)
	client.ReleaseName = releaseName
	client.Namespace = namespace

	// Load the chart
	chart, err := loader.Load(chartPath)
	if err != nil {
		spinner.Fail("Failed to load chart: " + err.Error())
		return err
	}

	_, err = client.Run(chart, map[string]interface{}{})
	if err != nil {
		spinner.Fail("Installation failed: " + err.Error())
		return err
	}

	spinner.Success(fmt.Sprintf("Release '%s' installed successfully", releaseName))
	return nil
}

func EnsureNamespace(namespace string, create bool) error {
	kubeconfig := filepath.Join(homedir.HomeDir(), ".kube", "config")
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}

	_, err = clientset.CoreV1().Namespaces().Get(context.TODO(), namespace, metav1.GetOptions{})
	if err == nil {
		return nil // Namespace exists
	}

	if create {
		nsSpec := &v1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: namespace}}
		_, err = clientset.CoreV1().Namespaces().Create(context.TODO(), nsSpec, metav1.CreateOptions{})
		if err != nil {
			return err
		}
	}
	return nil
}

func HelmUpgrade(releaseName, chartPath, namespace string, setValues []string, valuesFiles []string, createNamespace, atomic bool, timeout time.Duration, debug bool) error {
	spinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Upgrading release '%s'", releaseName))
	settings := cli.New()
	if debug {
		settings.Debug = true
	}

	if createNamespace {
		if err := EnsureNamespace(namespace, true); err != nil {
			spinner.Fail("Failed to ensure namespace: " + err.Error())
			return err
		}
	}

	actionConfig := new(action.Configuration)
	debugLog := func(format string, v ...interface{}) {
		if settings.Debug {
			fmt.Printf(format, v...)
		}
	}
	if err := actionConfig.Init(settings.RESTClientGetter(), namespace, os.Getenv("HELM_DRIVER"), debugLog); err != nil {
		spinner.Fail("Failed to initialize Helm action configuration: " + err.Error())
		return err
	}

	client := action.NewUpgrade(actionConfig)
	client.Namespace = namespace
	client.Atomic = atomic
	client.Timeout = timeout

	chart, err := loader.Load(chartPath)
	if err != nil {
		spinner.Fail("Failed to load chart: " + err.Error())
		return err
	}

	vals := make(map[string]interface{})
	for _, f := range valuesFiles {
		currentVals, err := chartutil.ReadValuesFile(f)
		if err != nil {
			spinner.Fail("Failed to read values file: " + err.Error())
			return err
		}
		chartutil.CoalesceTables(vals, currentVals)
	}

	for _, set := range setValues {
		if err := strvals.ParseInto(set, vals); err != nil {
			spinner.Fail("Failed to parse set values: " + err.Error())
			return err
		}
	}

	_, err = client.Run(releaseName, chart, vals)
	if err != nil {
		spinner.Fail("Upgrade failed: " + err.Error())
	} else {
		spinner.Success(fmt.Sprintf("Release '%s' upgraded successfully", releaseName))
	}
	return err
}

func HelmList(namespace string) ([]*release.Release, error) {
	spinner, _ := pterm.DefaultSpinner.Start("Listing releases")
	actionConfig := new(action.Configuration)
	if err := actionConfig.Init(settings.RESTClientGetter(), namespace, "secrets", nil); err != nil {
		spinner.Fail(err.Error())
		return nil, err
	}
	client := action.NewList(actionConfig)
	// client.Namespace = namespace

	releases, err := client.Run()
	if err != nil {
		spinner.Fail(err.Error())
	} else {
		spinner.Success("Releases listed successfully")
	}
	return releases, err
}

func HelmUninstall(releaseName, namespace string) error {
	spinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Uninstalling release '%s'", releaseName))

	actionConfig := new(action.Configuration)
	// Initialize the configuration with a debug function for better tracing
	debugLog := func(format string, v ...interface{}) {
		fmt.Printf(format, v...)
	}
	if err := actionConfig.Init(settings.RESTClientGetter(), namespace, os.Getenv("HELM_DRIVER"), debugLog); err != nil {
		spinner.Fail("Failed to initialize Helm action configuration: " + err.Error())
		return err
	}

	// Create a new uninstall action
	client := action.NewUninstall(actionConfig)
	if client == nil {
		spinner.Fail("Failed to create Helm uninstall client")
		return fmt.Errorf("failed to create Helm uninstall client")
	}

	// Perform the uninstall
	_, err := client.Run(releaseName)
	if err != nil {
		spinner.Fail("Uninstall failed: " + err.Error())
		return err
	}

	spinner.Success(fmt.Sprintf("Release '%s' uninstalled successfully", releaseName))
	return nil
}

func HelmLint(chartPath string) error {
	spinner, _ := pterm.DefaultSpinner.Start("Linting chart")
	actionConfig := new(action.Configuration)
	_ = actionConfig
	client := action.NewLint()

	result := client.Run([]string{chartPath}, nil)
	if len(result.Messages) > 0 {
		for _, msg := range result.Messages {
			fmt.Println(msg)
		}
		spinner.Fail("Linting issues found")
	} else {
		spinner.Success("No linting issues")
	}
	return nil
}

// HelmTemplate renders the Helm templates for a given chart
func HelmTemplate(releaseName, chartPath, namespace string) error {
	settings := cli.New() // Ensures that we have a new Helm environment setup
	actionConfig := new(action.Configuration)

	if err := actionConfig.Init(settings.RESTClientGetter(), namespace, os.Getenv("HELM_DRIVER"), nil); err != nil {
		pterm.DefaultBasicText.WithStyle(pterm.NewStyle(pterm.FgRed)).Println(err.Error())
		return err
	}

	client := action.NewInstall(actionConfig)
	client.DryRun = true // Important: Dry run must be true for template rendering
	client.ReleaseName = releaseName
	client.Namespace = namespace
	client.Replace = true    // Re-use the release name without an error for repeated calls
	client.ClientOnly = true // Do not validate against the Kubernetes cluster

	// Load the chart
	chart, err := loader.Load(chartPath)
	if err != nil {
		pterm.DefaultBasicText.WithStyle(pterm.NewStyle(pterm.FgRed)).Println(err.Error())
		return err
	}

	spinner, _ := pterm.DefaultSpinner.Start("Rendering Helm templates...")
	// Run the installation simulation which renders the template
	rel, err := client.Run(chart, nil) // Pass nil because we do not have any overridden values
	if err != nil {
		spinner.Fail(err.Error())
		return err
	}
	spinner.Success("Templates rendered successfully")

	// Output the rendered templates
	green := color.New(color.FgGreen).SprintFunc()
	fmt.Println(green(rel.Manifest))

	return nil
}

// Provision provisions a Helm chart by installing or upgrading it, linting it, and rendering its templates
func Provision(releaseName, chartPath, namespace string) error {
	settings := cli.New()
	actionConfig := new(action.Configuration)
	if err := actionConfig.Init(settings.RESTClientGetter(), namespace, os.Getenv("HELM_DRIVER"), nil); err != nil {
		return err
	}

	client := action.NewList(actionConfig)
	results, err := client.Run()
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	var installErr, upgradeErr, lintErr, templateErr error

	exists := false
	for _, result := range results {
		if result.Name == releaseName {
			exists = true
			break
		}
	}

	wg.Add(1)
	if exists {
		go func() {
			defer wg.Done()
			upgradeErr = HelmUpgrade(releaseName, chartPath, namespace, nil, nil, false, false, 0, false)
		}()
	} else {
		go func() {
			defer wg.Done()
			installErr = HelmInstall(releaseName, chartPath, namespace)
		}()
	}

	wg.Add(2)
	go func() {
		defer wg.Done()
		lintErr = HelmLint(chartPath)
	}()

	go func() {
		defer wg.Done()
		templateErr = HelmTemplate(releaseName, chartPath, namespace)
	}()

	wg.Wait()

	if installErr != nil || upgradeErr != nil || lintErr != nil || templateErr != nil {
		if installErr != nil {
			pterm.Error.Println("Install failed:", installErr)
		}
		if upgradeErr != nil {
			pterm.Error.Println("Upgrade failed:", upgradeErr)
		}
		if lintErr != nil {
			pterm.Error.Println("Lint failed:", lintErr)
		}
		if templateErr != nil {
			pterm.Error.Println("Template rendering failed:", templateErr)
		}
		// make this return error in red color
		return fmt.Errorf(color.RedString("Provisioning failed"))
	}

	pterm.Success.Println("Provisioning completed successfully.")
	return nil
}
