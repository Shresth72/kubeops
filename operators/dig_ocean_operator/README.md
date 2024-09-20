## Kluster

### shresth72.dev

v1alpha1

- apiVersion
- kind
- metadata:
- spec: #controlled by the user
  - ...
  - ...
- status: #controlled only by the controller
  - clusterId #persisted data
  - progress
  - kubeconfig
- role-test:
  - resources: Kluster/status #only allowed spec subresource

#### resource

- apis/apps/v1/namespaces/<ns>/deployments
- v1/namespaces/<ns>/pods/<podname>

#### subresources

- v1/namespaces/<ns>/pods/<podname>/logs

#### generate

1. deep copy objects
2. clientset
3. informers
4. lister

#### in production

interface {
CreateCluster
DeleteCluster
GetKubeConfig
}

DigitalOcean, AKS implement these interfaces
