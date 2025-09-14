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
  "${K8s_DIR}/pki-client-jobs.yaml" \
  "${K8s_DIR}/pki-client-deployment.yaml"; do

  if kubectl get -f "$file" &> /dev/null; then
    echo "Deleting resources from $file"
    kubectl delete -f "$file"
  else
    echo "No resources found in $file, skipping"
  fi
done

# Create the server secret in Kubernetes
kubectl create secret generic go-mtls-server-certs \
  --from-file=server.chain.pem=${PKI_GO_DIR}/certs/server/server.chain.pem \
  --from-file=server.key.pem=${PKI_GO_DIR}/certs/server/server.key.pem \
  --from-file=root.cert.pem=${PKI_GO_DIR}/certs/server/root.cert.pem \
  --from-file=intermediate.cert.pem=${PKI_GO_DIR}/certs/server/intermediate.cert.pem

# # Create the client secret in Kubernetes
kubectl create secret generic go-mtls-client-certs \
  --from-file=client.cert.pem=${PKI_GO_DIR}/certs/client/client.cert.pem \
  --from-file=client.key.pem=${PKI_GO_DIR}/certs/client/client.key.pem \
  --from-file=root.cert.pem=${PKI_GO_DIR}/certs/server/root.cert.pem \
  --from-file=inter-root-combined.cert.pem=${PKI_GO_DIR}/certs/client/inter-root-combined.cert.pem

# # Deploy the server application

kubectl apply -f ${K8s_DIR}/pki-server-deployment.yaml
kubectl apply -f ${K8s_DIR}/pki-server-svc.yaml
kubectl apply -f ${K8s_DIR}/pki-client-jobs.yaml

sleep 5

# To verify the deployment, you can check the status of the pods
kubectl get pods

# Wait for job pod to complete
JOB_NAME="go-mtls-client-job"

echo "Waiting for Job $JOB_NAME to complete..."
kubectl wait --for=condition=complete job/$JOB_NAME --timeout=120s

# Get the name of the pod created by the Job
POD_NAME=$(kubectl get pods --selector=job-name=$JOB_NAME \
  --output=jsonpath='{.items[0].metadata.name}')

# Show status and logs
echo "===== Pod Status ====="
kubectl get pod "$POD_NAME" -o wide

echo "===== Pod Logs ====="
kubectl logs "$POD_NAME"

# Check CronJobs

# Get the latest job name created by the cronjob
LATEST_JOB=$(kubectl get jobs \
  --selector=job-name \
  --sort-by=.metadata.creationTimestamp \
  -o jsonpath='{.items[-1].metadata.name}')

# Get its pod and logs
POD_NAME=$(kubectl get pods --selector=job-name=$LATEST_JOB \
  -o jsonpath='{.items[0].metadata.name}')

kubectl logs $POD_NAME
