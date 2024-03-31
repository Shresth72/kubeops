#!/bin/bash

# Build the kubeadm cluster
kubeadm init --config=/kubeadm-config.yaml --upload-certs \
ignore-preflight-errors ExternalEtcdVersion 2>&1 || true

kubeadm taint nodes -all node-role.kubernetes.io/control-plane-