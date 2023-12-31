# Deploy OAP server and UI with default settings

In this example, we will deploy an OAP server and UI to Kubernetes cluster with 
default settings specified by their Custom Resource Defines(CRD).

## Install Operator

Follow [Operator installation instrument](../../README.md#operator) to install the operator.

## Deploy OAP server and UI with default setting

Clone this repo, then change current directory to [samples](../../operator/config/samples).

Issue the below command to deploy an OAP server and UI.

```shell
kubectl apply -f default.yaml
```

Get created custom resources as below:

```shell
$ kubectl get oapserver,ui
NAME                                               INSTANCES   RUNNING   ADDRESS
oapserver.operator.stackinsights.apache.org/default   1           1         default-oap.stackinsights-swck-system

NAME                                        INSTANCES   RUNNING   INTERNALADDRESS                     EXTERNALIPS   PORTS
ui.operator.stackinsights.apache.org/default   1           1         default-ui.stackinsights-swck-system                 [80]
```

## View the UI
In order to view the UI from your browser, you should get the external address from the ingress generated by the UI custom resource firstly.

```shell
$ kubectl get ingresses
NAME         HOSTS                ADDRESS          PORTS   AGE
default-ui   demo.ui.stackinsights    <External_IP>   80      33h
```

Edit your local `/etc/hosts` to append the following host-ip mapping.

```
demo.ui.stackinsights <External_IP>
```

Finally, navigate your browser to `demo.ui.stackinsights` to access UI service.

Notice, please install an ingress controller to your Kubernetes environment.
