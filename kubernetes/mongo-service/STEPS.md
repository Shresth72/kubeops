# Steps to start and deploy the k8s services

## Minikube start-up

- Start the minikube cluster env to connect using kubectl

```bash
minikube start
```

- Now you can interact with the pods and containers using kubectl and view all the pods and their details

```bash
kubectl get pod

kubectl get pod -o wide
```

- Create all the resources such as the secrets and configs, inside the minikube cluster

```bash
kubectl apply -f secret.yaml
kubectl apply -f mongo-config.yaml
kubectl apply -f mongo-app.yaml
kubectl apply -f web-app.yaml
```

- To view all the secrets, config maps and services

```bash
kubectl get secret
kubectl get configmap
kubectl get svc
```

- To expose a service route IP

```bash
minikube service webapp-service
```

- To remove resources

```bash
kubectl delete deployment --all
kubectl delete service --all
kubectl delete secret --all
kubectl delete configmap --all
```