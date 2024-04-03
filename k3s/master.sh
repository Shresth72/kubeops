#!/bin/bash

# Install k3s and startup the cluster on the master node (for now we are using traefik for lb, later switch to nginx)
curl -sfL https://get.k3s.io | sh -s - --no-deploy traefik --write-kubeconfig-mode 644 --node-name k3s-master-01

# View all nodes
kubectl get nodes

# Get the token to connect worker to master
cat /var/lib/rancher/k3s/server/node-token

# Apply a nginx ingress controller pod for internet access
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v1.2.1/deploy/static/provider/cloud/deploy.yaml

# To edit the k3s service edit the file at
/etc/systemd/system/k3s.service
system restart k3s

# Apply the app deployment
kubectl apply -f manifests/deployment.yml
kubectl apply -f manifests/service.yml
kubectl apply -f manifests/ingress.yml

#
# Implementing HTTPS Certificates

# Install cert-manager
kubectl apply -f https://github.com/jetstack/cert-manager/releases/download/v1.5.3/cert-manager.yaml

