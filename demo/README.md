# Setup

## Create Namespaces
```
kubectl create ns newrelic
kubectl create ns ao-demo
```

## Install Cert Manager
```
helm install \
  cert-manager jetstack/cert-manager \
  --namespace cert-manager \
  --create-namespace \
  --version v1.11.0 \
  --set installCRDs=true
```

## Install `newrelic-agent-operator`
```
helm upgrade --install newrelic-agent-operator chart/ \
  --set licenseKey='xxxxxxxxxxxxxxxx' -n newrelic
```

## Create custom resource

```
kubectl apply -f customresource.yaml -n ao-demo
```

## Create license key for demo apps
```
kubectl create secret generic newrelic-key-secret -n ao-demo --from-literal=new_relic_license_key=xxxxxxxxxxxxxxxx
```

## Deploy demo apps

```
kubectl apply -f ./apps/. -n ao-demo
```