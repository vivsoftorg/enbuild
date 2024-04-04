---
title: "Kubectl Cheat Sheet"
description: "Kubectl Cheat Sheet"
summary: ""
draft: false
menu:
  docs:
    parent: "/troubelshooting/"
    identifier: "kubectl"
weight: 301
toc: true
seo:
  title: "" # custom title (optional)
  description: "" # custom description (recommended)
  canonical: "" # custom canonical URL (optional)
  noindex: false # false (default) or true
---

This cheat sheet provides commonly used `kubectl` commands for managing Kubernetes clusters.

***Basic Commands***

***List all pods***

```bash
kubectl get pods
```

***List all pods in a specific namespace***
```bash
kubectl get pods -n <namespace>
```

***List all nodes***
```bash
kubectl get nodes
```

***List all services***
```bash
kubectl get services
```

# How to check logs of the pods

```bash
kubectl logs <pod-name>
```

# How to check logs of the pods in a specific namespace

```bash
kubectl logs <pod-name> -n <namespace>
```

# How to check logs of the pods in a specific container

```bash
kubectl logs <pod-name> -c <container-name>
```

# How to check logs of the pods in a specific container in a specific namespace

```bash
kubectl logs <pod-name> -c <container-name> -n <namespace>
```

# How to check logs of the pods in a specific container in a specific namespace with timestamps

```bash
kubectl logs <pod-name> -c <container-name> -n <namespace> --timestamps
```

# How to check logs of the pods in a specific container in a specific namespace with timestamps and limit the output to 10 lines

```bash
kubectl logs <pod-name> -c <container-name> -n <namespace> --timestamps --tail 10
```