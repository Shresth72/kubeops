#!/bin/bash

# Connect to the scalable Postgres Instance on Yugabyte
psql -h 10.168.0.0 -p 5433 -U yugabyte

# Now start a k3s cluster with endpoint as the Yugabyte postgres instance
curl -sfL https://get.k3s.io | sh -s - server --write-kubeconfig-mode=644 \
	--token=sample_secret_token \
	--datastore-endpoint="postgres://yugabyte:password@10.168.0.0:5433/yugabyte"

# Apply services
k3s kubectl apply -k ./app/service.yaml

