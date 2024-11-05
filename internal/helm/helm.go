package helm

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pterm/pterm"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/chartutil"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/release"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
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

func HelmUpgrade(releaseName, chartPath, namespace string) error {
	spinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Upgrading release '%s'", releaseName))
	actionConfig := new(action.Configuration)
	if err := actionConfig.Init(settings.RESTClientGetter(), namespace, "secrets", nil); err != nil {
		spinner.Fail(err.Error())
		return err
	}
	client := action.NewUpgrade(actionConfig)
	client.Namespace = namespace

	chart, err := loader.Load(chartPath)
	if err != nil {
		spinner.Fail(err.Error())
		return err
	}

	_, err = client.Run(releaseName, chart, nil)
	if err != nil {
		spinner.Fail(err.Error())
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
	if err := actionConfig.Init(settings.RESTClientGetter(), namespace, "secrets", nil); err != nil {
		spinner.Fail(err.Error())
		return err
	}
	client := action.NewUninstall(actionConfig)

	_, err := client.Run(releaseName)
	if err != nil {
		spinner.Fail(err.Error())
	} else {
		spinner.Success(fmt.Sprintf("Release '%s' uninstalled successfully", releaseName))
	}
	return err
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
