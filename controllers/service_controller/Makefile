.PHONEY: run

# Create Nginx Service and Ingress
run:
	go build && sudo ./ingress

watch-services:
	sudo watch k3s kubectl get service -n default

get-services:
	sudo k3s kubectl get service -n default

get-deploy:
	sudo k3s kubectl get deployments -n default

create-deploy:
	sudo k3s kubectl create deployment nginx -n default --image nginx

delete-deploy:
	sudo k3s kubectl delete deployment nginx -n default 

delete-service:
	sudo k3s kubectl delete services nginx

nginx-pf:
	sudo k3s kubectl port-forward -n default svc/nginx 8080:80

nginx-controller:
	sudo kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v1.11.2/deploy/static/provider/do/deploy.yaml
