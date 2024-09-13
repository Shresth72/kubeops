.PHONY: start-nginx-k3s delete-nginx-k3s start-nginx-k8s delete-nginx-k8s

# K3s: Start Nginx pod
start-nginx-k3s:
	sudo k3s kubectl run nginx --image=nginx:latest --restart=Never

# K3s: Delete Nginx pod
delete-nginx-k3s:
	sudo k3s kubectl delete pod nginx

# K3s: Create Nginx deployment
deploy-nginx-k3s:
	sudo k3s kubectl create deployment nginx-deployment --image=nginx:latest










# Kubernetes: Start Nginx pod
start-nginx-k8s:
	kubectl run nginx --image=nginx --restart=Never

# Kubernetes: Delete Nginx pod
delete-nginx-k8s:
	kubectl delete pod nginx

