---
title: "Helm"
description: "Reference for Helm"
summary: ""
draft: false
menu:
  docs:
    parent: "/reference/"
    identifier: "helm"
weight: 504
toc: true
seo:
  title: "" # custom title (optional)
  description: "" # custom description (recommended)
  canonical: "" # custom canonical URL (optional)
  noindex: false # false (default) or true
---

<picture><img src="/images/references/helm.svg" alt="Screenshot of Helm Logo" width="20%" height="20%"></img></picture>
<br />
<br />
The Helm CLI serves as a powerful package manager for Kubernetes applications, streamlining and simplifying the deployment and management of complex containerized applications. At its core, Helm utilizes charts—pre-configured packages of Kubernetes resources—to encapsulate and version entire applications. With the Helm CLI, users can effortlessly install, upgrade, and roll back applications, ensuring consistent and reproducible deployments across different environments. Its templating system allows for easy customization of configurations, while the Helm repository facilitates the sharing and distribution of charts. Whether orchestrating microservices or deploying scalable applications, Helm proves to be an indispensable tool for developers and operators seeking efficiency and consistency in Kubernetes environments.
<br />

### Sample Helm Commands

Below is a list of common Helm CLI commands for managing Kubernetes applications:

#### Initializing Helm

Initialize Helm in your Kubernetes cluster:

```bash
 helm init
```

#### Adding a Helm Repository

Add a repository containing Helm Charts:

```bash
 helm repo add stable https://charts.helm.sh/stable
```

#### Searching for Helm Charts

Searching for Helm charts in the added repositories:

```bash
 helm search repo <keyword>
```

#### Installing a Helm Chart

Install a Helm chart into your Kubernetes cluster:

```bash
helm install <release-name> <chart-name>
```

#### Listing Installed Helm Releases

List all installed Helm releases:

```bash
helm list
```

#### Getting Information about a Release

Get information about a specific Helm release:

```bash
helm status <release-name>
```

#### Upgrading a Helm Release

Upgrade a Helm release to a new version:

```bash
helm upgrade <release-name> <chart-name>
```

#### Deleting a Helm Release

Delete a Helm release from your Kubernetes cluster:

```bash
helm delete <release-name>
```

#### Viewing Helm Release History

View the history of changes for a Helm release:

```bash
helm history <release-name>
```

#### Uninstalling Helm

To uninstall Helm from your Kubernetes cluster, run:

```bash
helm reset
```

## Further Reading

- Read [helm](https://helm.sh/) Official Documentation
