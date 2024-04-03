# Ingress and Service Pods Networking

appPod <-- Service <-- Ingress <-- IngressController <-- LoadBalancer

IngressController
- Evaluate all the rules
- Manages redirections
- Entrypoint to the cluster
- Many third party implementations

LoadBalancer
- Cloud Provider Load Balancer
- External Proxy Server
