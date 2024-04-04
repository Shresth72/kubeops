# HELM

Package Manager for Kubernetes to package YAML Files and distribute them in public and private repositories

## Helm Chart

- Bundle of YAML Files
- Push to Helm Repository
- Install on Chart and reuse it

Also works as a templating engine for Kubernetes YAML Files to make them reusable as blueprints and avoid duplication of similar YAML Files
The files can be parameterized and values can be passed to the templates to make them reusable. Values defined either via `yaml` file or with `--set` flag

Useful in CI/CD pipelines to deploy applications on Kubernetes Clusters

```bash
mychart/
    Chart.yaml
    deafult_values.yaml
    charts/
    templates/
        deployment.yaml
        service.yaml
        ingress.yaml
```

## Release Management

- Helm Versioning is called `Release` and can be managed with `helm list` command

- Helm Version 2 comes in two parts: `Tiller` and `Helm Client`. Tiller is the server side component that runs on the Kubernetes Cluster and interacts with the Kubernetes API Server. Helm Client is the client side component that runs on the local machine and interacts with the Tiller Server.

- `Tiller` stores the copy of the configuration that `Client` sets, and creates a history of charts execution. It also stores the release information in the `ConfigMap` in the `kube-system` namespace.

- This way the changes are applied to the existing deployment instead of creating a new one. This is useful for rollback and versioning.

```bash
helm install mychart
helm upgrade mychart
```

- Helm Version 3 removes the `Tiller` component and stores the release information in the `ConfigMap` in the `kube-system` namespace. This makes the Helm more secure because the `Tiller` component was a security risk.
