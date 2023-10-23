 StackInsights Cloud on Kubernetes
============

A bridge project between [StackInsights](https://github.com/apache/stackinsights) and Kubernetes.

SWCK is a platform for the StackInsights user that provisions, upgrades, maintains StackInsights relevant components, and makes them work natively on Kubernetes.

# Features

* Java Agent Injector: Inject the java agent into the application pod natively.
  * Inject the java agent into the application pod.
  * Leverage a global configuration to simplify the agent and injector setup.
  * Use the annotation to customize specific workloads.
  * Synchronize injecting status to `JavaAgent` CR for monitoring purposes.
* Operator: Provision and maintain StackInsights backend components.
* Custom Metrics Adapter: Provides custom metrics coming from StackInsights OAP cluster for autoscaling by Kubernetes HPA

# Quick Start

There are two ways to install swck.
* Go to the [download page](https://demo.stackinsights.ai/downloads/#StackInsightsCloudonKubernetes) to download the latest release binary, `stackinsights-swck-<SWCK_VERSION>-bin.tgz`. Unarchive the package to a folder named `stackinsights-swck-<SWCK_VERSION>-bin`
* Apply the kustomization directory from github.

## Java Agent Injector

* Install the [Operator](#operator)
* Label the namespace with `swck-injection=enabled`

```shell
$ kubectl label namespace default(your namespace) swck-injection=enabled
```

* Add label `swck-java-agent-injected: "true"` to the workloads

For more details, please read [Java agent injector](/docs/java-agent-injector.md)

## Operator

* To install the operator in an existing cluster, ensure you have [`cert-manager`](https://cert-manager.io/docs/installation/) installed.
* Apply the manifests for the Controller and CRDs in release/config:

 ```
 kubectl apply -f stackinsights-swck-<SWCK_VERSION>-bin/config/operator-bundle.yaml
 ```

* Also, you could deploy the operator quickly based on **Master Branch** or **Stable Release**:
 
 ```
 kubectl apply -k "github.com/humancloud/si-swck/operator/config/default"
 ```

or

 ```
 kubectl apply -k "github.com/humancloud/si-swck-swck/operator/config/default?ref=v0.8.0"
 ```

For more details, please refer to [deploy operator](docs/operator.md)

## Custom Metrics Adapter
  
* Deploy the OAP server by referring to Operator Quick Start.
* Apply the manifests for an adapter in release/adapter/config:

 ```
 kubectl apply -f stackinsights-swck-<SWCK_VERSION>-bin/config/adapter-bundle.yaml
 ```
* Also, you could deploy the adapter quickly based on **Master Branch** or **Stable Release**:
 
 ```
 kubectl apply -k "github.com/humancloud/si-swck/adapter/config"
 ```

or

 ```
 kubectl apply -k "github.com/humancloud/si-swck/adapter/config?ref=v0.8.0"
 ```

For more details, please read [Custom metrics adapter](docs/custom-metrics-adapter.md)
