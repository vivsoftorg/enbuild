---
title: "Installing  Nginx Ingress"
description: "Steps to Install Nginx Ingress"
summary: "Steps to Install Nginx Ingress"
draft: false
menu:
  docs:
    parent: "docs/how-to-guides/"
    identifier: "installNginx Ingress"
weight: 209
toc: true
seo:
  title: "" # custom title (optional)
  description: "" # custom description (recommended)
  canonical: "" # custom canonical URL (optional)
  noindex: false # false (default) or true
---

# Prerequisites
Before you begin, ensure that you have the following prerequisites in place:
- Kubernetes cluster is up and running and you have access to it via `kubectl` command.
- Helm 3 installed on your system.



## Create Namespace and Image pull secret

To install the Nginx Ingress Controller to your cluster, you’ll first need to add its repository to Helm by running:

```
helm repo add ingress-nginx https://kubernetes.github.io/ingress-nginx
```
The output will be:

```
"ingress-nginx" has been added to your repositories
```

Update your system to let Helm know what it contains:


```
helm repo update ingress-nginx
```

## install the Nginx ingress:

run the following command to install the Nginx ingress:

```
helm install nginx-ingress ingress-nginx/ingress-nginx --set controller.publishService.enabled=true
```

This command installs the Nginx Ingress Controller from the stable charts repository, names the Helm release nginx-ingress, and sets the publishService parameter to true.

Once it has run, you will receive an output similar to this:


```
NAME: nginx-ingress
LAST DEPLOYED: Wed Apr 10 16:19:24 2024
NAMESPACE: default
STATUS: deployed
REVISION: 1
TEST SUITE: None
NOTES:
The ingress-nginx controller has been installed.
It may take a few minutes for the load balancer IP to be available.
You can watch the status by running 'kubectl get service --namespace default nginx-ingress-ingress-nginx-controller --output wide --watch'

An example Ingress that makes use of the controller:
  apiVersion: networking.k8s.io/v1
  kind: Ingress
  metadata:
    name: example
    namespace: foo
  spec:
    ingressClassName: nginx
    rules:
      - host: www.example.com
        http:
          paths:
            - pathType: Prefix
              backend:
                service:
                  name: exampleService
                  port:
                    number: 80
              path: /
    # This section is only required if TLS is to be enabled for the Ingress
    tls:
      - hosts:
        - www.example.com
        secretName: example-tls

If TLS is enabled for the Ingress, a Secret containing the certificate and key must also be provided:

  apiVersion: v1
  kind: Secret
  metadata:
    name: example-tls
    namespace: foo
  data:
    tls.crt: <base64 encoded cert>
    tls.key: <base64 encoded key>
  type: kubernetes.io/tls
...
```

Helm has logged what resources it created in Kubernetes as a part of the chart installation.


### Optional - Install CertManager 

To use TLS with ingress we need to create a certificate with propr TLS data. You can create that manually. But for reference we will use the Cert-manager to manage our certificates. 


Here, we will install CertManager and use self-signed-certificate.

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
  namespace: kube-system
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

Follow the the official [Documentation](https://cert-manager.io/) for using cert-manager issuers other than self-sign.

## Install ENBUILD with exposing service using Ingress

To install ENBUILD with exposing frontend serivice using Ingress you can use these [example helm input file](https://github.com/vivsoftorg/enbuild/blob/main/examples/enbuild/with_ingress.yaml).


:exclamation: **Note:** Make sure to change the `domain` and the certificate name as per your requirments.

## Create DNS records 

Run this command to watch the ingress-ingress-nginx-controller Load Balancer become available:
```
kubectl --namespace default get services -o wide -w nginx-ingress-ingress-nginx-controller
```
This command fetches the Nginx Ingress service in the default namespace and outputs its information, but the command does not exit immediately. With the -w argument, it watches and refreshes the output when changes occur.

While waiting for the Load Balancer to become available, you may receive a pending response:

After some time has passed, the IP address of your newly created Load Balancer will appear:
```
NAME                                     TYPE           CLUSTER-IP      EXTERNAL-IP     PORT(S)                      AGE     SELECTOR
nginx-ingress-ingress-nginx-controller   LoadBalancer   10.43.254.211   192.168.0.108   80:31730/TCP,443:32755/TCP   5m56s   app.kubernetes.io/component=controller,app.kubernetes.io/instance=nginx-ingress,app.kubernetes.io/name=ingress-nginx
```


Next, you’ll need to ensure that your domains like `enbuild.ijuned.com` and `rabbitmq.ijuned.com` are pointed to the Load Balancer via A records. This is done through your DNS provider. To configure your DNS records follow your DNS provider documentaion.

You’ve installed the Nginx Ingress maintained by the Kubernetes community. It will route HTTP and HTTPS traffic from the Load Balancer to appropriate back-end Services configured in the Ingress Resources.

Once the DNS / Host entry is added you can access the ENBUILD using the created ingress domain 

```
❯ kubectl get ing -n enbuild
NAME                  CLASS   HOSTS                ADDRESS         PORTS     AGE
enbuild-enbuild-ing   nginx   enbuild.ijuned.com   192.168.0.108   80, 443   10m
```