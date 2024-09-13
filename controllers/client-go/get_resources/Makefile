.PHONY: start-nginx-k3s delete-nginx-k3s deploy-nginx-k3s docker-build docker-tag docker-push create-manifests create-yaml create-deployment delete-deployment role-create role-bind view-logs

# Variables
VERSION := latest
IMAGE := nginx
CUSTOM_VERSION ?= 0.1.0
DEPLOYMENT_NAME := get_resources
NAMESPACE := default






# K3s: Start Nginx pod
start-nginx-k3s:
	sudo k3s kubectl run nginx --image=$(IMAGE):$(VERSION) --restart=Never -n $(NAMESPACE)

# K3s: Delete Nginx pod
delete-nginx-k3s:
	sudo k3s kubectl delete pod nginx -n $(NAMESPACE)

# K3s: Create Nginx deployment
deploy-nginx-k3s:
	sudo k3s kubectl create deployment nginx-deployment --image=$(IMAGE):$(VERSION) -n $(NAMESPACE)








# Docker: Build
docker-build:
	docker build -t $(DEPLOYMENT_NAME):$(CUSTOM_VERSION) .

# Docker: Tag
docker-tag:
	docker tag $(DEPLOYMENT_NAME):$(CUSTOM_VERSION) shresth144/$(DEPLOYMENT_NAME):$(CUSTOM_VERSION)

# Docker: Push
docker-push:
	docker push shresth144/$(DEPLOYMENT_NAME):$(CUSTOM_VERSION)








# Kubernetes: Create get_resources manifests
create-manifests:
	k create deployment $(DEPLOYMENT_NAME) --image shresth144/$(DEPLOYMENT_NAME):$(CUSTOM_VERSION) --dry-run=client -oyaml > manifests

# Kubernetes: Create get_resources yaml
create-yaml:
	k create deployment $(DEPLOYMENT_NAME) --image shresth144/$(DEPLOYMENT_NAME):$(CUSTOM_VERSION) --dry-run=client -oyaml > $(DEPLOYMENT_NAME).yaml

# Kubernetes: Create get_resources deployment
create-deployment:
	k create -f $(DEPLOYMENT_NAME).yaml

# Kubernetes: Delete get_resources deployment
delete-deployment:
	k delete -f $(DEPLOYMENT_NAME).yaml

# Kubernetes: Create Role for List Pod & Deployment
role-create:
	k create role poddepl --resource pods,deployment --verb list --namespace $(NAMESPACE)

# Kubernetes: Create Role Binding for the Default Servve Account
role-bind:
	k create rolebinding poddepl --role poddepl --serviceaccount default:default --namespace $(NAMESPACE)
	
# Kubernetes: View Deployment logs
view-logs:
	@echo "Fetching logs for the latest pod..."
	@POD_NAME=$$(kubectl get pods -n "$(NAMESPACE)" -l app="$(DEPLOYMENT_NAME)" -o jsonpath='{.items[0].metadata.name}'); \
	if [ "$$POD_NAME" = "" ]; then \
		echo "No pods found for deployment $(DEPLOYMENT_NAME)"; \
	else \
		kubectl logs -f $$POD_NAME -n "$(NAMESPACE)"; \
	fi










# Kubernetes: Start Nginx pod
start-nginx-k8s:
	kubectl run nginx --image=$(IMAGE) --restart=Never -n $(NAMESPACE)

# Kubernetes: Delete Nginx pod
delete-nginx-k8s:
	kubectl delete pod nginx -n $(NAMESPACE)

