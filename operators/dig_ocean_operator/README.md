Kluster

shresth72.dev
v1alpha1

generate

1. deep copy objects
2. clientset
3. informers
4. lister

#### In Production

interface {
CreateCluster
DeleteCluster
GetKubeConfig
}

DigitalOcean, AKS implement these interfaces
