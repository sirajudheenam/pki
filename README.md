# PKI

This is created to enhance my understanding of PKI using a simple client-server app in go.

## Create Ceritifcates

```bash
bash scripts/create-certs.sh
```

script | purpose
-------| -------
[scripts/container.sh](scripts/container.sh) | to perform container operations
[scripts/create-certs.sh](scripts/create-certs.sh) | to create root CA server client certs
[scripts/kill_server.sh](scripts/kill_server.sh) | to kill when server stuck and port occupied
[scripts/run-client.sh](scripts/run-client.sh) | to run openssl client to make TLS connection
[scripts/run-server.sh](scripts/run-server.sh) | to run openssl server to accept client connections

## pki-go - Go Client Server App using TLS

## Run Server on different terminal

```bash
go run server.go
```

## Run Client on different terminal

```bash
SERVER_URL=https://localhost:8443/hello go run client.go
```

## Run client as docker container

```bash
# use our client

docker run --rm \
  -v $(pwd)/certs/client:/app/certs/client:ro \
  --network host \
  sirajudheenam/go-mtls-client

# Output 
Server response: Hello, client1!

# use public curl container image and run
docker run --rm \
  -v $(pwd)/certs/client:/certs/client:ro \
  --network host \
  curlimages/curl:latest \
  curl -vk https://host.docker.internal:8443/hello \
    --cert /certs/client/client.cert.pem \
    --key /certs/client/client.key.pem \
    --cacert /certs/client/root.cert.pem

# both works

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
docker build -f Dockerfile.client -t sirajudheenam/go-mtls-client .

docker run -d --name go-mtls-server sirajudheenam/go-mtls-client

# Server
docker tag sirajudheenam/go-mtls-server:1.0.0 sirajudheenam/go-mtls-server:latest

# Client
docker tag sirajudheenam/go-mtls-client:1.0.0 sirajudheenam/go-mtls-client:latest

docker push sirajudheenam/go-mtls-server:1.0.0
docker push sirajudheenam/go-mtls-server:latest

docker push sirajudheenam/go-mtls-client:1.0.0
docker push sirajudheenam/go-mtls-client:latest

```

## Using `docker-compose`

```bash
# Instead use docker-compose
docker-compose up --build
```

## Using `Makefile`

```bash
# Build only:
make build VERSION=1.0.0
# Build + tag + push:
make release VERSION=1.0.0
# Clean up local images
make clean VERSION=1.0.0
```

## Using `Makefile2` with docker-compose

```bash
# Build only:
make build VERSION=1.0.0
# Build + tag + push:
make release VERSION=1.0.0
# Clean up local images
make clean VERSION=1.0.0
# Start services
make up
# Stop services
make down
# Tail logs
make logs

```

## Troubleshooting

```bash
# docker network create mtls-net

# docker run -d --name go-mtls-server --network mtls-net -p 8443:8443 sirajudheenam/go-mtls-server

# docker run --rm --name go-mtls-client --network mtls-net \
#   -e SERVER_URL=https://go-mtls-server:8443/hello \
#   sirajudheenam/go-mtls-client

# 2025/09/11 01:08:09 Request failed: Get "https://go-mtls-server:8443/hello": tls: failed to verify certificate: x509: certificate is valid for localhost, not go-mtls-server

# list all docker networkd
docker network ls

# remove a network
docker network rm mtls-net

# build using docker-compose
docker compose build
```