#!/bin/bash

set -e

# Get the directory where the script itself is located
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# Go one directory up
PARENT_DIR="$(dirname "$SCRIPT_DIR")"
PKI_GO_DIR="${PARENT_DIR}/pki-go"
K8s_DIR="${PARENT_DIR}/k8s"

echo "Deleting existing Kubernetes resources if they exist..."

# Function to delete resource only if it exists
delete_if_exists() {
  local resource_type=$1
  shift
  for name in "$@"; do
    if kubectl get "$resource_type" "$name" &> /dev/null; then
      echo "Deleting $resource_type/$name"
      kubectl delete "$resource_type" "$name"
    else
      echo "$resource_type/$name not found, skipping"
    fi
  done
}

# Delete secrets
delete_if_exists secret go-mtls-server-certs go-mtls-client-certs

# Delete resources from YAML files
for file in \
  "${K8s_DIR}/pki-server-deployment.yaml" \
  "${K8s_DIR}/pki-server-svc.yaml" \
  "${K8s_DIR}/pki-server-ingress.yaml" \
  "${K8s_DIR}/pki-client-jobs.yaml" \
  "${K8s_DIR}/pki-client-deployment.yaml"; do

  if kubectl get -f "$file" &> /dev/null; then
    echo "Deleting resources from $file"
    kubectl delete -f "$file"
  else
    echo "No resources found in $file, skipping"
  fi
done