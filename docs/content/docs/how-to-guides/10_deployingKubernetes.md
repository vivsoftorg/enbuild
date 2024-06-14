---
title: "Deploying Kubernetes"
description: "Steps to Deploy Kubernetes "
summary: "Steps to Deploy Kubernets"
draft: false
menu:
  docs:
    parent: "docs/how-to-guides/"
    identifier: ""
weight: 203
toc: true
seo:
  title: "" # custom title (optional)
  description: "" # custom description (recommended)
  canonical: "" # custom canonical URL (optional)
  noindex: false # false (default) or true
---

# Deploy BigBang

-  Login to Enbuild -[Enbuild](https://enbuild.vivplatform.io)
-  Click on the **Home**
-  Select the **Kubernetes**
-  Choose the component **RKE2** from the **DISTRO** category and click on the **VALUES** tab and provide the `Credentials`

<picture><img src="/images/how-to-guides/RKE2.png" alt="RKE2"></img></picture>

- After providing all the input values, provide the name for your deployment proceed to Infrastructure section,
- Select AWS as your cloud and provide your AWS credentials.

<picture><img src="/images/how-to-guides/AWS_Kubernetes.png" alt="AWS_Kubernetes"></img></picture>

Once all inputs are provided click on **Deploy Stack**

- Choose the component **EKS** from the **DISTRO** category and follow the previous steps.

<picture><img src="/images/how-to-guides/EKS.png" alt="EKS"></img></picture>