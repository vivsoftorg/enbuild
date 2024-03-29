---
title: "Installing  Istio"
description: "Steps to Install Istio"
summary: "Steps to Install Istio"
draft: false
menu:
  docs:
    parent: ""
    identifier: "installIstio"
weight: 802
toc: true
seo:
  title: "" # custom title (optional)
  description: "" # custom description (recommended)
  canonical: "" # custom canonical URL (optional)
  noindex: false # false (default) or true
---

# Prerequisites
Before you begin, ensure that you have the following prerequisites in place:
- Access to Platform1 registry
- Kubernetes cluster is up and running and you have access to it via `kubectl` command.
- Helm 3 installed on your system.


## Login to Platform1 registry
Login to Platform1 registry by using the following command:
  

```shell
helm registry login registry1.dso.mil/bigbang
WARNING: Kubernetes configuration file is group-readable. This is insecure. Location: /root/.kube/config
WARNING: Kubernetes configuration file is world-readable. This is insecure. Location: /root/.kube/config
Username: Juned_Memon
Password:
Login Succeeded
```

## Create Namespace and Image pull secret

We need to create the `istio-system` and `istio-operator` namespaces.
```
kubectl create ns istio-operator
kubectl create ns istio-system
```
Next, we need to create the imagePullSecret for pulling the images from Platform1 registry.

First export the REGISTRY1_USER and REGISTRY1_PASSWORD with your P1 credentials.

```
export REGISTRY1_USER=<YOUR_REGISTRY1_USER>
export REGISTRY1_PASSWORD=<YOUR_REGISTRY1_PASSWORD>
```

Next, create the imagePullSecret for pulling the images from Platform1 registry.


```
kubectl create secret -n istio-operator docker-registry private-registry --docker-server=registry1.dso.mil --docker-username=$REGISTRY1_USER --docker-password=$REGISTRY1_PASSWORD

kubectl create secret -n istio-system docker-registry private-registry --docker-server=registry1.dso.mil --docker-username=$REGISTRY1_USER --docker-password=$REGISTRY1_PASSWORD
```

## Install istio-operator Helm charts

Now install `istio-operator` Helm charts using the P1 Helm chart:

```
helm upgrade --install --namespace istio-operator istio-operator oci://registry1.dso.mil/bigbang/istio-operator --version 1.20.4-bb.0 --set imagePullSecrets[0]="private-registry"  --set createNamespace=false 
```

Verify the `istio-operator` pod is up and running:


```
# kubectl get po -n istio-operator
NAME                              READY   STATUS    RESTARTS   AGE
istio-operator-7b5fff8cfb-h6w4k   1/1     Running   0          18s
```

### Optional - Install CertManager 

Before installing the stio-controlplane we need a `wildcard-cert` secret containing the SSL certificate for the domain on which are planning to expose the virtual services.

You can use `CertManager` manage that TLS certificates and keys, or you can create them manually using openssl.

Here, we will install CertManager and use self-signed-certificat. 

To deploy a proper certificate using `CertManager` refer the official [Documentation](https://cert-manager.io/)

```
kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.14.4/cert-manager.yaml


# kubectl get pods --namespace cert-manager
NAME                                       READY   STATUS    RESTARTS   AGE
cert-manager-67c98b89c8-g428w              1/1     Running   0          5m12s
cert-manager-cainjector-5c5695d979-7qczq   1/1     Running   0          5m12s
cert-manager-webhook-7f9f8648b9-2bt85      1/1     Running   0          5m12s
```

Create a self signed cluster issuer

```
---
apiVersion: cert-manager.io/v1
kind: ClusterIssuer
metadata:
  name: selfsigned-ca-issuer
spec:
  selfSigned: {}
---
apiVersion: cert-manager.io/v1
kind: Certificate
metadata:
  name: selfsigned-ca
  namespace: cert-manager
spec:
  isCA: true
  commonName: selfsigned-ca
  secretName: root-secret
  privateKey:
    algorithm: ECDSA
    size: 256
  issuerRef:
    name: selfsigned-ca-issuer
    kind: ClusterIssuer
    group: cert-manager.io
---
apiVersion: cert-manager.io/v1
kind: ClusterIssuer
metadata:
  name: selfsigned
spec:
  ca:
    secretName: root-secret
---
apiVersion: cert-manager.io/v1
kind: Certificate
metadata:
  name: self-sign-cert
  namespace: istio-system
spec:
  secretName: wildcard-cert
  commonName: ijuned.com
  dnsNames:
  - ijuned.com
  - "*.ijuned.com"
  issuerRef:
    name: selfsigned
    kind: ClusterIssuer
---
```
Apply the configurations
```
kubectl apply -f self-sign-cert.yaml
```

## Install istio-controlplane

Before installing the stio-controlplane we need a `wildcard-cert` secret containing the SSL certificate for the domain on which are planning to expose the virtual services. 

You can create them manaully if you have the `tls.key` and `tls.cert` of your private certificate.

Otherwise you can use [Certmanager](#optional---install-certmanager)

To install the istio controlplane Helm chart.

The domain input that you provide, will be used to create a `host` entry in the istio `Gateway` named `main`

```
helm upgrade --install --namespace istio-system istio oci://registry1.dso.mil/bigbang/istio --version 1.20.4-bb.0 --set imagePullSecrets[0]="private-registry" --set domain="ijuned.com"
```

