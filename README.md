# planton-cloud-kube-agent

Planton Cloud Agent that runs on customers' kubernetes clusters. This agent acts as a proxy between the customers'
kubernetes clusters and Planton Cloud Service.

## local development

### setup environment variables

```shell
rm -f .env_export
cat > .env_export << EOF
export LOG_LEVEL=debug
export PLANTON_CLOUD_KUBE_AGENT_MACHINE_ACCOUNT_EMAIL=machine-account-email
export PLANTON_CLOUD_KUBE_AGENT_CLIENT_SECRET=client-secret
export PLANTON_CLOUD_KUBE_AGENT_HOSTING_CLUSTER_ID=ho-planton-gcp-as1-host-a
export PLANTON_CLOUD_SERVICE_API_ENDPOINT=api.dev.planton.cloud:443
export OPEN_COST_API_ENDPOINT=http://localhost:9003
export OPEN_COST_POLLING_INTERVAL_SECONDS=10
export TOKEN_EXPIRATION_BUFFER_MINUTES=5
export TOKEN_EXPIRATION_CHECK_INTERVAL_SECONDS=120
EOF
```

### start kubecost port-forward

```shell
kubectl port-forward -n kubecost svc/kubecost-cost-analyzer 9003:9003
```

### start microservice

```shell
make run
```
