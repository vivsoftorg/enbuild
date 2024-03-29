---
title: "Kubectl Cheat Sheet"
description: "Kubectl Cheat Sheet"
summary: ""
draft: false
menu:
  docs:
    parent: "/troubelshooting/"
    identifier: "kubectl"
weight: 901
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
