# [ENBUILD HELM CHART](https://github.com/vivsoftorg/enbuild_helm_chart.git)

This helm chart installs the [ENBUILD application](https://gitlab.com/enbuild-staging/vivsoft-platform-ui).

# Installing the Chart

This Helm chart repository enables you to install a ENBUILD
Helm chart directly from it into your Kubernetes cluster. Please refer to the
[ENBUILD documentation](https://vivsoftorg.github.io/enbuild/) for all
the additional details required.

```shell
# Let helm the command line tool know about a Helm chart repository
# that we decide to name enbuild.
❯ helm repo add vivsoft https://vivsoftorg.github.io/enbuild

# Update the Helm chart repository.
❯ helm repo update vivsoft

# Search for the ENBUILD Helm chart in the enbuild Helm chart repository.
❯ helm search repo  vivsoft/enbuild
NAME           	CHART VERSION	APP VERSION	DESCRIPTION
vivsoft/enbuild	0.0.12        	1.0.10      	A Helm chart for ENBUILD

# Simplified example on how to install a Helm chart from a Helm chart repository
# named vivsoft in a namespace named enbuild. See the Helm chart's documentation for additional details
# required.
❯ helm upgrade --install  enbuild vivsoft/enbuild --namespace enbuild --create-namespace 

# To install a specific version of the Helm chart.
❯ helm upgrade --install  enbuild vivsoft/enbuild --namespace enbuild --create-namespace  --version 0.0.12
```

# Uninstalling the Chart

To uninstall/delete the `enbuild` deployment:

```bash
❯ helm delete --namespace enbuild enbuild
```

## Parameters

### Global parameters

| Name                                         | Description                                                                                                                                                                                      | Value                 |
| -------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ | --------------------- |
| `global.AppVersion`                          | [default: ""] Provide custom appVersion, to override the default one. All the ENBUILD images will be of the same version. To use indidual tag for each service set the tag on per service basis. | `""`                  |
| `global.domain`                              | What domain to use to expose the ENBUILD using istio or Ingress                                                                                                                                  | `ijuned.com`          |
| `global.disable_tls_gitlab`                  | Set to true if you are using self-signed certificates                                                                                                                                            | `false`               |
| `global.ingress.enabled`                     | Should we create the Ingress Resources ?                                                                                                                                                         | `false`               |
| `global.ingress.tls`                         | Is Ingress TLS enabled ?                                                                                                                                                                         | `false`               |
| `global.ingress.tls_secret`                  | If Ingress is TLS enabled, Provide the Secret for the TLS Certificate.                                                                                                                           | `""`                  |
| `global.ingress.classname`                   | Ingress classname if enabled.                                                                                                                                                                    | `""`                  |
| `global.ingress.annotations`                 | Ingress annotations if enabled.                                                                                                                                                                  | `[]`                  |
| `global.istio.enabled`                       | Should we create the Istio Resources ?                                                                                                                                                           | `false`               |
| `global.istio.gateway`                       | Istio gateway to use for creating Virtual Service.                                                                                                                                               | `istio-system/main`   |
| `global.image.registry`                      | Container registry to pull images from                                                                                                                                                           | `registry.gitlab.com` |
| `global.image.pullPolicy`                    | Container imagePullPolicy                                                                                                                                                                        | `Always`              |
| `global.image.registry_credentials`          | if the image.registry is private container registry, provide the credentials                                                                                                                     | `{}`                  |
| `global.image.registry_credentials.username` | Container registry Username                                                                                                                                                                      | `""`                  |
| `global.image.registry_credentials.password` | Container registry password                                                                                                                                                                      | `""`                  |
| `global.monitoring.enabled`                  | Should we install loki-stack for monitoring , if yes set to true , set the lok-stack values in the values.yaml                                                                                   | `true`                |

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

| Name                          | Description                                                                                                                    | Value                                         |
| ----------------------------- | ------------------------------------------------------------------------------------------------------------------------------ | --------------------------------------------- |
| `mongodb.enabled`             | Set to true to Deploy the MongoDB.                                                                                             | `false`                                       |
| `mongodb.mongo_root_username` | DB username. If `mongodb.enabled` this is used to to set the username. Else this is username for existing Cosmos or DocumentDB | `""`                                          |
| `mongodb.mongo_root_password` | DB Password. If `mongodb.enabled` this is used to to set the password. Else this is password for existing Cosmos or DocumentDB | `""`                                          |
| `mongodb.mongo_server`        | If `mongodb.enabled` is false , provide the right cosmosDB/DocumentDB endpoint                                                 | `""`                                          |
| `mongodb.image.repository`    | Container repository for mongodb Container                                                                                     | `enbuild-staging/vivsoft-platform-ui/mongodb` |
| `mongodb.image.tag`           | Container tag for mongodb Container                                                                                            | `4.4.5`                                       |

### ENBUILD UI Services parameters

| Name                         | Description                                                                      | Value                                                         |
| ---------------------------- | -------------------------------------------------------------------------------- | ------------------------------------------------------------- |
| `enbuildUi.image.repository` | Container repository for enbuildUi                                               | `enbuild-staging/vivsoft-platform-ui/enbuild-frontend`        |
| `enbuildUi.image.tag`        | Container image tag. Skip to use the HelmChart appVersion as Image Tag           | `undefined`                                                   |
| `enbuildUi.replicas`         | Container enbuildUI Replicas                                                     | `1`                                                           |
| `enbuildUi.service_type`     | enbuildUI service_type                                                           | `ClusterIP`                                                   |
| `enbuildUi.node_port`        | enbuildUI node_port                                                              | `30080`                                                       |
| `enbuildUi.hostname`         | enbuild service hostname. `enbuildUi.hostname`.`global.domain` becomes your FQDN | `enbuild`                                                     |
| `enbuildUi.kiali_url`        | kiali_url                                                                        | `/kiali/`                                                     |
| `enbuildUi.grafana_url`      | grafana_url                                                                      | `/grafana/d/os6Bh8Omk/kubernetes-cluster?orgId=1&refresh=30s` |
| `enbuildUi.loki_url`         | loki_url                                                                         | `/grafana/d/liz0yRCZz/logs-app?orgId=1`                       |
| `enbuildUi.kubecost_url`     | kubecost_url                                                                     | `kubecost/overview.html`                                      |

### ENBUILD Backend Services parameters

| Name                         | Description                                                            | Value                                                 |
| ---------------------------- | ---------------------------------------------------------------------- | ----------------------------------------------------- |
| `enbuildBk.image.repository` | Container repository for enbuildBk                                     | `enbuild-staging/vivsoft-platform-ui/enbuild-backend` |
| `enbuildBk.image.tag`        | Container image tag. Skip to use the HelmChart appVersion as Image Tag | `undefined`                                           |
| `enbuildBk.replicas`         | Container enbuildBk Replicas                                           | `1`                                                   |
| `enbuildBk.service_type`     | enbuildBk service_type                                                 | `ClusterIP`                                           |
| `enbuildBk.encryption_key`   | encryption_key to be used by Backend                                   | `encryption_key`                                      |

### ENBUILD USER Services parameters

| Name                           | Description                                                            | Value                                              |
| ------------------------------ | ---------------------------------------------------------------------- | -------------------------------------------------- |
| `enbuildUser.image.repository` | Container repository for enbuildUser                                   | `enbuild-staging/vivsoft-platform-ui/enbuild-user` |
| `enbuildUser.image.tag`        | Container image tag. Skip to use the HelmChart appVersion as Image Tag | `undefined`                                        |
| `enbuildUser.replicas`         | Container enbuildUser Replicas                                         | `1`                                                |
| `enbuildUser.service_type`     | enbuildUser service_type                                               | `ClusterIP`                                        |

### ENBUILD MQ Consumer Services parameters

| Name                               | Description                                                            | Value                                                     |
| ---------------------------------- | ---------------------------------------------------------------------- | --------------------------------------------------------- |
| `enbuildConsumer.image.repository` | Container repository for enbuildConsumer                               | `enbuild-staging/vivsoft-platform-ui/enbuild-mq-consumer` |
| `enbuildConsumer.image.tag`        | Container image tag. Skip to use the HelmChart appVersion as Image Tag | `undefined`                                               |
| `enbuildConsumer.replicas`         | Container enbuildConsumer Replicas                                     | `1`                                                       |

### ENBUILD ML Services parameters

| Name                         | Description                                                                                            | Value                                            |
| ---------------------------- | ------------------------------------------------------------------------------------------------------ | ------------------------------------------------ |
| `enbuildMl.enabled`          | Should we create the ENBUILD ML microservice, which also controls whether or not to install jupyterhub | `false`                                          |
| `enbuildMl.image.repository` | Container repository for enbuildMl                                                                     | `enbuild-staging/vivsoft-platform-ui/enbuild-ml` |
| `enbuildMl.image.tag`        | Container image tag. Skip to use the HelmChart appVersion as Image Tag                                 | `undefined`                                      |
| `enbuildMl.replicas`         | Container enbuildMl Replicas                                                                           | `1`                                              |
| `enbuildMl.service_type`     | enbuildMl service_type                                                                                 | `ClusterIP`                                      |

### ENBUILD GenAI Services parameters

| Name                            | Description                                                                                            | Value                                               |
| ------------------------------- | ------------------------------------------------------------------------------------------------------ | --------------------------------------------------- |
| `enbuildGenAI.enabled`          | Should we create the ENBUILD GenAI microservice, which also controls whether or not to install Ollama. | `false`                                             |
| `enbuildGenAI.image.repository` | Container repository for enbuildGenAI                                                                  | `enbuild-staging/vivsoft-platform-ui/enbuild-genai` |
| `enbuildGenAI.image.tag`        | Container image tag. Skip to use the HelmChart appVersion as Image Tag                                 | `undefined`                                         |
| `enbuildGenAI.replicas`         | Container enbuildGenAI Replicas                                                                        | `1`                                                 |
| `enbuildGenAI.service_type`     | enbuildGenAI service_type                                                                              | `ClusterIP`                                         |
| `enbuildGenAI.api_key`          | api_key [default: "dummy"] for OpenAI service.                                                         | `dummy`                                             |
| `enbuildGenAI.model_name`       | model_name for OpenAI service.                                                                         | `"ollama/llama3.1"`                                 |
| `enbuildGenAI.ollama_endpoint`  | ollama_endpoint for OpenAI service.                                                                    | `"http://open-webui-ollama:11434"`                  |

### ENBUILD AI Services parameters

| Name                         | Description                                                                                         | Value                                            |
| ---------------------------- | --------------------------------------------------------------------------------------------------- | ------------------------------------------------ |
| `enbuildAI.enabled`          | Should we create the ENBUILD AI microservice, which also controls whether or not to install Ollama. | `false`                                          |
| `enbuildAI.image.repository` | Container repository for enbuildAI                                                                  | `enbuild-staging/vivsoft-platform-ui/enbuild-ai` |
| `enbuildAI.image.tag`        | Container image tag. Skip to use the HelmChart appVersion as Image Tag                              | `undefined`                                      |
| `enbuildAI.replicas`         | Container enbuilAI Replicas                                                                         | `1`                                              |
| `enbuildAI.service_type`     | enbuildAI service_type                                                                              | `ClusterIP`                                      |

<!-- # ----- -->