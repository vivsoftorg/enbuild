---
title: "Deploying ENBUILD exposing the service using Istio"
description: "Steps to deploy ENBUILD on top of Istio"
summary: ""
date: 2024-03-28T16:04:48+02:00
lastmod: 2024-03-28T16:04:48+02:00
draft: false
menu:
  docs:
    parent: ""
    identifier: "deployEnbuildWithIstio"
weight: 800
toc: true
seo:
  title: "" # custom title (optional)
  description: "" # custom description (recommended)
  canonical: "" # custom canonical URL (optional)
  noindex: false # false (default) or true
---

Follow these step-by-step instructions to deploy ENBUILD locally for testing.

## Prerequisites

### Existing Kubernetes Cluster

Ensure that you have access to a Kubernetes cluster and obtain the KubeConfig file.

You can use [rancher-desktop](https://docs.rancherdesktop.io/getting-started/installation/) or [k3d](https://k3d.io/v5.6.0/) to spin up a local Kubernetes cluster.

### ENBUILD Container Images

Access to the ENBUILD container images are required for this deployment.
These images are published to the VivSoft managed container reigistry on `registry.gitLab.com`.
Make sure that you have the necessary credentials to pull these images.

### Helm CLI

The [Helm](https://helm.sh/) streamlines and automates Kubernetes deployments by managing charts, enabling users to easily package, version, and deploy complex applications.

### Istio is deployed on your cluster

Istio is a service mesh that enhances connectivity, security, traffic management, and observability for microservices.

You can follow the [docs](../how-to-guides/installing-istio/) to install Istio on your cluster using the P1 Chart and Images


## Deployment Steps:

Following are the steps you will need to take to deploy ENBUILD to your Kubernetes cluster.

### Add ENBUILD Helm Chart Repository

To add the ENBUILD Helm chart repository, run the following command:

```bash
helm repo add vivsoft https://vivsoftorg.github.io/enbuild

"vivsoft" has been added to your repositories
```

### Configure ENBUILD Helm Values

Before deploying ENBUILD to the Kubernetes cluster, you will need to create a custom values.yaml file so that we can specify configurations unique to this deployment.

For local deployment however we require minimum deployment values.

:exclamation: **Note:** For more information about the complete set of ENBUILD Helm values click [here](/docs/getting-started/helm-values/)!

```yaml
imageCredentials:
  registry: registry.gitlab.com
  username: registry_user_name
  password: registry_password
global:
  domain: ijuned.com         # Set the proper doamin.
  istio:
    enabled: true
    gateway: istio-system/main # set to the proper istio gateway. This istio gateway must have above domain added as `hosts` 
```

:zap: **Note:** The `imageCredentials` section is only required until the images are available publically.

### Deploy ENBUILD HELM Chart

Make sure you update the values input to reference the values you created in Step 2.
Execute the command below.

```bash
helm upgrade --install --namespace enbuild enbuild vivsoft/enbuild --create-namespace -f target/quick_install.yaml

Release "enbuild" does not exist. Installing it now.
NAME: enbuild
LAST DEPLOYED: Fri Mar 22 17:37:23 2024
NAMESPACE: enbuild
STATUS: deployed
REVISION: 1
TEST SUITE: None
NOTES:
1. Get the application URL by running these commands:
  echo "Visit http://127.0.0.1:3000 to use your application after starting the port forward"
  kubectl --namespace enbuild port-forward svc/enbuild-enbuild-ui 3000:80
```

### Validate ENBUILD Deployment

Use the following commands to validate the ENBUILD pods are up and running.

```bash
kubectl get pods -n enbuild

NAME                                       READY   STATUS    RESTARTS         AGE
enbuild-enbuild-genai-8488c86d6f-csfmn     1/1     Running   0                76m
enbuild-enbuild-ui-56f5667d5b-4xckt        1/1     Running   0                76m
enbuild-mongodb-0                          1/1     Running   0                76m
enbuild-rabbitmq-0                         1/1     Running   0                76m
enbuild-enbuild-backend-66676f8cd8-hxtbr   1/1     Running   0                76m
enbuild-enbuild-user-b87d95b45-c79p6       1/1     Running   0                76m
enbuild-enbuild-request-7c47c6d67b-j2fnd   1/1     Running   1 (73m ago)      76m
enbuild-enbuild-ml-6f944ff759-ztdj6        1/1     Running   1 (73m ago)      76m
enbuild-rabbitmq-1                         1/1     Running   0                73m
enbuild-rabbitmq-2                         1/1     Running   0                72m
enbuild-enbuild-mq-575c965764-zcnlg        1/1     Running   18 (6m24s ago)   76m

```

:exclamation: **Note:** You might see restarts of the enbuild-enbuild-mq-\* pod until the RabbitMQ service is up and running.

**Validate the ENBUILD services are setup correctly**

```bash
kubectl get services -n enbuild

NAME                        TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)                                 AGE
enbuild-rabbitmq-headless   ClusterIP   None            <none>        4369/TCP,5672/TCP,25672/TCP,15672/TCP   80s
enbuild-mongo               ClusterIP   10.43.230.6     <none>        27017/TCP                               80s
enbuild-enbuild-user        ClusterIP   10.43.140.228   <none>        80/TCP                                  80s
enbuild-enbuild-ui          ClusterIP   10.43.110.47    <none>        80/TCP                                  80s
enbuild-enbuild-backend     ClusterIP   10.43.146.20    <none>        80/TCP                                  80s
enbuild-rabbitmq            ClusterIP   10.43.54.197    <none>        5672/TCP,4369/TCP,25672/TCP,15672/TCP   80s
```

### Access ENBUILD

Once the ENBUILD is installed, there will be two virtul services created. 
You can create a DNS record/host entry on your local machine pointing to these `virtul services` to the istio gateway loadbalancer/ip that you have used to install istio.
`global.istio.gateway`  

```bash
kubectl get vs -A

NAMESPACE   NAME                 GATEWAYS                          HOSTS                     AGE
enbuild     enbuild-enbuild-ui   ["istio-system/public-gateway"]   ["enbuild.ijuned.com"]    4s
enbuild     enbuild-rabbitmq     ["istio-system/public-gateway"]   ["rabbitmq.ijuned.com"]   4s
```

To find out the loadbalancer/ip of the istio gateway use the following command
```bash
kubectl get svc -n istio-system istio-ingressgateway
NAMESPACE        NAME                                 TYPE           CLUSTER-IP      EXTERNAL-IP   PORT(S)                                      AGE
istio-system     istio-ingressgateway                 LoadBalancer   10.43.60.241    <pending>     15021:32686/TCP,80:31687/TCP,443:30260/TCP   4m29s

```

After the DNS entry / host entry , you can access the ENBUILD using the virtual service e.g. 

https://enbuild.ijuned.com

:exclamation: **Note:** If you have used self signed certificate, the browser will complain for it.




### Uninstall ENBUILD

Use the following command to uninstall ENBUILD from your Kubernetes cluster.

```bash
helm uninstall enbuild -n enbuild

release "enbuild" uninstalled
```
