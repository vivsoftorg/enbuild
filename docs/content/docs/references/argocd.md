---
title: "ArgoCD"
description: "Reference for ArgoCD"
summary: ""
date: 2023-09-07T16:13:18+02:00
lastmod: 2023-09-07T16:13:18+02:00
draft: false
menu:
  docs:
    parent: "/reference/"
    identifier: "argocd"
weight: 910
toc: true
seo:
  title: "" # custom title (optional)
  description: "" # custom description (recommended)
  canonical: "" # custom canonical URL (optional)
  noindex: false # false (default) or true
---

<picture><img src="/images/references/argocd.png" alt="Screenshot of ArgoCD Logo" width="30%" height="30%"></img></picture>
<br />
<br />
ArgoCD, a declarative, GitOps continuous delivery tool for Kubernetes, revolutionizes the application deployment process. At its essence, ArgoCD enables organizations to manage and automate the deployment of applications through Git repositories, promoting a declarative approach to configuration and versioning. It continuously monitors the desired state of applications defined in Git and automatically reconciles any divergences with the current state in the Kubernetes cluster. This ensures that applications are consistently deployed and maintained across different environments. ArgoCD's intuitive user interface provides visibility into the deployment status, allowing for easy tracking and rollbacks. With support for multiple clusters and repositories, ArgoCD empowers teams to achieve efficient, scalable, and auditable continuous delivery workflows in Kubernetes environments.

### Sample ArgoCD Commands

Below is a list of common ArgoCD CLI commands for managing Kubernetes applications:

#### Logging in to Argo CD Server

To log in to the Argo CD server, use the following command:

```bash
 argocd login <argocd-server-url> --username <username> --password <password>
```

#### Setting the Current Context

Before executing any Argo CD CLI commands, you need to set the current context to the Argo CD server:

```bash
 argocd context <argocd-server-url>
```

#### Listing Applications

List all applications managed by Argo CD:

```bash
 argocd app list
```

#### Getting Information about an Application

Retrieve detailed information about a specific application:

```bash
 argocd app get <application-name>
```

#### Syncing an Application

Manually sync an application with its target state:

```bash
 argocd app sync <application-name>
```

#### Setting Sync Options

You can set various sync options for an application. For example:

```bash
 argocd app set <application-name> --sync-policy <policy>
```

#### Deleting an Application

Delete an application from Argo CD:

```bash
 argocd app delete <application-name>
```

#### Configuring Auto-Sync

Enable or disable auto-sync for an application:

```bash
 argocd app auto-sync <application-name> --<enable/disable>
```

#### Viewing Application Resources

List Kubernetes resources managed by an application:

```bash
 argocd app resources <application-name>
```

#### Accessing Argo CD Web UI

You can also access the Argo CD Web UI by running:

```bash
 argocd open
```

#### Logging out of Argo CD Server

To log out of the Argo CD server, run:

```bash
 argocd logout
```

## Further Reading

- Read [argocd](https://argoproj.github.io/cd/) Official Documentation
