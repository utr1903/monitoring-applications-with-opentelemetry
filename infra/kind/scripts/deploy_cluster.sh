#!/bin/bash

###################
### Infra Setup ###
###################

kind create cluster \
  --name test \
  --config ../config/kind-config.yaml \
  --image=kindest/node:v1.27.3
