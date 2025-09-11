# PKI

This Repo is a side-effect of my pursuit to enhance my understanding of PKI using a simple client-server app in go.

We shall create certificates manually first to understand what does it do.


```bash
go mod init github.com/sirajudheenam/pki
go mod tidy
```
## Create Ceritifcates

```bash
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
go run server.go
```

## Run Client on different terminal

```bash
SERVER_URL=https://localhost:8443/hello go run client.go
```

## Run tests
```bash
go test ./internal/... -v
```

## Create Docker Image

```bash

# Build image
docker build -f Dockerfile.server -t sirajudheenam/go-mtls-server .

# Run container
docker run -d --name go-mtls-server -p 8443:8443 sirajudheenam/go-mtls-server

# Run container
# docker run --rm -p 8443:8443 go-mtls-server 

# Build image
docker build -f Dockerfile.client -t sirajudheenam/go-mtls-client:1.0.0 -t sirajudheenam/go-mtls-client:latest .

# Tag your explicitely with different tag if needed, though above step does that.
# Server
docker tag sirajudheenam/go-mtls-server:1.0.0 sirajudheenam/go-mtls-server:latest
# Client
docker tag sirajudheenam/go-mtls-client:1.0.0 sirajudheenam/go-mtls-client:latest

docker push sirajudheenam/go-mtls-server:1.0.0
docker push sirajudheenam/go-mtls-server:latest
docker push sirajudheenam/go-mtls-client:1.0.0
docker push sirajudheenam/go-mtls-client:latest

```

## Run server in a docker container

```bash
docker run --rm \
  -v $(pwd)/certs/server:/app/certs/server:ro \
  --network host \
  sirajudheenam/go-mtls-server
```

```bash
## Run client as docker container

docker run --rm \
  -v $(pwd)/certs/client:/app/certs/client:ro \
  --network host \
  sirajudheenam/go-mtls-client

# Output 
Server response: Hello, client1!

# use public curl container image and run
# below needs /etc/host entry on macOS 'localhost host.docker.internal'
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

## Using `docker-compose`

```bash
# Instead use docker-compose
docker-compose up --build
```

## Using `Makefile`

```bash
# Build only:
make build VERSION=1.0.1
# Build + tag + push:
make release VERSION=1.0.1
# Clean up local images
make clean VERSION=1.0.1
```

## Using `Makefile2` with docker-compose

```bash
# Build only:
make build VERSION=1.0.1
# Build + tag + push:
make release VERSION=1.0.1
# Clean up local images
make clean VERSION=1.0.1
# Start services
make up
# Stop services
make down
# Tail logs
make logs

```

## Deploy it in k8s

```bash
cd pki-go/certs/server
kubectl create secret generic go-mtls-server-certs \
  --from-file=server.chain.pem \
  --from-file=server.key.pem \
  --from-file=root.cert.pem \
  --from-file=intermediate.cert.pem
# Output
# secret/go-mtls-server-certs created

cd ../../../k8s
kubectl create -f pki-server.yaml

# Output
# deployment.apps/go-mtls-server created
# service/go-mtls-server-service created
# ingress.networking.k8s.io/go-mtls-server-ingress created


# on minikube setup on macOS with Docker driver
# run this on separate terminal
kubectl port-forward svc/go-mtls-server-service 8443:8443

# Output
# Forwarding from 127.0.0.1:8443 -> 8443
# Forwarding from [::1]:8443 -> 8443
# Handling connection for 8443

# run from macOS

cd ../pki-go
SERVER_URL=https://go-mtls-server:8443 go run cmd/client/main.go

# make sure you have host entry at /etc/hosts as below

# check with cat /etc/hosts | grep go-mtls-server
# 127.0.0.1   go-mtls-server
# else 
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


# # client
cd pki-go
kubectl create secret generic go-mtls-client-certs \
  --from-file=client.cert.pem=./certs/client/client.cert.pem \
  --from-file=client.key.pem=./certs/client/client.key.pem \
  --from-file=inter-root-combined.cert.pem=./certs/client/inter-root-combined.cert.pem \
  --from-file=root.cert.pem=../demo-pki/root/root.cert.pem

# Output
# secret/go-mtls-client-certs created

cd ../k8s
kubectl create -f .

# Output
# deployment.apps/go-mtls-client created
# job.batch/go-mtls-client-job created

```

We have successfully deployed the app in Docker and k8s

What is next?

Perhaps create a helm chart for it.