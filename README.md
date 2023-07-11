# newrelic-agent-operator
<strong>newrelic-agent-operator</strong> is an implementation of a Kubernetes [Operator](https://kubernetes.io/docs/concepts/extend-kubernetes/operator), that manages auto-instrumentation of the workloads using New Relic APM agents.

## Description
Customers expect the APM agent deployment process to be straightforward and hassle-free. They prefer agents that are easy to install and configure. The idea behind <strong>newrelic-agent-operator</strong> is to provide easy and efficient installation to the users. It uses
* Kubernetest operator to create required custom resource for the agent installation. It creates <strong>newrelic-instrumentation</strong> 
  resource.
* It uses [Controllers](https://kubernetes.io/docs/concepts/architecture/controller/),which provide a reconcile function responsible for 
  synchronizing resources until the desired state is reached on the cluster.  
* Init containers that executes required one-time New Relic agent setup in application pods.
* Adding annotation ```instrumentation.newrelic.com/inject-<python>: "true"``` in application deployment yaml does the magic. This 
  annotation helps to trigger language specific installation steps as part of the init container.

## Getting Started
You’ll need a Kubernetes cluster to run against. You can use [KIND](https://sigs.k8s.io/kind) or [Minikube](https://minikube.sigs.k8s.io/docs/start/) to get a local cluster for testing, or run against a remote cluster.
**Note:** Your controller will automatically use the current context in your kubeconfig file (i.e. whatever cluster `kubectl cluster-info` shows).

### Running on the cluster
1. Install Instances of Custom Resources:

```sh
kubectl apply -f config/samples/
```

### Update newrelic-agent-operator
1. Build and push your image to the location specified by `IMG`:

```sh
make docker-build docker-push IMG=<some-registry>/newrelic-agent-operator:tag
```

2. Deploy the controller to the cluster with the image specified by `IMG`:

```sh
make deploy IMG=<some-registry>/newrelic-agent-operator:tag
```

### Uninstall CRDs
To delete the CRDs from the cluster:

```sh
make uninstall
```

### Undeploy controller
Undeploy the controller from the cluster:

```sh
make undeploy
```

### Test It Out
1. Install the CRDs into the cluster:

```sh
make install
```

2. Run your controller (this will run in the foreground, so switch to a new terminal if you want to leave it running):

```sh
make run
```

**NOTE:** You can also run this in one step by running: `make install run`

### Modifying the API definitions
If you are editing the API definitions, generate the manifests such as CRs or CRDs using:

```sh
make manifests
```

**NOTE:** Run `make --help` for more information on all potential `make` targets

More information can be found via the [Kubebuilder Documentation](https://book.kubebuilder.io/introduction.html)

## Things to consider
* Currently, <strong>newrelic-agent-operator</strong> supports instrumentation for [Python](https://docs.newrelic.com/docs/apm/agents/python-agent/configuration/python-agent-configuration/), [Java](https://docs.newrelic.com/docs/apm/agents/java-agent/getting-started/introduction-new-relic-java/), [Node.js](https://docs.newrelic.com/docs/apm/agents/nodejs-agent/getting-started/introduction-new-relic-nodejs/) and [.NET](https://docs.newrelic.com/docs/apm/agents/net-agent/getting-started/introduction-new-relic-net/).
* [cert-manager](https://cert-manager.io/docs/installation/) is required to be presented on cluster before starting newrelic-agent-operator


## Contributing
We encourage your contributions to improve New Relic agent operator! Keep in mind when you submit your pull request, you'll need to sign the CLA via the click-through using CLA-Assistant. You only have to sign the CLA one time per project. If you have any questions, or to execute our corporate CLA, required if your contribution is on behalf of a company, please drop us an email at opensource@newrelic.com.

A note about vulnerabilities

As noted in our security policy, New Relic is committed to the privacy and security of our customers and their data. We believe that providing coordinated disclosure by security researchers and engaging with the security community are important means to achieve our security goals.

If you believe you have found a security vulnerability in this project or any of New Relic's products or websites, we welcome and greatly appreciate you reporting it to New Relic through HackerOne.


## License

Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
