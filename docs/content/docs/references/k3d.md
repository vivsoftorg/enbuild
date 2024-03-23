---
title: "K3D"
description: "Reference for K3D"
summary: ""
date: 2023-09-07T16:13:18+02:00
lastmod: 2023-09-07T16:13:18+02:00
draft: false
menu:
  docs:
    parent: "/reference/"
    identifier: "k3d"
weight: 910
toc: true
seo:
  title: "" # custom title (optional)
  description: "" # custom description (recommended)
  canonical: "" # custom canonical URL (optional)
  noindex: false # false (default) or true
---

<picture><img src="/images/references/k3d.png" alt="Screenshot of K3D Logo" width="30%" height="30%"></img></picture>
<br />
K3D, a lightweight and versatile tool, simplifies the management and deployment of Kubernetes clusters by bringing the power of Kubernetes into a single-node environment. Designed for simplicity and speed, K3D allows developers and operators to spin up Kubernetes clusters with ease, making it an ideal choice for local development, testing, and CI/CD pipelines. Leveraging containerd and other containerization technologies, K3D offers a minimalistic yet efficient Kubernetes experience. Users can create, scale, and delete clusters effortlessly, making it a valuable tool for scenarios where resource constraints or rapid cluster provisioning are crucial. With K3D, developers can focus on building and testing applications in a Kubernetes-like environment without the complexity of managing large-scale clusters, thereby accelerating the development lifecycle.

### Sample K3D Commands

Below is a list of basic K3D CLI commands for managing Kubernetes clusters:

#### Creating a Kubernetes Cluster

Create a new Kubernetes cluster using k3d:

```bash
 k3d cluster create <cluster-name>
```

#### Listing Kubernetes Clusters

List all existing Kubernetes clusters managed by k3d:

```bash
 k3d cluster list
```

#### Getting Information about a Cluster

Retrieve detailed information about a specific Kubernetes cluster:

```bash
 k3d cluster get <cluster-name>
```

#### Accessing Kubernetes Cluster

Set the kubeconfig context to the newly created Kubernetes cluster:

```bash
 export KUBECONFIG="$(k3d kubeconfig write <cluster-name>)"
```

#### Deleting a Kubernetes Cluster

Delete a Kubernetes cluster managed by k3d:

```bash
 k3d cluster delete <cluster-name>
```

#### Starting a Kubernetes Cluster

Start a previously stopped Kubernetes cluster:

```bash
 k3d cluster start <cluster-name>
```

#### Stopping a Kubernetes Cluster

Stop a running Kubernetes cluster:

```bash
 k3d cluster stop <cluster-name>
```

#### Scaling Nodes

Scale the number of worker nodes in the cluster:

```bash
 k3d node create <node-name> --replicas <num-replicas>
```

#### Exporting kubeconfig

Export the kubeconfig file for a cluster:

```bash
 k3d kubeconfig write <cluster-name>
```

## Further Reading

- Read [k3d](https://k3d.io/) Official Documentation
