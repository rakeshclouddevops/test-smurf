# #!/bin/bash
# set -e

# # Log in to Docker Hub if credentials are provided
# if [ -n "$DOCKER_USERNAME" ] && [ -n "$DOCKER_PASSWORD" ]; then
#   echo "$DOCKER_PASSWORD" | docker login --username "$DOCKER_USERNAME" --password-stdin
#   echo "Logged into Docker Hub successfully."
# fi

# # Execute the smurf application with provided commands
# exec "/usr/local/bin/smurf" "$@"

#!/bin/bash

set -e  # Exit on error

# Read inputs from the GitHub Actions workflow
KUBECONFIG_PATH=$1
HELM_CHART=$2
RELEASE_NAME=$3
NAMESPACE=$4
SMURF_VERSION=$5
ACTION=$6

echo "Using Kubeconfig: $KUBECONFIG_PATH"
echo "Using Helm chart: $HELM_CHART"
echo "Release Name: $RELEASE_NAME"
echo "Namespace: $NAMESPACE"
echo "Smurf Selm Version: $SMURF_VERSION"
echo "Action: $ACTION"

# Set up kubeconfig
export KUBECONFIG=$KUBECONFIG_PATH

# Download the specific version of smurf selm if specified
if [ "$SMURF_VERSION" != "latest" ]; then
    curl -sSL https://github.com/smurf/selm/releases/download/v$SMURF_VERSION/smurf-selm-linux-amd64 -o smurf-selm && \
    chmod +x smurf-selm && \
    mv smurf-selm /usr/local/bin/
else
    echo "Using the latest version of smurf selm"
fi

# Run the corresponding smurf selm action
case $ACTION in
  create)
    echo "Creating Helm chart at $HELM_CHART"
    smurf selm create --chart-dir "$HELM_CHART"
    ;;
  install)
    echo "Installing Helm chart at $HELM_CHART"
    smurf selm install --chart-dir "$HELM_CHART" --namespace "$NAMESPACE"
    ;;
  upgrade)
    echo "Upgrading Helm chart at $HELM_CHART"
    smurf selm upgrade --chart-dir "$HELM_CHART" --namespace "$NAMESPACE"
    ;;
  provision)
    echo "Provisioning Helm chart at $HELM_CHART"
    smurf selm provision --chart-dir "$HELM_CHART" --namespace "$NAMESPACE"
    ;;
  *)
    echo "Invalid action: $ACTION"
    exit 1
    ;;
esac

echo "Smurf Selm action completed successfully!"
