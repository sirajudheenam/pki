#!/bin/bash
VERSION=$(git describe --tags --always) # or use `date +%Y.%m.%d`

docker build -f Dockerfile.server -t sirajudheenam/go-mtls-server:$VERSION .
docker build -f Dockerfile.client -t sirajudheenam/go-mtls-client:$VERSION .

docker compose up
