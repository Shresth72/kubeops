#!/bin/bash

# Connect to the master node
curl -sfL https://get.k3s/io | K3S_NODE_NAME=k3s-worker-01 K3S_URL=https://<MASTER_IP>:6443 K3S_TOKEN=<TOKEN> sh -


