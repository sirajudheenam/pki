#!/bin/bash
PKI_DIR=pki-go
case "$1" in
  build-server)
    docker build -f $PKI_DIR/Dockerfile.server -t sirajudheenam/go-mtls-server .
    ;;
  build-client)
    docker build -f $PKI_DIR/Dockerfile.client -t sirajudheenam/go-mtls-client .
    ;;
  run-server)
    # docker run --rm -p 8443:8443 sirajudheenam/go-mtls-server 
    docker run -d --name go-mtls-server -p 8443:8443 sirajudheenam/go-mtls-server:1.0.0
    ;;
  run-client)
    docker run --rm sirajudheenam/go-mtls-client:1.0.0

    ;;
  cleanup)
    docker ps -a --filter "ancestor=pki-go-server" -q | xargs -r docker rm -f
    docker ps -a --filter "ancestor=pki-go-client" -q | xargs -r docker rm -f
    docker ps -a --filter "ancestor=sirajudheenam/go-mtls-server:1.0.0" -q | xargs -r docker rm -f
    docker ps -a --filter "ancestor=sirajudheenam/go-mtls-client:1.0.0" -q | xargs -r docker rm -f
    # docker rmi -f $(docker images -f "dangling=true" -q)
    # docker rmi -f sirajudheenam/go-mtls-server sirajudheenam/go-mtls-client
    ;;
  rmi)
    docker rmi -f sirajudheenam/go-mtls-server:1.0.0 sirajudheenam/go-mtls-client:1.0.0
    ;;
  *)
    echo "Usage: $0 [build-server|build-client|run-server|run-client|cleanup|rmi]"
    ;;
esac


# docker build -f Dockerfile.server -t sirajudheenam/go-mtls-server .
# docker build -f Dockerfile.client -t sirajudheenam/go-mtls-client .