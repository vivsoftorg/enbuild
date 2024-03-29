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

The [Quick install of ENBUILD](docs/how-to-guides/deploying-enbuild-for-local-testing/) is good for local testing. 
Where you can access the ENBUILD UI using the port-forwarding of the ENBUILD UI service.

In production scenarios this is not feasible. There are following options available options to expose the UI service out of the kubernetes cluster.

# Expose UI using kubernetes service type LoadBalancer
Just set the service type to LB
# Use the service type as NodePort
Set the the service type as NodePort and access it using the node port
# Use the ingress controller 
For this you need to have the ingress controller installed and configure in your cluster. 
And then expose the service using ingress
# Expose using istio virtual service
For this you need to have the istio installed and configured in your cluster. 
Refer [instaling istio](docs/how-to-guides/installing-istio/) for steps to installing the istio. 

For exposing the UI service using istio virtual service you need to set the `istio.enabled` to true and need to provide the [refer example input file](https://github.com/vivsoftorg/enbuild/blob/main/examples/enbuild/with_istio.yaml)
