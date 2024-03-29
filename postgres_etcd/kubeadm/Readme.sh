# Migrating from etcd to Postgres for Kubernetes states and storage asas backend!

## Kine
- Kine 	us used to translate the etcd APi calls to SQL database calls, and can be used with K3s or Kubeadm for the etcd external endpoint so the clusters can store their states to the database at the endpoint.

## To setup K3s with Postgres

### Setup a Postgres server and connect it's endpoint to Kine

```bash
sudo -u postgres psql
```

- Connect using go script in the kine!

```bash
go run -x . --endpoint "postgres://user:password@127.0.0.1:5432/postgres"
```
