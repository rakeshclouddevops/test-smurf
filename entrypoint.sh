#!/bin/bash
set -e

# Log in to Docker Hub if credentials are provided
if [ -n "$DOCKER_USERNAME" ] && [ -n "$DOCKER_PASSWORD" ]; then
  echo "$DOCKER_PASSWORD" | docker login --username "$DOCKER_USERNAME" --password-stdin
  echo "Logged into Docker Hub successfully."
fi

# Execute the smurf application with provided commands
exec "/usr/local/bin/smurf" "$@"