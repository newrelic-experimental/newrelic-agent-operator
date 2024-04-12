# New Relic Agent Operator Installation and Configuration

## [Prerequisite] Install Cert Manager

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
helm upgrade --install newrelic-agent-operator newrelic-agent-operator/newrelic-agent-operator --set licenseKey='<NEW RELIC INGEST LICENSE KEY>' -n newrelic
```

## Configure Auto-Instrumentation

The `Instrumentation` custom resource defines the New Relic instrumentation containers that are injected as init containers into your application pods.  The instrumentation container versions correspond with New Relic APM agent versions so you have full control over which agent versions are deployed.  

**Important:** Each namespace in your cluster will need an `Instrumentation` resource.

You can configure the agent settings globally in this file using ENV variables.  Global agent settings can be overridden in your deployment manifest if a different configuration is required **TODO**: add example of this.

```
apiVersion: newrelic.com/v1alpha1
kind: Instrumentation
metadata:
  labels:
    app.kubernetes.io/name: instrumentation
    app.kubernetes.io/created-by: newrelic-agent-operator
  name: newrelic-instrumentation
spec:
  # ### Start OTel collector config ###
  # # A required but separate opentelemetry collector is required for Golang autoinstrumentation with the port set to 4318.
  # exporter:
  #   endpoint: http://opentelemetry-collector.ao-demo:4318
  # propagators:
  #   - tracecontext
  # sampler:
  #   type: always_on
  # resource:
  #   resourceAttributes:
  #     cluster.name: "your-cluster-name"
  # ### End OTel config ###
  java:
    image: ghcr.io/newrelic-experimental/newrelic-agent-operator/instrumentation-java:8.10.0
    env:
    - name: NEW_RELIC_APPLICATION_LOGGING_FORWARDING_ENABLED
      value: "false"
  nodejs:
    image: ghcr.io/newrelic-experimental/newrelic-agent-operator/instrumentation-nodejs:11.15.0
  python:
    image: ghcr.io/newrelic-experimental/newrelic-agent-operator/instrumentation-python:9.8.0
  dotnet:
    image: ghcr.io/newrelic-experimental/newrelic-agent-operator/instrumentation-dotnet:10.23.0
  php:
    image: ghcr.io/newrelic-experimental/newrelic-agent-operator/instrumentation-php:10.19.0.9
  # go:
  #   image: ghcr.io/open-telemetry/opentelemetry-go-instrumentation/autoinstrumentation-go:latest
```


## Create a secret containing your New Relic License Key

Create a secret containing the license key associated with the New Relic account where the APM agents, instrumenting your applications, will send their data.  A secret will be required for each namespace containing applications that will be targeted by the operator.

```
kubectl create secret generic newrelic-key-secret -n ao-demo --from-literal=new_relic_license_key=<NEW RELIC INGEST LICENSE KEY>
```

## Add annotations to existing Deployments

The `newrelic-agent-operator` looks for language-specific annotations when your Pods are scheduled to the cluster.  Think of this as "opting in" to auto-instrumentation.

Below are the currently supported annotations:


```
instrumentation.newrelic.com/inject-java: "true"
instrumentation.newrelic.com/inject-nodejs: "true"
instrumentation.newrelic.com/inject-python: "true"
instrumentation.newrelic.com/inject-dotnet: "true"
instrumentation.newrelic.com/inject-php: "true"
instrumentation.newrelic.com/inject-go: "true"
```

**TODO:** Additional instructions for Golang coming soon...

Add the appropriate annotation to your Deployment.  This tells the operator to inject the newrelic instrumentation init container and auto-instrument the application.

```
apiVersion: apps/v1
kind: Deployment
metadata:
  name: spring-petclinic
spec:
  selector:
    matchLabels:
      app: spring-petclinic
  replicas: 1
  template:
    metadata:
      labels:
        app: spring-petclinic
      annotations:
        instrumentation.newrelic.com/inject-java: "true"
    spec:
      containers:
        - name: spring-petclinic
          image: ghcr.io/pavolloffay/spring-petclinic:latest
          ports:
            - containerPort: 8080
          env:
          - name: NEW_RELIC_APP_NAME
            value: spring-petclinic-demo
```

## Instrumentation Validation

Once your workload is deployed and the Pod(s) are in a `Running` state, you can validate that the `newrelic-instrumentation` container was injected by evaluating the Pod spec.

```
$ kubectl get pod -l app=spring-petclinic -n ao-demo -o jsonpath='{.items[*].spec.initContainers[*].name}'                                                                                                          newrelic-instrumentation%
```

Or, you can simply use a command like `kubectl describe pod` to manually find it in the output.  

## Validating APM Transaction events in New Relic

Run the following NRQL to identify that instrumentation has been successful and your workloads are generating `Transaction` events.

```
FROM Transaction SELECT count(*) facet appName TIMESERIES SINCE 10 minutes ago
```