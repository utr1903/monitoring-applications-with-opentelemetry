#!/bin/bash

# Get commandline arguments
while (( "$#" )); do
  case "$1" in
    --registry)
      containerRegistry="${2}"
      shift
      ;;
    --username)
      containerRegistryUsername="${2}"
      shift
      ;;
    *)
      shift
      ;;
  esac
done

# Container registery
if [[ $containerRegistry == "" ]]; then
  echo "Container registery [--registry] is not provided! Using default [ghcr.io]..."
  containerRegistry="ghcr.io"
fi

# Container registery username
if [[ $containerRegistryUsername == "" ]]; then
  echo "Container registery username [--username] is not provided! Using default [utr1903]..."
  containerRegistryUsername="utr1903"
fi

# Prefix
project="monitoring-otel"

# Services
otelcollector="otelcollector"
grpcserver="grpcserver"
grpcclient="grpcclient"
httpserver="httpserver"
httpclient="httpclient"

# Images
grpcserverImageName="${containerRegistry}/${containerRegistryUsername}/${project}-${grpcserver}:latest"
grpcclientImageName="${containerRegistry}/${containerRegistryUsername}/${project}-${grpcclient}:latest"
httpclientImageName="${containerRegistry}/${containerRegistryUsername}/${project}-${httpserver}:latest"
httpserverImageName="${containerRegistry}/${containerRegistryUsername}/${project}-${httpclient}:latest"

###################
### Deploy Helm ###
###################

# otelcollector
helm upgrade ${otelcollector} \
  --install \
  --wait \
  --debug \
  --set name=${otelcollector} \
  --set replicas=1 \
  "./${otelcollector}"

# grpcserver
helm upgrade ${grpcserver} \
  --install \
  --wait \
  --debug \
  --set imageName=${grpcserverImageName} \
  --set imagePullPolicy="Always" \
  --set name=${grpcserver} \
  --set replicas=1 \
  "./${grpcserver}"

# grpcclient
helm upgrade ${grpcclient} \
  --install \
  --wait \
  --debug \
  --set imageName=${grpcclientImageName} \
  --set imagePullPolicy="Always" \
  --set name=${grpcclient} \
  --set replicas=1 \
  "./${grpcclient}"
