---
title: "ENBUILD Helm Chart Installation in Airgapped Kubernetes Environment"
description: "Steps to install ENBUILD in Airgapped Kubernetes Environment"
summary: "Steps to install ENBUILD in Airgapped Kubernetes Environment"
draft: false
menu:
  docs:
    parent: "docs/how-to-guides/"
    identifier: "enbuildHaulerAirGap"
weight: 202
toc: true
seo:
  title: "" # custom title (optional)
  description: "" # custom description (recommended)
  canonical: "" # custom canonical URL (optional)
  noindex: false # false (default) or true
---

# ENBUILD Helm Chart Installation in Airgapped Kubernetes Environment

This document provides a step-by-step guide for deploying the ENBUILD Helm chart in a disconnected (airgapped) Kubernetes environment using an ENBUILD haul file.

## Prerequisites
- Kubernetes cluster (EKS, RKE2, or K3s) installed in the airgapped environment.
- A micro instance or one of the nodes to act as the Hauler registry.
- Access to the ENBUILD haul package (how to transfer the haul file to the airgapped environment is beyond the scope of this document).

## Steps

### 1. Install Kubernetes in the Airgapped Environment
Install the Kubernetes environment in the client's disconnected environment. Follow the installation guide for the respective Kubernetes distribution:
- [EKS Installation Guide](https://docs.aws.amazon.com/eks/latest/userguide/install-kubectl.html)
- [RKE2 Installation Guide](https://docs.rke2.io/install/)
- [K3s Installation Guide](https://docs.k3s.io/installation/)

### 2. Create a Micro Instance for the Hauler Registry
Set up a micro instance in the airgapped environment to host the Hauler registry. For K3s or RKE2, this can be one of the existing nodes.

### 3. Install Hauler, kubectl, and Helm Binaries on the Registry Node
Install the necessary binaries on the node that will act as the Hauler registry:

- **Hauler**: [Hauler Installation](https://docs.hauler.dev)
- **kubectl**: [Install kubectl](https://kubernetes.io/docs/tasks/tools/)
- **Helm**: [Install Helm](https://helm.sh/docs/intro/install/)

### 4. Download the ENBUILD Haul File
The ENBUILD haul file is available as a [public artifact](https://enbuild-haul.s3.us-east-1.amazonaws.com/enbuild-0.0.20.tar.zst) in the Vivsoft AWS S3 bucket. This file should be downloaded and transferred to a client-side S3 bucket (how to transfer the file is beyond the scope of this document). Once the file is available in the client’s side S3 bucket, download it on the registry server using the following command:

```bash
ENBUILD_HELM_CHART_VERSION=0.0.20
ENBUILD_HAULER_URL="https://your-client-side-s3-bucket/enbuild-${ENBUILD_HELM_CHART_VERSION}.tar.zst"
curl -O ${ENBUILD_HAULER_URL}
```

### 5. Load the ENBUILD Haul File and Start the Registry

Once the haul file is available in the airgapped environment, load the haul file into the Hauler registry and start the registry service:

```bash
hauler store load enbuild-${ENBUILD_HELM_CHART_VERSION}.tar.zst
hauler store serve registry &
sleep 5
```

This command will start a registry that serves the images from the loaded haul file.

### 6. Configure Kubernetes to Use the Private Registry
Configure your Kubernetes environment to use the private Hauler registry for pulling images. This may involve setting `insecure_skip_verify` for the node’s private IP address (`<PRIVATE_IP_LOCAL_REGISTRY_SERVER>:5000`).

Refer to the appropriate guide for your environment:

- **K3s**: [K3s Private Registry Setup](https://docs.k3s.io/installation/private-registry)
- **RKE2**: [RKE2 Containerd Registry Configuration](https://docs.rke2.io/install/containerd_registry_configuration)
- **EKS/Generic**: [Containerd Registry Setup](https://github.com/containerd/containerd/blob/main/docs/cri/registry.md)

#### Example: `registries.yaml` for K3s , with `192.168.0.111` as the PRIVATE_IP_LOCAL_REGISTRY_SERVER

```yaml
mirrors:
  "192.168.0.111:5000":
    endpoint:
      - "http://192.168.0.111:5000"
  "registry.gitlab.com":
    endpoint:
      - "http://192.168.0.111:5000"

configs:
  "192.168.0.111:5000":
    tls:
      insecure_skip_verify: true
```

Restart the K3s or RKE2 service to apply the registry configuration.

### 7. Create Configuration Files for ENBUILD Installation
Create the necessary Helm configuration file for ENBUILD installation:

```bash
cat <<EOF > quick_install_hauler.yaml
global:
  image:
    registry: ${PRIVATE_IP_LOCAL_REGISTRY_SERVER}:5000
    pullPolicy: Always
rabbitmq:
  image:
    registry: ${PRIVATE_IP_LOCAL_REGISTRY_SERVER}:5000
    repository: bitnami/rabbitmq
    tag: 3.11.13-debian-11-r0
EOF
```

### 8. Install the ENBUILD Helm Chart
Finally, install the ENBUILD Helm chart using the private registry:

```bash
helm upgrade --install --namespace enbuild enbuild --plain-http oci://${PRIVATE_IP_LOCAL_REGISTRY_SERVER}:5000/hauler/enbuild --version ${ENBUILD_HELM_CHART_VERSION} -f quick_install_hauler.yaml --create-namespace
```

This command will install the ENBUILD release in the `enbuild` namespace of your Kubernetes cluster.

---

This completes the setup of the ENBUILD Helm chart in an airgapped environment using a local Hauler registry.
