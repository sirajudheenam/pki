# PKI

This Repo is a side-effect of my pursuit to enhance my understanding of PKI using a simple client-server app in go.

We shall create certificates manually first to understand what does it do.


## Create Ceritifcates for local development

```bash
# This creates certificates for `localhost`
bash scripts/create-certs.sh
```

script | purpose
-------| -------
[scripts/container.sh](scripts/container.sh) | to perform container operations
[scripts/create-certs.sh](scripts/create-certs.sh) | to create root CA server client certs
[scripts/create-certs-k8s.sh](scripts/create-certs-k8s.sh) | to create root CA server client certs for k8s containers with hostname other than localhost.
[scripts/kill_server.sh](scripts/kill_server.sh) | to kill when server stuck and port occupied
[scripts/run-client.sh](scripts/run-client.sh) | to run openssl client to make TLS connection
[scripts/run-server.sh](scripts/run-server.sh) | to run openssl server to accept client connections

## pki-go - Go Client Server App using TLS

```bash
cd pki-go
go mod init github.com/sirajudheenam/pki/pki-go
go mod tidy
```

## Run Server on different terminal

```bash
go run cmd/server/main.go
```

## Run Client on different terminal

```bash
# watch here SERVER_NAME is localhost
SERVER_NAME=localhost SERVER_PORT=8443 SERVER_ROOT_PATH=/hello CLIENT_CERTS=./certs/client go run cmd/client/main.go
```

## Run tests

```bash
go test ./internal/... -v

# Output 
# ok      github.com/sirajudheenam/pki/pki-go/internal/client     0.705s
# ok      github.com/sirajudheenam/pki/pki-go/internal/server     0.878s

```

## Create Docker Image for Server

```bash

# Build image
docker build -f Dockerfile.server -t sirajudheenam/go-mtls-server:1.0.0 -t sirajudheenam/go-mtls-server:local .

# Run container
docker run --rm --name go-mtls-server \
  --network host -v $(pwd)/certs/server:/app/certs/server:ro \
  -p 8443:8443 sirajudheenam/go-mtls-server:local

# Tag your explicitely with different tag if needed, though above step does that.
docker tag sirajudheenam/go-mtls-server:1.0.0 sirajudheenam/go-mtls-server:local

# Push Container
docker push sirajudheenam/go-mtls-server:1.0.0
docker push sirajudheenam/go-mtls-server:local

```

## Run server in a docker container locally

```bash
docker run --rm \
  -v $(pwd)/certs/server:/app/certs/server:ro \
  --network host \
  sirajudheenam/go-mtls-server
```

## Create Docker Image for Client

```bash
# Build image
docker build -f Dockerfile.client -t sirajudheenam/go-mtls-client:local .
# Tag your explicitely with different tag if needed, though above step does that.
# docker tag sirajudheenam/go-mtls-client:1.0.0 sirajudheenam/go-mtls-client:local

# since this is local build, there is no need to push
# docker push sirajudheenam/go-mtls-client:local

```

## Run client as docker container to test locally

```bash
docker run --rm \
  -v $(pwd)/certs/client:/app/certs/client:ro \
  --network host \
  -e SERVER_NAME=localhost -e SERVER_PORT=8443 \
  -e SERVER_ROOT_PATH=/hello -e CLIENT_CERTS=./certs/client/ \
  sirajudheenam/go-mtls-client:local

# Output 
# Server response: Hello, client1!
```

## Test public curl container image and run

```bash
# below needs /etc/host entry on macOS 'localhost host.docker.internal'
echo "localhost   host.docker.internal" | sudo tee -a /etc/hosts

docker run --rm \
  -v $(pwd)/certs/client:/certs/client:ro \
  --network host \
  curlimages/curl:latest \
  curl -vk https://host.docker.internal:8443/hello \
    --cert /certs/client/client.cert.pem \
    --key /certs/client/client.key.pem \
    --cacert /certs/client/root.cert.pem

docker run --rm \
  -v $(pwd)/certs/client:/certs/client:ro \
  --network host \
  curlimages/curl:latest \
  curl -vk https://localhost:8443/hello \
    --cert /certs/client/client.cert.pem \
    --key /certs/client/client.key.pem \
    --cacert /certs/client/root.cert.pem

# both works

```

## Using `Makefile` with `docker-compose`

```bash
# Build only:
make build VERSION=local
# Build + tag + push:
make release VERSION=local
# Clean up local images
make clean VERSION=local
# Start services
make up
# Stop services
make down
# Tail logs
make logs

```

## Deploy in kubernetes cluster

### Prerequisites

- minikube running k8s locally
- access to any cloud

In order to run this on minikube or Cloud, let us use `go-mtls-server` as `SERVER_NAME`

### Generate Certificates with `go-mtls-server` as hostname

```bash

bash scripts/create-certs-k8s.sh
```

### Test code locally

```bash
cd pki-go

# don't close the terminal where this is running
go run cmd/server/main.go

# on another terminal 
SERVER_NAME=go-mtls-server SERVER_PORT=8443 SERVER_ROOT_PATH=/hello go run cmd/client/main.go

```

### Build and push docker images after testing locally

```bash
# Build image
docker build -f Dockerfile.client -t sirajudheenam/go-mtls-client:1.0.0 -t sirajudheenam/go-mtls-client:latest .
# Tag your explicitely with different tag if needed, though above step does that.
# docker tag sirajudheenam/go-mtls-client:1.0.0 sirajudheenam/go-mtls-client:latest

# since this is local build, there is no need to push
docker push sirajudheenam/go-mtls-client:1.0.0
docker push sirajudheenam/go-mtls-client:latest

```

### Deploy it k8s

```bash
# # server
cd pki-go
# delete if you have secret
# kubectl delete secret go-mtls-server-certs
# kubectl delete secret go-mtls-client-certs
kubectl create secret generic go-mtls-server-certs \
  --from-file=server.chain.pem=./certs/server/server.chain.pem \
  --from-file=server.key.pem=./certs/server/server.key.pem \
  --from-file=root.cert.pem=./certs/server/root.cert.pem \
  --from-file=intermediate.cert.pem=./certs/server/intermediate.cert.pem
# Output
# secret/go-mtls-server-certs created

kubectl create secret generic go-mtls-client-certs \
  --from-file=client.cert.pem=./certs/client/client.cert.pem \
  --from-file=client.key.pem=./certs/client/client.key.pem \
  --from-file=inter-root-combined.cert.pem=./certs/client/inter-root-combined.cert.pem \
  --from-file=root.cert.pem=./certs/server/root.cert.pem

# Output
# secret/go-mtls-client-certs created

cd ../k8s
# clean up first 
kubectl delete -f .
# create 
kubectl create -f .

# Output


# on minikube setup on macOS with Docker driver
# run this on separate terminal
kubectl port-forward svc/go-mtls-server-service 8443:8443

# Output
# Forwarding from 127.0.0.1:8443 -> 8443
# Forwarding from [::1]:8443 -> 8443
# Handling connection for 8443


# make sure you have host entry at /etc/hosts as below

# check with 
cat /etc/hosts | grep go-mtls-server
# 127.0.0.1   go-mtls-server
# if doesn't exist run
echo "127.0.0.1   go-mtls-server" | sudo tee -a /etc/hosts


# Output
# Server response: Hello, client1!

curl -vk https://localhost:8443/hello \
  --cert ./certs/client/client.cert.pem \
  --key ./certs/client/client.key.pem \
  --cacert ./certs/client/root.cert.pem

# Output 
# Hello, client1!

# Hurray!!  It works !

cd pki-go
docker run --rm \
  -v $(pwd)/certs/client:/certs/client:ro \
  curlimages/curl:latest \
  curl -vk https://localhost:8443/hello \
    --cert /certs/client/client.cert.pem \
    --key /certs/client/client.key.pem \
    --cacert /certs/client/root.cert.pem

docker run --rm \
  -v $(pwd)/certs/client:/certs/client:ro \
  curlimages/curl:latest \
  curl -vk https://go-mtls-server:8443/hello \
    --cert /certs/client/client.cert.pem \
    --key /certs/client/client.key.pem \
    --cacert /certs/client/root.cert.pem


cd ../k8s
kubectl delete -f pki-client-deployment.yaml 

kubectl create -f pki-client-deployment.yaml 

# Output
# deployment.apps/go-mtls-client created
# job.batch/go-mtls-client-job created

```

We have successfully deployed the app in Docker and k8s

What is next?

Perhaps create a helm chart for it.

