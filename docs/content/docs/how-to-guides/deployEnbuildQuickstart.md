---
title: "Deploying ENBUILD (Quickstart from GitLab)"
description: "Steps to deploy ENBUILD"
summary: ""
date: 2023-09-07T16:04:48+02:00
lastmod: 2023-09-07T16:04:48+02:00
draft: false
menu:
  docs:
    parent: ""
    identifier: "deployEnbuildQuickstart"
weight: 800
toc: true
seo:
  title: "" # custom title (optional)
  description: "" # custom description (recommended)
  canonical: "" # custom canonical URL (optional)
  noindex: false # false (default) or true
---

Follow these step-by-step instructions to deploy ENBUILD on a Kubernetes cluster from GitLab.

**Prerequisites:**

- **Existing Kubernetes Cluster** – Ensure that you have access to a Kubernetes cluster and obtain the KubeConfig file. You will need the necessary permissions to deploy resources to the cluster.
- **ENBUILD Helm Chart** – You will need access to the ENBUILD Helm Chart. You can either download it locally or use the VivSoft hosted version on GitLab.com.
- **ENBUILD Container Images** – Access to the ENBUILD container images are required for this deployment. These images will need to be staged in an existing container registry or accessed from the VivSoft managed container reigistry on Registry.GitLab.com. Make sure that you have the necessary credentials to pull these images.
- **Helm CLI** - The Helm CLI streamlines and automates Kubernetes deployments by managing charts, enabling users to easily package, version, and deploy complex applications.

:exclamation: **Note:** For more information about any of the prerequisites above click [here](../../reference)!

**Deployment Steps:**

1. **Add ENBUILD Helm Chart**\
   The ENBUILD helm chart is available and hosted on GitLab.com. You will need to add the HELM repo locally to your workstation. **You will need to provide your GitLab user name and GitLab API Access Token.**

   ```bash
   export GITLAB_USERNAME=MY_GITLAB_USERNAME
   export GITLAB_PASSWORD=MY_GITLAB_TOKEN
   helm repo add --username $GITLAB_USERNAME --password $GITLAB_PASSWORD vivsoft https://gitlab.com/api/v4/projects/30816323/packages/helm/stable
   ```

   ```bash
   "vivsoft" has been added to your repositories
   ```

2. **Configure ENBUILD Helm Values**\
   Before deploying ENBUILD to the Kubernetes cluster, you will need to create a custom values.yaml file so that we can specify configurations unique to this deployment. The minimum deployment values file can be found below. **In Step 3, this YAML file will be referenced as "enbuild-demo-values.yaml".**

   ```yaml
   global:
     domain: enbuild.demo
     disable_tls_gitlab: false # Set to true if you are using self-signed certificates
     ingress:
       enabled: true
     istio:
       enabled: false
       gateway: istio-system/public-gateway
     auth_plugin: local
   imageCredentials:
     registry: registry.gitlab.com
     username: MY_GITLAB_USERNAME
     password: MY_GITLAB_TOKEN
   rabbitmq:
     enabled: true
     replicaCount: 3
     auth:
       username: admin
       password: password
       erlangCookie: lamba
     clustering:
       forceBoot: true
     env: demo
     queue_prefix: enbuild-demo
   mongodb:
     enabled: true
     mongo_root_username: "enbuild-demo"
     mongo_root_password: "password"
     mongo_server: "If you are using cosmosDB then set the right cosmosDB endpoint as mongo_server and set the enabled=false"
     image:
       repository: registry.gitlab.com/enbuild-staging/container_registry/opensource-mongodb/mongodb
       tag: 4.4.5
       pullPolicy: Always
   enbuildUi:
     image:
       repository: registry.gitlab.com/enbuild-staging/vivsoft-platform-ui/enbuild-frontend
       pullPolicy: Always
     replicas: 1
     service_type: ClusterIP
     node_port: 30080
     hostname: enbuild # hostname.domain become FQDN
   enbuildBk:
     image:
       repository: registry.gitlab.com/enbuild-staging/vivsoft-platform-ui/enbuild-backend
       pullPolicy: Always
     replicas: 1
     service_type: ClusterIP
     encryption_key: "demokey123"
   enbuildMl:
     image:
       repository: registry.gitlab.com/enbuild-staging/vivsoft-platform-ui/enbuild-ml
       pullPolicy: Always
     replicas: 1
     service_type: ClusterIP
   enbuildGenAI:
     api_key: "dummy"
     image:
       repository: registry.gitlab.com/enbuild-staging/vivsoft-platform-ui/enbuild-genai
       pullPolicy: Always
     replicas: 1
     service_type: ClusterIP
   enbuildRequest:
     image:
       repository: registry.gitlab.com/enbuild-staging/vivsoft-platform-ui/enbuild-request
       pullPolicy: Always
     replicas: 1
     service_type: ClusterIP
   enbuildUser:
     image:
       repository: registry.gitlab.com/enbuild-staging/vivsoft-platform-ui/enbuild-user
       pullPolicy: Always
     replicas: 1
     service_type: ClusterIP
   enbuildSync:
     image:
       repository: registry.gitlab.com/enbuild-staging/vivsoft-platform-ui/enbuild-cronjob
       pullPolicy: Always
     replicas: 1
   ```

   :zap: **Click [here](https://gitlab.com/enbuild-staging/enbuild-helm-chart/-/blob/main/values.yaml) to see the latest example values.yaml available in the ENBUILD repository.**

3. **Deploy ENBUILD**\
    Make sure you update the values input to reference the values you created in Step 2.
   Execute the command below.

   ```bash
   helm upgrade --install --namespace enbuild enbuild vivsoft/enbuild --version 1.0.5 --create-namespace -f enbuild-demo-values.yaml
   ```

   ```bash
   Release "enbuild" does not exist. Installing it now.
   NAME: enbuild
   LAST DEPLOYED: Tue Jan 16 15:28:49 2024
   NAMESPACE: enbuild
   STATUS: deployed
   REVISION: 1
   TEST SUITE: None
   NOTES:
     Get the application URL by running these commands:
     echo "Visit http://127.0.0.1:8080 to use your application after starting the port forward"
     kubectl --namespace enbuild port-forward svc/enbuild-enbuild-ui 8080:80
   ```

4. **Validate ENBUILD Deployment**\
    Use the following commands to validate the ENBUILD deployment.\
    \
    **Validate the ENBUILD pods are up and running.**

   ```bash
   kubectl get pods -n enbuild
   ```

   ```bash
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

   :exclamation: **Note:** You might see restarts of the enbuild-enbuild-mq-\* pod until the VCS is configured.
   \
    \
    **Validate the ENBUILD services are up and running.**

   ```bash
     kubectl get services -n enbuild
   ```

   ```bash
   NAME                        TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)                                 AGE
   enbuild-rabbitmq-headless   ClusterIP   None            <none>        4369/TCP,5672/TCP,25672/TCP,15672/TCP   80m
   enbuild-mongo               ClusterIP   10.43.25.176    <none>        27017/TCP                               80m
   enbuild-enbuild-user        ClusterIP   10.43.111.206   <none>        80/TCP                                  80m
   enbuild-enbuild-request     ClusterIP   10.43.5.125     <none>        80/TCP                                  80m
   enbuild-enbuild-ml          ClusterIP   10.43.227.97    <none>        80/TCP                                  80m
   enbuild-rabbitmq            ClusterIP   10.43.102.160   <none>        5672/TCP,4369/TCP,25672/TCP,15672/TCP   80m
   enbuild-enbuild-genai       ClusterIP   10.43.134.22    <none>        80/TCP                                  80m
   enbuild-enbuild-backend     ClusterIP   10.43.80.216    <none>        80/TCP                                  80m
   enbuild-enbuild-ui          ClusterIP   10.43.191.2     <none>        80/TCP                                  80m
   ```

5. **Access ENBUILD**\
    Use the following command to access the ENBUILD deployment using your local web browser.

   ```bash
     kubectl --namespace enbuild port-forward svc/enbuild-enbuild-ui 8080:80
   ```

   ```bash
   Forwarding from 127.0.0.1:8080 -> 8080
   Forwarding from [::1]:8080 -> 8080
   Handling connection for 8080
   Handling connection for 8080
   Handling connection for 8080
   Handling connection for 8080
   Handling connection for 8080
   Handling connection for 8080
   Handling connection for 8080
   Handling connection for 8080
   Handling connection for 8080

   ```

   Navigate your web browser to **127.0.0.1:8080**.

   <picture><img src="/images/deployEnbuildQuickstart/initial-login.png" alt="Screenshot of ENBUILD Login Screen"></img></picture>

6. **Uninstall ENBUILD**\
    Use the following command to uninstall ENBUILD from your Kubernetes cluster.

   ```bash
     helm uninstall enbuild -n enbuild
   ```

   ```bash
    release "enbuild" uninstalled
   ```
