---
title: "Helm Values"
description: "Various configuration details for deploying ENBUILD using Helm."
summary: ""
date: 2023-09-07T16:04:48+02:00
lastmod: 2023-09-07T16:04:48+02:00
draft: false
menu:
  docs:
    parent: "docs"
    identifier: "enbuildHelmValues"
weight: 6
toc: true
seo:
  title: "" # custom title (optional)
  description: "" # custom description (recommended)
  canonical: "" # custom canonical URL (optional)
  noindex: false # false (default) or true
---

# ENBUILD Helm Chart Values

The following key value pairs are used to configure ENBUILD.

## Parameters

### Global parameters

| Name                         | Description                                                                                                                                                                                      | Value               |
| ---------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ | ------------------- |
| `global.AppVersion`          | [default: ""] Provide custom appVersion, to override the default one. All the ENBUILD images will be of the same version. To use indidual tag for each service set the tag on per service basis. | `""`                |
| `global.domain`              | What domain to use to expose the ENBUILD using istio or Ingress                                                                                                                                  | `ijuned.com`        |
| `global.disable_tls_gitlab`  | Set to true if you are using self-signed certificates                                                                                                                                            | `false`             |
| `global.ingress.enabled`     | Should we create the Ingress Resources ?                                                                                                                                                         | `false`             |
| `global.ingress.tls`         | Is Ingress TLS enabled ?                                                                                                                                                                         | `false`             |
| `global.ingress.tls_secret`  | If Ingress is TLS enabled, Provide the Secret for the TLS Certificate.                                                                                                                           | `""`                |
| `global.ingress.classname`   | Ingress classname if enabled.                                                                                                                                                                    | `""`                |
| `global.ingress.annotations` | Ingress annotations if enabled.                                                                                                                                                                  | `[]`                |
| `global.istio.enabled`       | Should we create the Istio Resources ?                                                                                                                                                           | `false`             |
| `global.istio.gateway`       | Istio gateway to use for creating Virtual Service.                                                                                                                                               | `istio-system/main` |

### Container Registry Parameters

| Name                        | Description                                                                         | Value |
| --------------------------- | ----------------------------------------------------------------------------------- | ----- |
| `imageCredentials`          | Should we use a private container registry? if yes provide the following parameters | `{}`  |
| `imageCredentials.registry` | Container registry Path                                                             | `""`  |
| `imageCredentials.username` | Container registry Username                                                         | `""`  |
| `imageCredentials.password` | Container registry password                                                         | `""`  |

### Jupyterhub Parameters

| Name                      | Description        | Value   |
| ------------------------- | ------------------ | ------- |
| `jupyterhub.cull.enabled` | Deploy Jupyterhub  | `false` |

### ENBUILD RabbitMQ parameters

| Name                         | Description                                                          | Value         |
| ---------------------------- | -------------------------------------------------------------------- | ------------- |
| `rabbitmq.enabled`           | Set to false to use existing RabbitMQ                                | `true`        |
| `rabbitmq.replicaCount`      | RabbitMQ replicaCount                                                | `1`           |
| `rabbitmq.auth.username`     | RabbitMQ username                                                    | `admin`       |
| `rabbitmq.auth.password`     | RabbitMQ password                                                    | `SuperSecret` |
| `rabbitmq.auth.erlangCookie` | RabbitMQ erlangCookie                                                | `lamba`       |
| `rabbitmq.host`              | If `rabbitmq.enabled` is false , provide the right rabbitmq endpoint | `""`          |
| `rabbitmq.queue_prefix`      | Queue Prefix for all RabbitMQ Queues                                 | `enbuild`     |

### ENBUILD Backend/DB parameters

| Name                          | Description                                                                                                                    | Value                                                                               |
| ----------------------------- | ------------------------------------------------------------------------------------------------------------------------------ | ----------------------------------------------------------------------------------- |
| `mongodb.enabled`             | Set to true to Deploy the MongoDB.                                                                                             | `false`                                                                             |
| `mongodb.mongo_root_username` | DB username. If `mongodb.enabled` this is used to to set the username. Else this is username for existing Cosmos or DocumentDB | `""`                                                                                |
| `mongodb.mongo_root_password` | DB Password. If `mongodb.enabled` this is used to to set the password. Else this is password for existing Cosmos or DocumentDB | `""`                                                                                |
| `mongodb.mongo_server`        | If `mongodb.enabled` is false , provide the right cosmosDB/DocumentDB endpoint                                                 | `""`                                                                                |
| `mongodb.image.repository`    | Container repository for mongodb Container                                                                                     | `registry.gitlab.com/enbuild-staging/container_registry/opensource-mongodb/mongodb` |
| `mongodb.image.tag`           | Container tag for mongodb Container                                                                                            | `4.4.5`                                                                             |
| `mongodb.image.pullPolicy`    | Container ImagePullPolicy for mongodb Container                                                                                | `Always`                                                                            |

### ENBUILD UI Services parameters

| Name                         | Description                                                                      | Value                                                                      |
| ---------------------------- | -------------------------------------------------------------------------------- | -------------------------------------------------------------------------- |
| `enbuildUi.image.repository` | Container repository for enbuildUi                                               | `registry.gitlab.com/enbuild-staging/vivsoft-platform-ui/enbuild-frontend` |
| `enbuildUi.image.pullPolicy` | Container image pullPolicy                                                       | `Always`                                                                   |
| `enbuildUi.image.tag`        | Container image tag. Skip to use the HelmChart appVersion as Image Tag           | `undefined`                                                                |
| `enbuildUi.replicas`         | Container enbuildUI Replicas                                                     | `1`                                                                        |
| `enbuildUi.service_type`     | enbuildUI service_type                                                           | `ClusterIP`                                                                |
| `enbuildUi.node_port`        | enbuildUI node_port                                                              | `30080`                                                                    |
| `enbuildUi.hostname`         | enbuild service hostname. `enbuildUi.hostname`.`global.domain` becomes your FQDN | `enbuild`                                                                  |
| `enbuildUi.kiali_url`        | kiali_url                                                                        | `https://kiali.ijuned.com/kiali/`                                          |
| `enbuildUi.grafana_url`      | grafana_url                                                                      | `https://grafana.ijuned.com/`                                              |
| `enbuildUi.loki_url`         | loki_url                                                                         | `https://grafana.ijuned.com/d/liz0yRCZz/logs-app?orgId=1`                  |
| `enbuildUi.kubecost_url`     | kubecost_url                                                                     | `https://kubecost.ijuned.com/overview.html`                                |

### ENBUILD Backend Services parameters

| Name                         | Description                                                            | Value                                                                     |
| ---------------------------- | ---------------------------------------------------------------------- | ------------------------------------------------------------------------- |
| `enbuildBk.image.repository` | Container repository for enbuildBk                                     | `registry.gitlab.com/enbuild-staging/vivsoft-platform-ui/enbuild-backend` |
| `enbuildBk.image.pullPolicy` | Container image pullPolicy                                             | `Always`                                                                  |
| `enbuildBk.image.tag`        | Container image tag. Skip to use the HelmChart appVersion as Image Tag | `undefined`                                                               |
| `enbuildBk.replicas`         | Container enbuildBk Replicas                                           | `1`                                                                       |
| `enbuildBk.service_type`     | enbuildBk service_type                                                 | `ClusterIP`                                                               |
| `enbuildBk.encryption_key`   | encryption_key to be used by Backend                                   | `encryption_key`                                                          |

### ENBUILD USER Services parameters

| Name                           | Description                                                            | Value                                                                  |
| ------------------------------ | ---------------------------------------------------------------------- | ---------------------------------------------------------------------- |
| `enbuildUser.image.repository` | Container repository for enbuildUser                                   | `registry.gitlab.com/enbuild-staging/vivsoft-platform-ui/enbuild-user` |
| `enbuildUser.image.pullPolicy` | Container image pullPolicy                                             | `Always`                                                               |
| `enbuildUser.image.tag`        | Container image tag. Skip to use the HelmChart appVersion as Image Tag | `undefined`                                                            |
| `enbuildUser.replicas`         | Container enbuildUser Replicas                                         | `1`                                                                    |
| `enbuildUser.service_type`     | enbuildUser service_type                                               | `ClusterIP`                                                            |

### ENBUILD Sync Services parameters

| Name                           | Description                                                            | Value                                                                     |
| ------------------------------ | ---------------------------------------------------------------------- | ------------------------------------------------------------------------- |
| `enbuildSync.image.repository` | Container repository for enbuildSync                                   | `registry.gitlab.com/enbuild-staging/vivsoft-platform-ui/enbuild-cronjob` |
| `enbuildSync.image.pullPolicy` | Container image pullPolicy                                             | `Always`                                                                  |
| `enbuildSync.image.tag`        | Container image tag. Skip to use the HelmChart appVersion as Image Tag | `undefined`                                                               |
| `enbuildSync.replicas`         | Container enbuildSync Replicas                                         | `1`                                                                       |

### ENBUILD ML Services parameters

| Name                         | Description                                                            | Value                                                                |
| ---------------------------- | ---------------------------------------------------------------------- | -------------------------------------------------------------------- |
| `enbuildMl.enabled`          | Should we create the ENBUILD ML microservice ?                         | `false`                                                              |
| `enbuildMl.image.repository` | Container repository for enbuildMl                                     | `registry.gitlab.com/enbuild-staging/vivsoft-platform-ui/enbuild-ml` |
| `enbuildMl.image.pullPolicy` | Container image pullPolicy                                             | `Always`                                                             |
| `enbuildMl.image.tag`        | Container image tag. Skip to use the HelmChart appVersion as Image Tag | `undefined`                                                          |
| `enbuildMl.replicas`         | Container enbuildMl Replicas                                           | `1`                                                                  |
| `enbuildMl.service_type`     | enbuildMl service_type                                                 | `ClusterIP`                                                          |

### ENBUILD GenAI Services parameters

| Name                            | Description                                                            | Value                                                                   |
| ------------------------------- | ---------------------------------------------------------------------- | ----------------------------------------------------------------------- |
| `enbuildGenAI.enabled`          | Should we create the ENBUILD GenAI microservice ?                      | `false`                                                                 |
| `enbuildGenAI.image.repository` | Container repository for enbuildGenAI                                  | `registry.gitlab.com/enbuild-staging/vivsoft-platform-ui/enbuild-genai` |
| `enbuildGenAI.image.pullPolicy` | Container image pullPolicy                                             | `Always`                                                                |
| `enbuildGenAI.image.tag`        | Container image tag. Skip to use the HelmChart appVersion as Image Tag | `undefined`                                                             |
| `enbuildGenAI.replicas`         | Container enbuildGenAI Replicas                                        | `1`                                                                     |
| `enbuildGenAI.service_type`     | enbuildGenAI service_type                                              | `ClusterIP`                                                             |
| `enbuildGenAI.api_key`          | api_key for OpenAI service.                                            | `dummy`                                                                 |

### ENBUILD Request Services parameters

| Name                              | Description                                                            | Value                                                                     |
| --------------------------------- | ---------------------------------------------------------------------- | ------------------------------------------------------------------------- |
| `enbuildRequest.enabled`          | Should we create the ENBUILD Request microservice ?                    | `false`                                                                   |
| `enbuildRequest.image.repository` | Container repository for enbuildRequest                                | `registry.gitlab.com/enbuild-staging/vivsoft-platform-ui/enbuild-request` |
| `enbuildRequest.image.pullPolicy` | Container image pullPolicy                                             | `Always`                                                                  |
| `enbuildRequest.image.tag`        | Container image tag. Skip to use the HelmChart appVersion as Image Tag | `undefined`                                                               |
| `enbuildRequest.replicas`         | Container enbuildRequest Replicas                                      | `1`                                                                       |
| `enbuildRequest.service_type`     | enbuildRequest service_type                                            | `ClusterIP`                                                               |
