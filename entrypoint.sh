#!/bin/sh
set -e

# Log in to Docker Hub if credentials are provided
if [ -n "$DOCKER_USERNAME" ] && [ -n "$DOCKER_PASSWORD" ]; then
  echo "$DOCKER_PASSWORD" | docker login --username "$DOCKER_USERNAME" --password-stdin
  echo "Logged into Docker Hub successfully."
fi

# # Log in to GitHub Container Registry if credentials are provided
# if [ -n "$GITHUB_USERNAME" ] && [ -n "$TOKEN" ]; then
#   echo "$TOKEN" | docker login ghcr.io -u "$GITHUB_USERNAME" --password-stdin
#   echo "Logged into GitHub Container Registry successfully."
# fi

# Execute the smurf application with provided commands
exec "/go/src/app/smurf" "$@"












