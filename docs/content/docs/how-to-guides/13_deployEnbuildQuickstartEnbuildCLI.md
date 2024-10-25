---
title: "Deploying ENBUILD using ENBUILD CLI"
description: "Steps to Configure ENBUILD CLI and deploy ENBUILD using ENBUILD CLI"
summary: "Configure ENBUILD CLI and deploy ENBUILD using ENBUILD CLI"
draft: false
menu:
  docs:
    parent: "docs/how-to-guides/"
    identifier: "configureEnbuildCLI"
weight: 202
toc: true
seo:
  title: "" # custom title (optional)
  description: "" # custom description (recommended)
  canonical: "" # custom canonical URL (optional)
  noindex: false # false (default) or true
---

## Configure ENBUILD CLI

Follow these step-by-step instructions to configure ENBUILD CLI.

### Prerequisites

Make sure you install the following dependencies.

1. [Docker](https://docs.docker.com/engine/install/)
    - Install docker by following these [steps](https://docs.docker.com/engine/install/).
    - Make sure that docker engine is running before going using the Enbuild CLI.

2. [yq cli](https://mikefarah.gitbook.io/yq)
    - Install yq cli following these [steps](https://github.com/mikefarah/yq/#install).
    - Enbuild cli is using it internally for creating bigbang catalog template values file.

**Note:** Additional dependencies such as k3d, Helm, jq, grep, sed, curl, and iproute2 will be checked and installed by the ENBUILD CLI when executing the commands.


### Configuration

1. Download the ENBUILD CLI binary compatible with your operating system from this [link](https://github.com/vivsoftorg/enbuild/releases/tag/v0.0.11)

2. Extract the downloaded folder

3. Add the `enbuild` command to the PATH environment variable

    ```bash
    export PATH=$PATH:<path-to-the-above-extracted-enbuild-directory>
    ```

3. Verify that `enbuild` cli is ready to use by running these commands

    ```bash
    enbuild -v
    ```
4. For more information on enbuild cli commands, please run

    ```bash
    enbuild -h
    ```

## Deploy ENBUILD using ENBUILD CLI

### Deployment Steps

To Create a k3d kubernetes cluster with ENBUILD installed on your local machine, run the command

  ```bash
  enbuild demo up
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

Use the port forwarding command to access the ENBUILD UI using your web browser.

```bash
kubectl --namespace enbuild port-forward svc/enbuild-enbuild-ui 3000:80

Forwarding from 127.0.0.1:3000 -> 8080
Forwarding from [::1]:3000 -> 8080
```

Navigate your web browser to **http://127.0.0.1:3000**. and set the admin password.

<picture><img src="/images/getting-started/initial-login.png" alt="Screenshot of ENBUILD Login Screen"></img></picture>

After you set the initial admin password, you should see the ENBUILD home page with BigBang Catalog.

<picture><img src="/images/getting-started/enbuild_home_page_first_login.png" alt="Screenshot of ENBUILD Home Screen"></img></picture>


:zap: ***[Proceed to Configureing ENBUILD](../configuring-enbuild/)***

### Uninstall ENBUILD using ENBUILD CLI

To Uninstall ENBUILD on local k3d cluster and stop the k3d cluster on your local machine

  ```bash
  enbuild demo down
  ```

### Destroy the k3d cluster using ENBUILD CLI

To completely Remove k3d cluster with ENBUILD installed on your local machine

  ```bash
  enbuild demo destroy
  ```

