# Setup

## Pre-requisites

The following pre-requisites are required.

- [kubectl](https://kubernetes.io/docs/tasks/tools/#kubectl)
- [Helm](https://helm.sh/)
- [Minikube](https://minikube.sigs.k8s.io/docs/start/) (or cluster of your choice)

## Create Namespaces

Create namespaces for hosting the agent operator and related demo apps.

```
kubectl create ns newrelic
kubectl create ns ao-demo
```

## Install Cert Manager

Add the Cert Manager helm repo and install Cert Manager into your cluster.

```
helm repo add jetstack https://charts.jetstack.io --force-update
helm install \
  cert-manager jetstack/cert-manager \
  --namespace cert-manager \
  --create-namespace \
  --version v1.14.4 \
  --set installCRDs=true
```

## Install the New Relic Agent Operator

Install the `newrelic-agent-operator` into your cluster using Helm.

```
helm repo add newrelic-agent-operator https://newrelic-experimental.github.io/newrelic-agent-operator
helm upgrade --install newrelic-agent-operator newrelic-agent-operator/newrelic-agent-operator --set licenseKey='xxxxxxxxxxxxxxxx' -n newrelic
```

## Create your Instrumentation custom resource

The Instrumentation custom resources referenced below may require different container images based on your CPU architecture for local testing.  Use the appropriate file based on your test system.

### arm64
```
kubectl apply -f ./demo/customresource_arm64.yaml -n ao-demo
```

### amd64
```
kubectl apply -f ./demo/customresource.yaml -n ao-demo
```


## Create license key for demo apps

A New Relic License (Ingest) Key will be stored in a secret within the `ao-demo` namespace.

```
kubectl create secret generic newrelic-key-secret -n ao-demo --from-literal=new_relic_license_key=xxxxxxxxxxxxxxxx
```

## Deploy demo apps and loadgen

The demo apps are accompanied by some basic load generation using locust.  Deploy the following files and you should be up and running quickly!

```
kubectl apply -f ./demo/apps/. -n ao-demo
kubectl apply -f ./demo/loadgen/locust.yaml -n ao-demo
```

## Validate load via NRQL

Run the following NRQL to identify that instrumentation has been successful and the applications are generating `Transaction` events.

```
FROM Transaction SELECT count(*) facet appName TIMESERIES SINCE 10 minutes ago
```