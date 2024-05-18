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

# Monitoring
prometheus="prometheus"
tempo="tempo"
loki="loki"
grafana="grafana"
otelcollector="otelcollector"

# Services
grpcserver="grpcserver"
grpcclient="grpcclient"
httpserver="httpserver"
httpclient="httpclient"

# Images
grpcserverImageName="${containerRegistry}/${containerRegistryUsername}/${project}-${grpcserver}:latest"
grpcclientImageName="${containerRegistry}/${containerRegistryUsername}/${project}-${grpcclient}:latest"
httpserverImageName="${containerRegistry}/${containerRegistryUsername}/${project}-${httpserver}:latest"
httpclientImageName="${containerRegistry}/${containerRegistryUsername}/${project}-${httpclient}:latest"

###################
### Deploy Helm ###
###################

helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo add grafana https://grafana.github.io/helm-charts
helm repo update

# prometheus
helm upgrade ${prometheus} \
  --install \
  --wait \
  --debug \
  --set alertmanager.enabled=false \
  --set kube-state-metrics.enabled=false \
  --set prometheus-node-exporter.enabled=false \
  --set prometheus-pushgateway.enabled=false \
  --values ./${prometheus}/values.yaml \
  --version "25.21.0" \
  "prometheus-community/prometheus"

# tempo
helm upgrade ${tempo} \
  --install \
  --wait \
  --debug \
  --set tempo.metricsGenerator.enabled=true \
  --set tempo.metricsGenerator.remoteWriteUrl="http://prometheus-server.default.svc.cluster.local:80/api/v1/write" \
  --version "1.7.3" \
  "grafana/tempo"

# loki
helm upgrade ${loki} \
  --install \
  --wait \
  --debug \
  --values ./${loki}/values.yaml \
  --version "6.5.2" \
  "grafana/loki"

# grafana
helm upgrade ${grafana} \
  --install \
  --wait \
  --debug \
  --set adminUser=admin \
  --set adminPassword=admin123 \
  --values ./${grafana}/values.yaml \
  --version "7.3.11" \
  "grafana/grafana"

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

# httpserver
helm upgrade ${httpserver} \
  --install \
  --wait \
  --debug \
  --set imageName=${httpserverImageName} \
  --set imagePullPolicy="Always" \
  --set name=${httpserver} \
  --set replicas=1 \
  "./${httpserver}"

# httpclient
helm upgrade ${httpclient} \
  --install \
  --wait \
  --debug \
  --set imageName=${httpclientImageName} \
  --set imagePullPolicy="Always" \
  --set name=${httpclient} \
  --set replicas=1 \
  "./${httpclient}"
