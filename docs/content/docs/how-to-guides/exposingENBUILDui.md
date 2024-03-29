---
title: "Exposing ENBUILD UI"
description: "Exposing ENBUILD UI"
summary: "Steps to Install Istio"
draft: false
menu:
  docs:
    parent: "docs/how-to-guides/deploying-enbuild-for-local-testing/"
    identifier: "ExposingENBUILDUI"
weight: 803
toc: true
seo:
  title: "" # custom title (optional)
  description: "" # custom description (recommended)
  canonical: "" # custom canonical URL (optional)
  noindex: false # false (default) or true
---
# Introduction:
In production Kubernetes environments, exposing the ENBUILD UI service requires careful consideration to ensure accessibility and security. While the [Quick install of ENBUILD](docs/how-to-guides/deploying-enbuild-for-local-testing/) facilitates local testing through port forwarding, deploying in a production scenario demands a more robust approach. This document outlines various options available for exposing the ENBUILD UI service outside the Kubernetes cluster.

## Option 1: Expose UI using Kubernetes Service Type LoadBalancer

Setting the service type to LoadBalancer enables external access to the ENBUILD UI service.
Simply configure the service type as LB to allow external traffic. 
Refer to the [example helm input file](https://github.com/vivsoftorg/enbuild/blob/main/examples/enbuild/loadbalancer.yaml) for guidance.

## Option 2: Use Service Type NodePort

Configuring the service type as NodePort provides accessibility by exposing a specific port on all nodes in the cluster.
Access the ENBUILD UI using the designated node port.

Refer to the [example helm input file](https://github.com/vivsoftorg/enbuild/blob/main/examples/enbuild/nodePort.yaml) for guidance.

## Option 3: Use Ingress Controller

Installation and configuration of an Ingress controller within the Kubernetes cluster are prerequisites for this option.
Expose the ENBUILD UI service through Ingress configuration for enhanced routing and management of external traffic.
Refer to the [example helm input file](https://github.com/vivsoftorg/enbuild/blob/main/examples/enbuild/with_ingress.yaml) for guidance.

## Option 4: Expose Using Istio Virtual Service

(Istio installation)[(docs/how-to-guides/installing-istio/)] and configuration are required for leveraging this option.
Set the `istio.enabled` parameter to true and provide the necessary configurations, such as the Istio Virtual Service, to expose the ENBUILD UI service.
Refer to the [example helm input file](https://github.com/vivsoftorg/enbuild/blob/main/examples/enbuild/with_istio.yaml) for guidance.

Refer the [detailed guide of installing the ENBUILD on top of istio](docs/how-to-guides/deploying-enbuild-exposing-the-service-using-istio/)