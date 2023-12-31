 StackInsights Cloud on Kubernetes
============

A bridge project between [StackInsights](https://github.com/apache/stackinsights) and Kubernetes.

SWCK is a platform for the StackInsights user, provisions, upgrades, maintains StackInsights relevant components, and makes them work natively on Kubernetes.

# Features

1. Java Agent Injector: Inject the java agent into the application pod natively.
1. Operator: Provision and maintain StackInsights backend components.
1. Custom Metrics Adapter: Provides custom metrics come from StackInsights OAP cluster for autoscaling by Kubernetes HPA

# Build images

Issue below instrument to get the docker image:

```
make
```

or 

```
make build
```

To onboard operator or adapter, you should push the image to a registry where the kubernetes cluster can pull it.

## Onboard Java Agent Injector and Operator

The java agent injector and operator share a same binary. To onboard them, you should follow:

* To install the java agent injector and operator in an existing cluster, make sure you have  [`cert-manager`](https://cert-manager.io/docs/installation/) installed.
* Apply the manifests for the Controller and CRDs in `config`:

 ```
 kubectl apply -f config/operator-bundle.yaml
 ```

## Onboard Custom Metrics Adapter

* Deploy OAP server by referring to Operator Quick Start.
* Apply the manifests for an adapter in `config`:

 ```
 kubectl apply -f config/adapter-bundle.yaml
 ```
