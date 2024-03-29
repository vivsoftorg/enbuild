---
title: "Install and Configure Keycloak for ENBUILD SSO"
description: "Install and Configure Keycloak for ENBUILD SSO"
summary: "Install and Configure Keycloak for ENBUILD SSO"
draft: false
menu:
  docs:
    parent: ""
    identifier: "installKeycloak"
weight: 807
toc: true
seo:
  title: "" # custom title (optional)
  description: "" # custom description (recommended)
  canonical: "" # custom canonical URL (optional)
  noindex: false # false (default) or true
---

# Prerequisites
Before you begin, ensure that you have the following prerequisites in place:
- Access to Platform1 registry
- Kubernetes cluster is up and running and you have access to it via `kubectl` command.
- Helm 3 installed on your system.


## Login to Platform1 registry
Login to Platform1 registry by using the following command:
  

```shell
helm registry login registry1.dso.mil/bigbang
WARNING: Kubernetes configuration file is group-readable. This is insecure. Location: /root/.kube/config
WARNING: Kubernetes configuration file is world-readable. This is insecure. Location: /root/.kube/config
Username: Juned_Memon
Password:
Login Succeeded
```