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
    --platform)
      platform="$2"
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

# Container platform
if [[ $platform == "" ]]; then
  # Default is amd64
  platform="amd64"
else
  if [[ $platform != "amd64" && $platform != "arm64" ]]; then
    echo "Platform [--platform] can either be 'amd64' or 'arm64'."
    exit 1
  fi
fi

# Prefix
project="monitoring-otel"

# Services
grpcserver="grpcserver"
grpcclient="grpcclient"
httpserver="httpserver"
httpclient="httpclient"

# Images
grpcserverImageName="${containerRegistry}/${containerRegistryUsername}/${project}-${grpcserver}:latest"
grpcclientImageName="${containerRegistry}/${containerRegistryUsername}/${project}-${grpcclient}:latest"
httpclientImageName="${containerRegistry}/${containerRegistryUsername}/${project}-${httpserver}:latest"
httpserverImageName="${containerRegistry}/${containerRegistryUsername}/${project}-${httpclient}:latest"

####################
### Build & Push ###
####################

# grpcserver
docker build \
  --platform "linux/${platform}" \
  --tag "${grpcserverImageName}" \
  --build-arg="APP_NAME=${grpcserver}" \
  "./."
docker push "${grpcserverImageName}"

# grpcclient
docker build \
  --platform "linux/${platform}" \
  --tag "${grpcclientImageName}" \
  --build-arg="APP_NAME=${grpcclient}" \
  "./."
docker push "${grpcclientImageName}"

# # httpserver
# docker build \
#   --platform "linux/${platform}" \
#   --tag "${httpserverImageName}" \
#   --build-arg="APP_NAME=${httpserver}" \
#   "./."
# docker push "${httpserverImageName}"

# # httpclient
# docker build \
#   --platform "linux/${platform}" \
#   --tag "${httpclientImageName}" \
#   --build-arg="APP_NAME=${httpclient}" \
#   "./."
# docker push "${httpclientImageName}"
