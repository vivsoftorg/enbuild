---
title: "Platform One Big Bang"
description: "Reference for Platform One's Big Bang"
summary: ""
draft: false
menu:
  docs:
    parent: "/reference/"
    identifier: "bigbang"
weight: 502
toc: true
seo:
  title: "" # custom title (optional)
  description: "" # custom description (recommended)
  canonical: "" # custom canonical URL (optional)
  noindex: false # false (default) or true
---

<picture><img src="/images/emma/emma-grogu.png" alt="Screenshot of Platform One" width="50%" height="50%"></img></picture>
<br />
<br />
Big Bang is a declarative, continuous delivery tool for deploying DoD hardened and approved packages into a Kubernetes cluster.

## Usage & Scope

Big Bang's scope is to provide publicly available installation manifests for packages required to adhere to the DoD DevSecOps Reference Architecture and additional useful utilities. Big Bang packages are broken into three categories:

- Core: [Core packages](https://github.com/DoD-Platform-One/bigbang/blob/master/docs/understanding-bigbang/package-architecture/README.md##Core) are a group of capabilities required by the DoD DevSecOps Reference Architecture, that are supported directly by the Big Bang development team. The specific capabilities that are considered core currently are Service Mesh, Policy Enforcement, Logging, Monitoring, and Runtime Security.

- Addons: [Addon packages](https://github.com/DoD-Platform-One/bigbang/blob/master/docs/understanding-bigbang/package-architecture/README.md##Addons) are any packages/capabilities that the Big Bang development team directly supports that do not fall under the above core definition. These serve to extend the functionality/features of Big Bang.

- Community: [Community packages](https://repo1.dso.mil/big-bang/product/community) are any packages that are maintained by the broader Big Bang community (users, vendors, etc). These packages could be alternatives to core or addon packages, or even entirely new packages to help extend usage/functionality of Big Bang.

In order for an installation of Big Bang to be a valid installation/configuration you must install/deploy a core package of each category (for additional details on categories and options see [here](https://github.com/DoD-Platform-One/bigbang/blob/master/docs/understanding-bigbang/package-architecture/README.md##Core)).

Big Bang also builds tooling around the testing and validation of Big Bang packages. These tools are provided as-is, without support.

Big Bang is intended to be used for deploying and maintaining a DoD hardened and approved set of packages into a Kubernetes cluster. Deployment and configuration of ingress/egress, load balancing, policy auditing, logging, monitoring, etc. are handled via Big Bang. Additional packages (e.g. ArgoCD, GitLab) can also be enabled and customized to extend Big Bang's baseline. Once deployed, the Kubernetes cluster can be used to add mission specific applications.

## Further Reading

- Read [big bang](https://docs-bigbang.dso.mil/latest/) Official Documentation
- Checkout [big bang](https://github.com/DoD-Platform-One/bigbang) on GitHub
