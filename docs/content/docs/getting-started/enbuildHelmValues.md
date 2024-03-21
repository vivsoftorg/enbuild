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

A Helm chart for ENBUILD 1.0.5

## Helm Chart Dependencies

The ENBUILD Helm Chart has a dependency on the following Helm Charts.
| Repository | Name | Version |
| --------------------------------------- | ---------- | ------- |
| https://charts.bitnami.com/bitnami | rabbitmq | 11.13.0 |
| https://jupyterhub.github.io/helm-chart | jupyterhub | 2.0.0 |

## ENBUILD Helm Chart Values

The following key value pairs are used to configure ENBUILD.

### Global parameters

| Name                         | Description                                                            | Value                         |
| ---------------------------- | ---------------------------------------------------------------------- | ----------------------------- |
| `global.domain`              | What domain to use to expose the ENBUILD using istio or Ingress        | `ijuned.com`                  |
| `global.disable_tls_gitlab`  | Set to true if you are using self-signed certificates                  | `false`                       |
| `global.ingress.enabled`     | Should we create the Ingress Resources ?                               | `false`                       |
| `global.ingress.tls`         | Is Ingress TLS enabled ?                                               | `false`                       |
| `global.ingress.tls_secret`  | If Ingress is TLS enabled, Provide the Secret for the TLS Certificate. | `""`                          |
| `global.ingress.classname`   | Ingress classname if enabled.                                          | `""`                          |
| `global.ingress.annotations` | Ingress annotations if enabled.                                        | `[]`                          |
| `global.istio.enabled`       | Should we create the Istio Resources ?                                 | `false`                       |
| `global.istio.gateway`       | Istio gateway to use for creating Virtual Service.                     | `istio-system/public-gateway` |
| `global.auth_plugin`         | What Auth Plugin to use either of local,keycloak or okta.              | `local`                       |

### Global Keycloak parameters

| Name                                        | Description                                    | Value       |
| ------------------------------------------- | ---------------------------------------------- | ----------- |
| `global.keycloak`                           | Needed only if auth_plugin is set to keycloak. | `undefined` |
| `global.keycloak.keycloak_url`              | Keycloak URL                                   | `""`        |
| `global.keycloak.keycloak_backend_secret`   | Keycloak Client Secret                         | `""`        |
| `global.keycloak.keycloak_backend_username` | Keycloak Admin Username                        | `""`        |
| `global.keycloak.keycloak_backend_password` | Keycloak Admin Password                        | `""`        |

### Global OKTA parameters

| Name                        | Description                                | Value       |
| --------------------------- | ------------------------------------------ | ----------- |
| `global.okta`               | Needed only if auth_plugin is set to okta. | `undefined` |
| `global.okta.client_url`    | OKTA Client URL                            | `""`        |
| `global.okta.client_id`     | OKTA Client ID                             | `""`        |
| `global.okta.client_secret` | OKTA Client Secret                         | `""`        |
| `global.okta.client_token`  | OKTA Client TOKEN                          | `""`        |
| `global.okta.client_token`  | OKTA Client TOKEN                          | `""`        |
| `global.okta.base_url`      | OKTA Client Base URL                       | `""`        |

### Container Registry Parameters

| Name                        | Description                 | Value                 |
| --------------------------- | --------------------------- | --------------------- |
| `imageCredentials.registry` | Container registry Path     | `registry.gitlab.com` |
| `imageCredentials.username` | Container registry Username | `registry_user_name`  |
| `imageCredentials.password` | Container registry password | `registry_password`   |

### Jupyterhub Parameters

| Name                      | Description       | Value   |
| ------------------------- | ----------------- | ------- |
| `jupyterhub.cull.enabled` | Deploy Jupyterhub | `false` |

### ENBUILD RabbitMQ parameters

| Name                         | Description                                                          | Value          |
| ---------------------------- | -------------------------------------------------------------------- | -------------- |
| `rabbitmq.enabled`           | Set to false to use existing RabbitMQ                                | `true`         |
| `rabbitmq.replicaCount`      | RabbitMQ replicaCount                                                | `3`            |
| `rabbitmq.auth.username`     | RabbitMQ username                                                    | `juned`        |
| `rabbitmq.auth.password`     | RabbitMQ password                                                    | `memon`        |
| `rabbitmq.auth.erlangCookie` | RabbitMQ erlangCookie                                                | `lamba`        |
| `rabbitmq.host`              | If `rabbitmq.enabled` is false , provide the right rabbitmq endpoint | `""`           |
| `rabbitmq.queue_prefix`      | Queue Prefix for all RabbitMQ Queues                                 | `enbuild-prod` |

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

| Name                         | Description                                                                      | Value                                                                               |
| ---------------------------- | -------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------- |
| `enbuildUi.image.repository` | Container repository for enbuildUi                                               | `registry.gitlab.com/enbuild-staging/vivsoft-platform-ui/enbuild-frontend`          |
| `enbuildUi.image.pullPolicy` | Container image pullPolicy                                                       | `Always`                                                                            |
| `enbuildUi.image.tag`        | Container image tag. Skip to use the HelmChart appVersion as Image Tag           | `undefined`                                                                         |
| `enbuildUi.replicas`         | Container enbuildUI Replicas                                                     | `1`                                                                                 |
| `enbuildUi.service_type`     | enbuildUI service_type                                                           | `ClusterIP`                                                                         |
| `enbuildUi.node_port`        | enbuildUI node_port                                                              | `30080`                                                                             |
| `enbuildUi.hostname`         | enbuild service hostname. `enbuildUi.hostname`.`global.domain` becomes your FQDN | `enbuild`                                                                           |
| `enbuildUi.kiali_url`        | kiali_url                                                                        | `https://kiali.vivplatform.io/kiali/`                                               |
| `enbuildUi.grafana_url`      | grafana_url                                                                      | `https://grafana.vivplatform.io/d/os6Bh8Omk/kubernetes-cluster?orgId=1&refresh=30s` |
| `enbuildUi.loki_url`         | loki_url                                                                         | `https://grafana.vivplatform.io/d/liz0yRCZz/logs-app?orgId=1`                       |
| `enbuildUi.kubecost_url`     | kubecost_url                                                                     | `https://kubecost.vivplatform.io/overview.html`                                     |

### ENBUILD Backend Services parameters

| Name                         | Description                                                            | Value                                                                     |
| ---------------------------- | ---------------------------------------------------------------------- | ------------------------------------------------------------------------- |
| `enbuildBk.image.repository` | Container repository for enbuildBk                                     | `registry.gitlab.com/enbuild-staging/vivsoft-platform-ui/enbuild-backend` |
| `enbuildBk.image.pullPolicy` | Container image pullPolicy                                             | `Always`                                                                  |
| `enbuildUi.image.tag`        | Container image tag. Skip to use the HelmChart appVersion as Image Tag | `undefined`                                                               |
| `enbuildBk.replicas`         | Container enbuildBk Replicas                                           | `1`                                                                       |
| `enbuildBk.service_type`     | enbuildBk service_type                                                 | `ClusterIP`                                                               |
| `enbuildBk.encryption_key`   | encryption_key to be used by Backend                                   | `encryption_key`                                                          |

### ENBUILD ML Services parameters

| Name                         | Description                                                            | Value                                                                |
| ---------------------------- | ---------------------------------------------------------------------- | -------------------------------------------------------------------- |
| `enbuildMl.image.repository` | Container repository for enbuildMl                                     | `registry.gitlab.com/enbuild-staging/vivsoft-platform-ui/enbuild-ml` |
| `enbuildMl.image.pullPolicy` | Container image pullPolicy                                             | `Always`                                                             |
| `enbuildUi.image.tag`        | Container image tag. Skip to use the HelmChart appVersion as Image Tag | `undefined`                                                          |
| `enbuildMl.replicas`         | Container enbuildMl Replicas                                           | `1`                                                                  |
| `enbuildMl.service_type`     | enbuildMl service_type                                                 | `ClusterIP`                                                          |

### ENBUILD GenAI Services parameters

| Name                            | Description                                                            | Value                                                                   |
| ------------------------------- | ---------------------------------------------------------------------- | ----------------------------------------------------------------------- |
| `enbuildGenAI.image.repository` | Container repository for enbuildGenAI                                  | `registry.gitlab.com/enbuild-staging/vivsoft-platform-ui/enbuild-genai` |
| `enbuildGenAI.image.pullPolicy` | Container image pullPolicy                                             | `Always`                                                                |
| `enbuildUi.image.tag`           | Container image tag. Skip to use the HelmChart appVersion as Image Tag | `undefined`                                                             |
| `enbuildGenAI.replicas`         | Container enbuildGenAI Replicas                                        | `1`                                                                     |
| `enbuildGenAI.service_type`     | enbuildGenAI service_type                                              | `ClusterIP`                                                             |
| `enbuildGenAI.api_key`          | api_key for OpenAI service.                                            | `dummy`                                                                 |

### ENBUILD Request Services parameters

| Name                              | Description                                                            | Value                                                                     |
| --------------------------------- | ---------------------------------------------------------------------- | ------------------------------------------------------------------------- |
| `enbuildRequest.image.repository` | Container repository for enbuildRequest                                | `registry.gitlab.com/enbuild-staging/vivsoft-platform-ui/enbuild-request` |
| `enbuildRequest.image.pullPolicy` | Container image pullPolicy                                             | `Always`                                                                  |
| `enbuildUi.image.tag`             | Container image tag. Skip to use the HelmChart appVersion as Image Tag | `undefined`                                                               |
| `enbuildRequest.replicas`         | Container enbuildRequest Replicas                                      | `1`                                                                       |
| `enbuildRequest.service_type`     | enbuildRequest service_type                                            | `ClusterIP`                                                               |

### ENBUILD USER Services parameters

| Name                           | Description                                                            | Value                                                                  |
| ------------------------------ | ---------------------------------------------------------------------- | ---------------------------------------------------------------------- |
| `enbuildUser.image.repository` | Container repository for enbuildUser                                   | `registry.gitlab.com/enbuild-staging/vivsoft-platform-ui/enbuild-user` |
| `enbuildUser.image.pullPolicy` | Container image pullPolicy                                             | `Always`                                                               |
| `enbuildUi.image.tag`          | Container image tag. Skip to use the HelmChart appVersion as Image Tag | `undefined`                                                            |
| `enbuildUser.replicas`         | Container enbuildUser Replicas                                         | `1`                                                                    |
| `enbuildUser.service_type`     | enbuildUser service_type                                               | `ClusterIP`                                                            |

### ENBUILD Sync Services parameters

| Name                           | Description                                                            | Value                                                                     |
| ------------------------------ | ---------------------------------------------------------------------- | ------------------------------------------------------------------------- |
| `enbuildSync.image.repository` | Container repository for enbuildSync                                   | `registry.gitlab.com/enbuild-staging/vivsoft-platform-ui/enbuild-cronjob` |
| `enbuildSync.image.pullPolicy` | Container image pullPolicy                                             | `Always`                                                                  |
| `enbuildUi.image.tag`          | Container image tag. Skip to use the HelmChart appVersion as Image Tag | `undefined`                                                               |
| `enbuildSync.replicas`         | Container enbuildSync Replicas                                         | `1`                                                                       |
