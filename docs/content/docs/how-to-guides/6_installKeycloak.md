---
title: "Install and Configure Keycloak for ENBUILD SSO"
description: "Install and Configure Keycloak for ENBUILD SSO"
summary: "Install and Configure Keycloak for ENBUILD SSO"
draft: false
menu:
  docs:
    parent: "docs/how-to-guides/"
    identifier: "installKeycloak"
weight: 206
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

We need to create the `keycloak` namespace
```
kubectl create ns keycloak
```
Next, we need to create the imagePullSecret for pulling the images from Platform1 registry.

First export the REGISTRY1_USER and REGISTRY1_PASSWORD with your P1 credentials.

```
export REGISTRY1_USER=<YOUR_REGISTRY1_USER>
export REGISTRY1_PASSWORD=<YOUR_REGISTRY1_PASSWORD>
```

Next, create the imagePullSecret for pulling the images from Platform1 registry.


```
kubectl create secret -n keycloak docker-registry private-registry --docker-server=registry1.dso.mil --docker-username=$REGISTRY1_USER --docker-password=$REGISTRY1_PASSWORD
```

## Install keycloak Helm charts

Now install `keycloak` Helm charts using the P1 Helm chart:

Create a keycloak-input-values.yaml

```
## Overrides the default args for the Keycloak container
args:
  - "start"
  - "--http-port=8080"
  - "--import-realm"

# Additional environment variables for Keycloak
# https://www.keycloak.org/server/all-config
extraEnv: |-
  - name: KC_HTTP_ENABLED
    value: "true"
  - name: KC_PROXY
    value: edge
  - name: KC_HOSTNAME_STRICT
    value: "false"
  - name: KC_HOSTNAME_STRICT_HTTPS
    value: "false"
  - name: KC_HTTP_RELATIVE_PATH
    value: /auth
```

```
helm upgrade --install --namespace keycloak keycloak oci://registry1.dso.mil/bigbang/keycloak --version 23.0.7-bb.2  -f keycloak-input-values.yaml
```

Verify the `keycloak` pod is up and running:


```
# kubectl get po -n keycloak
NAME                    READY   STATUS    RESTARTS   AGE
keycloak-0              1/1     Running   0          105s
keycloak-postgresql-0   1/1     Running   0          16m
```

# Access keycloak UI

Access the Keycloak web interface using Port forwarding

```
❯ kubectl port-forward svc/keycloak-http 9090:80 -n keycloak
Forwarding from 127.0.0.1:9090 -> 8080
Forwarding from [::1]:9090 -> 8080
Handling connection for 9090
Handling connection for 9090
Handling connection for 9090
```

# Get the Keycloak admin credentials

Use the following command to get the admin credentials for the `keycloak` instance:

```
USER=$(kubectl get secret keycloak-env -n keycloak -o jsonpath='{.data.KEYCLOAK_ADMIN_PASSWORD}' | base64 --decode)
PASSWORD=$(kubectl get secret keycloak-env -n keycloak -o jsonpath='{.data.KEYCLOAK_ADMIN}' | base64 --decode)
echo "Keycloak Admin user is $USER and password is $PASSWORD"
```

