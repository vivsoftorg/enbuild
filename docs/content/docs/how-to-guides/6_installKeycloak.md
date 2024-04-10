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

# Configure the Keycloak

## Create a realm

keyclaok realm is a logical container for a set of related users, applications, and groups. 

Access the `http://localhost:9090/auth/admin/#/realms` and click on `Add realm`.

To save we have provided a realm.json file. All you have to do is import it.

<picture><img src="/images/how-to-guides/keycloak/add-realm.png" alt="Add Realm"></img></picture>

Click on the `Browse` button and provide the [realm.json](https://github.com/vivsoftorg/enbuild/blob/main/examples/keycloak/realm.json) file.
<picture><img src="/images/how-to-guides/keycloak/add-realm2.png" alt="Add Realm2"></img></picture>

## Create a enbuild-ui Client

Clients are entities that can request Keycloak to authenticate a user. Most often, clients are applications and services that want to use Keycloak to secure themselves and provide a single sign-on solution. 

Clients can also be entities that just want to request identity information or an access token so that they can securely invoke other services on the network that are secured by Keycloak.

Click on the Clients menu from the left pane. All the available clients for the selected Realm will get listed here.

Rather creating Clients manually we will use `Import` so that the you do not have provide the details manually.

To create a new client, click on `Import client`. 

<picture><img src="/images/how-to-guides/keycloak/import_client.png" alt="import_client"></img></picture>

Click on the `Browse` button and provide the [enbuild-ui.json](https://github.com/vivsoftorg/enbuild/blob/main/examples/keycloak/kc-client-enbuild-ui.json)file.

<picture><img src="/images/how-to-guides/keycloak/import_client_1.png" alt="import_client"></img></picture>

and Save the client.

After import make sure you configure the valid redirect URI's. 

You have to provide the URI(FQDN) of ENBUILD that you have configured while installing the ENBUILD Helm chart.

<picture><img src="/images/how-to-guides/keycloak/configure-client.png" alt="configure-client"></img></picture>

## Create a enbuild Client

Use the same method to  `Import` the [enbuild.json](https://github.com/vivsoftorg/enbuild/blob/main/examples/keycloak/kc-client-enbuild.json) client.


After you save the `enbuild` client. 
Go to the Credentials tab of the `enbuild` client to get the Client Secret which is required for configuring the (Keycloak in ENBUILD admin setting)[../configuring-enbuild/#configure-keycloak]

<picture><img src="/images/how-to-guides/keycloak/credentials.png" alt="credentials"></img></picture>


## Realm Roles and Groups 

:exclamation: **Note:** All these roles and groups are already created in the realm, when we imported the realms.

ENBUILD support different roles , and based on these roles, differnt type of catalog can be visible. 

By default 
- ENBUILD_ADMIN
- ENBUILD-APPDEV
- ENBUILD-DATAOPS
- ENBUILD-DEVOPS
- ENBUILD_USER 

these Groups are present in the realm we imported. 

<picture><img src="/images/how-to-guides/keycloak/keycloak-groups.png" alt="keycloak-groups"></img></picture>

These Groups mapped to the following roles in one to one relationship.
- admin -- ENBUILD_ADMIN
- appdev -- ENBUILD-APPDEV
- dataops -- ENBUILD-DATAOPS
- devops -- ENBUILD-DEVOPS 
- user -- ENBUILD_USER 


<picture><img src="/images/how-to-guides/keycloak/keycloak-roles.png" alt="keycloak-roles"></img></picture>

## Create Users

Users are entities that are able to log into your system. 
They can have attributes associated with themselves like email, username, address, phone number, and birth day. 
<picture><img src="/images/how-to-guides/keycloak/crete_user.png" alt="crete_user.png"></img></picture>
They can be assigned group membership  by clicking on the `Join Groups` button or `Groups` tab from the profile. 

<picture><img src="/images/how-to-guides/keycloak/add_user_to_group.png" alt="crete_user.png"></img></picture>

Click on the `Add Member` button and select the group you want to add the user to, then click on the `Add` button.

<picture><img src="/images/how-to-guides/keycloak/add_user_to_group1.png" alt="add_user_to_group"></img></picture>


# Exposing the KeyCloak service in Kubernetes.
So far we have accessed the keycloak service from the browser by using the port-forwarding and then accessing it via local-port.
But for production usage and to [configure it as SSO inside ENBUILD](../configuring-enbuild/#configure-keycloak) we need to expose it on a public IP address.

There are multiple ways to achive that -

- Expose the `keycloak-http` service as LoadBalancer - This way the kubernertes will create a external loadbalancer and expose the service on a public IP address. 
- Create an Ingress resource for the `keycloak-http` service - This way we can use a ingress controller to route traffic to the keycloak service based on the hostname.
- If you have istio deployed and configured in your cluster, then expose it as virtual service. [This is the sample virtaul service configurations you can use](https://github.com/vivsoftorg/enbuild/blob/main/examples/keycloak/istio-virtaul-service-for-keycloak.yaml)

:exclamation: **Note:** Make sure to adjust the `hosts` to the right value to which you want to expose the service.

# Create a DNS record for the keycloak 
after you have exposed the `keycloak-http` service using any of the above methods, you can create a DNS record pointing to the public IP address or FQDN of the loadbalncer. 

e.g. 
if you have exposed the service using the istio virtual service you will need to create a DNS record for this service in your DNS provider.

The IP address of the DNS will be the EXTERNAL-IP of the istio-ingressgateway service, you can find that using the command below.

```bash
kubectl get svc -n istio-system istio-ingressgateway -o jsonpath='{.status.loadBalancer.ingress[*].ip}'
```
Create a DNS type A record using the IP address from the command above.

Once done you can proceed to [configure keycloak as SSO for ENBUILD](../configuring-enbuild/#configure-keycloak)
