#!/bin/bash

# To run kine
go run -x . --endpoint "postgres://shres:secret@127.0.0.1:5433/kubernetes"

# Get k3s and connect to the postgres instance
curl -sfL https://get.k3s.io | sh -s - server --write-kubeconfig-mode=644 \
	--token=sample_secret_token \
	--datastore-endpoint="postgres://shres:secret@127.0.0.1:5433/kubernetes"

# Get nodes
k3s kubectl get nodes

